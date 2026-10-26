package storage

import (
	"context"
	"errors"
	"math"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const connectionTimeout = 30 * time.Second

func timeNowPlusCloseTimeout() time.Time { return time.Now().Add(5 * time.Second) }

func remotePath(u *url.URL, smb bool) (root, share string, err error) {
	if !strings.HasPrefix(u.Path, "/") {
		return "", "", ErrUnsafePath
	}
	rel := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/"), "/")
	if rel != "" {
		if err := validateName(rel); err != nil {
			return "", "", err
		}
	}
	if smb {
		if rel == "" {
			return "", "", errors.New("SMB destination requires a share")
		}
		parts := strings.SplitN(rel, "/", 2)
		if len(parts) == 1 {
			return ".", parts[0], nil
		}
		return parts[1], parts[0], nil
	}
	return path.Join("/", rel), "", nil
}

func remoteAddress(u *url.URL, defaultPort string) string {
	port := u.Port()
	if port == "" {
		port = defaultPort
	}
	return net.JoinHostPort(u.Hostname(), port)
}

func openSMB(ctx context.Context, u *url.URL, opts Options) (StorageBackend, error) {
	root, shareName, err := remotePath(u, true)
	if err != nil {
		return nil, err
	}
	conn, err := (&net.Dialer{Timeout: connectionTimeout}).DialContext(ctx, "tcp", remoteAddress(u, "445"))
	if err != nil {
		return nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(connectionTimeout))
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()
	dialer := smb2.Dialer{Initiator: &smb2.NTLMInitiator{User: u.User.Username(), Password: opts.Password, Domain: opts.Domain}}
	session, err := dialer.DialContext(ctx, conn)
	if err != nil {
		_ = conn.Close()
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	share, err := session.WithContext(ctx).Mount(shareName)
	if err != nil {
		_ = session.Logoff()
		_ = conn.Close()
		return nil, err
	}
	b := &remoteBackend{root: root, base: ".", conn: conn, allowUnknown: opts.AllowUnknownSpace}
	b.ops = remoteOps{
		lstat: func(ctx context.Context, name string) (os.FileInfo, error) {
			return share.WithContext(ctx).Lstat(smbPath(name))
		},
		mkdir: func(ctx context.Context, name string) error { return share.WithContext(ctx).Mkdir(smbPath(name), 0700) },
		open: func(ctx context.Context, name string, flags int) (archiveFile, error) {
			return share.WithContext(ctx).OpenFile(smbPath(name), flags, 0600)
		},
		rename: func(ctx context.Context, from, to string) error {
			return share.WithContext(ctx).Rename(smbPath(from), smbPath(to))
		},
		free: func(ctx context.Context, name string) (uint64, error) {
			info, err := share.WithContext(ctx).Statfs(smbPath(name))
			if err != nil {
				return 0, err
			}
			return availableSMBSpace(info)
		},
		close: func() error {
			return errors.Join(share.WithContext(context.Background()).Umount(), session.WithContext(context.Background()).Logoff())
		},
	}
	if err := b.ensureDirectories(ctx, root, true); err != nil {
		_ = b.Close()
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		_ = b.Close()
		return nil, err
	}
	_ = conn.SetDeadline(time.Time{})
	return b, nil
}

func smbPath(name string) string {
	if name == "." {
		return ""
	}
	return strings.ReplaceAll(name, "/", `\`)
}

func availableSMBSpace(info smb2.FileFsInfo) (uint64, error) {
	if info.BlockSize() == 0 || info.FragmentSize() == 0 || info.AvailableBlockCount() > math.MaxInt64 {
		return 0, &SpaceUnknownError{errors.New("invalid SMB allocation units")}
	}
	return spaceProduct(info.AvailableBlockCount(), info.BlockSize(), info.FragmentSize())
}

func availableSFTPSpace(info *sftp.StatVFS) (uint64, error) {
	if info.Frsize == 0 {
		return 0, &SpaceUnknownError{errors.New("invalid SFTP fragment size")}
	}
	return spaceProduct(info.Bavail, info.Frsize)
}

func sshConfig(u *url.URL, opts Options) (*ssh.ClientConfig, error) {
	knownPath := opts.KnownHostsPath
	if knownPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		knownPath = filepath.Join(home, ".ssh", "known_hosts")
	}
	callback, err := knownhosts.New(knownPath)
	if err != nil {
		return nil, errors.New("SSH known_hosts could not be loaded; configure KnownHostsPath with a trusted host key")
	}
	var auth []ssh.AuthMethod
	if opts.PrivateKeyPath != "" {
		data, err := os.ReadFile(opts.PrivateKeyPath)
		if err != nil {
			return nil, errors.New("SSH private key could not be read")
		}
		signer, err := ssh.ParsePrivateKey(data)
		clear(data)
		if err != nil {
			return nil, errors.New("SSH private key could not be parsed (encrypted keys require an unlocked key file)")
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if opts.Password != "" {
		auth = append(auth, ssh.Password(opts.Password))
	}
	if len(auth) == 0 {
		return nil, errors.New("SFTP requires a password or private key")
	}
	return &ssh.ClientConfig{User: u.User.Username(), Auth: auth, HostKeyCallback: callback, Timeout: connectionTimeout}, nil
}

func openSFTP(ctx context.Context, u *url.URL, opts Options) (StorageBackend, error) {
	root, _, err := remotePath(u, false)
	if err != nil {
		return nil, err
	}
	config, err := sshConfig(u, opts)
	if err != nil {
		return nil, err
	}
	address := remoteAddress(u, "22")
	conn, err := (&net.Dialer{Timeout: connectionTimeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}
	_ = conn.SetDeadline(time.Now().Add(connectionTimeout))
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()
	cc, channels, requests, err := ssh.NewClientConn(conn, address, config)
	if err != nil {
		_ = conn.Close()
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	sshClient := ssh.NewClient(cc, channels, requests)
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		_ = sshClient.Close()
		_ = conn.Close()
		return nil, err
	}
	b := &remoteBackend{root: root, base: "/", conn: conn, allowUnknown: opts.AllowUnknownSpace}
	b.ops = remoteOps{
		lstat: func(ctx context.Context, name string) (os.FileInfo, error) {
			return remoteCall(ctx, conn, func() (os.FileInfo, error) { return client.Lstat(name) })
		},
		mkdir: func(ctx context.Context, name string) error {
			return remoteDo(ctx, conn, func() error { return client.Mkdir(name) })
		},
		open: func(ctx context.Context, name string, flags int) (archiveFile, error) {
			return remoteCall(ctx, conn, func() (archiveFile, error) { return client.OpenFile(name, flags) })
		},
		rename: func(ctx context.Context, from, to string) error {
			return remoteDo(ctx, conn, func() error {
				if _, ok := client.HasExtension("posix-rename@openssh.com"); ok {
					return client.PosixRename(from, to)
				}
				return client.Rename(from, to)
			})
		},
		free: func(ctx context.Context, name string) (uint64, error) {
			if _, ok := client.HasExtension("statvfs@openssh.com"); !ok {
				return 0, &SpaceUnknownError{errors.New("SFTP server does not advertise statvfs")}
			}
			return remoteCall(ctx, conn, func() (uint64, error) {
				info, err := client.StatVFS(name)
				if err != nil {
					return 0, err
				}
				return availableSFTPSpace(info)
			})
		},
		close: func() error { return errors.Join(closeResult(client.Close()), closeResult(sshClient.Close())) },
	}
	if err := b.ensureDirectories(ctx, root, true); err != nil {
		_ = b.Close()
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		_ = b.Close()
		return nil, err
	}
	_ = conn.SetDeadline(time.Time{})
	return b, nil
}
