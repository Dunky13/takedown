# Desktop Google OAuth authentication handoff

The auth component now uses Google's documented Desktop app OAuth flow through `golang.org/x/oauth2` v0.37.0 and its Google client parser. The old Rod/browser automation, credential submission, TOTP, cookie/session extraction, profile restoration, Takeout DOM parsing, and request interception implementation and fixtures were removed. Only `internal/auth` was changed by this implementation task. Dependency updates and other packages belong to the parent; no commits, delegated children, real keyring writes, browser login, or live account tests were performed.

## Integration API

```go
type Options struct {
    DataDir, ClientFile, Email string
    Scopes []string
    Bus *events.Bus
    OpenBrowser func(string) error // optional; defaults to OpenURL
}
func New(opts Options) (*Manager, error)
func OpenURL(raw string) error
func (*Manager) Login(context.Context) error
func (*Manager) Refresh(context.Context) error
func (*Manager) Forget(context.Context) error
func (*Manager) ProfileID(context.Context) (string, error)
func (*Manager) HTTPClient() *http.Client
func (*Manager) Close() error
var ErrAuthorizationRequired error
func IsTemporary(error) bool
func RetryDelay(error) time.Duration
```

`New` normalizes the absolute data directory and lowercase/trimmed email hint, copies/sorts/deduplicates scopes, and validates scope syntax without reading a client file, accessing the keyring, creating directories, starting a browser, or using the network. Missing client configuration and empty scopes are allowed during setup. `Login`, `Refresh`, and HTTP operations require configured scopes and an accessible Desktop app client file. They return actionable `ErrAuthorizationRequired` errors when configuration or consent is missing. Rebuild the manager when options change; validated client configuration is cached for its lifetime.

`ProfileID` loads the validated client through the manager gate without a keyring/token operation. It returns lowercase SHA256 hex of the JSON tuple `[absoluteDataDir, clientID, normalizedEmailHint]`, excluding scopes. It works with empty scopes. This identifies a local profile slot, not a verified account: the user can consent with a different account than the hint. Source history can share this ID across resource selections while separating directories, client IDs, and hints.

## Authorization and grant boundaries

Only explicit `Login` opens the system browser via `github.com/pkg/browser`. Google owns account selection, passwords, passkeys/MFA, and consent. Background HTTP and `Refresh` never open a browser. Login uses library-generated independent random state and PKCE verifier, library S256 challenge/verifier options, offline access, explicit consent, and `include_granted_scopes=false`. No identity scopes are appended. Data Portability scopes cannot be mixed with other scopes. `Email` is only a login hint and local storage slot; it is never treated as verified identity.

Client loading accepts a regular Desktop `installed` JSON file, at most 1 MiB, with client ID/secret and the library-required redirect list. Web/service-account documents are rejected. Authorization endpoints must exactly match Google's documented v2 endpoint or its official legacy authorization path; token endpoints must exactly match `https://oauth2.googleapis.com/token`. Endpoints with extra credentials, ports, paths, query strings, or foreign hosts are rejected. Login uses the canonical v2 authorization endpoint and always replaces the downloaded redirect with its own loopback address. Token HTTP redirects are refused, and auth style is fixed to Google's documented form parameters to avoid exchange probing/retries.

The callback binds `127.0.0.1` on an ephemeral IPv4 TCP port and requires the exact Host and `/oauth2/callback` path. It rejects other methods, absolute/encoded paths, malformed/duplicate/unknown query parameters, wrong/missing state, ambiguous code/error responses, and a foreign supplied issuer. Invalid callbacks do not consume the flow; one valid response does. Replies never reflect codes/state/provider descriptions and set no-store, no-referrer, and restrictive CSP headers. The listener closes on completion, denial, cancellation, or manager close.

Successful tokens must be usable bearer tokens with an explicit future expiry. If a token response supplies `scope`, every required scope must be present; empty, partial, null, wrongly typed, duplicated form scope, and mixed Data Portability/other grants are rejected. Omitted scope follows OAuth's same-as-request convention. The token metadata check supplements the official library; exchange, token parsing, refresh protocol, expiry checks, and rotation remain library-owned.

## Token persistence and concurrency

Access and refresh tokens are stored only in the OS keyring through `go-keyring`, under the dedicated OAuth service. The account name hashes the absolute directory, normalized hint, validated client ID, and sorted scopes. Versioned records also bind the required scopes. There is no token file, cookie jar, encrypted-file token store, or plaintext fallback. Keyring failures are actionable and their raw diagnostics are discarded.

A manager gate serializes authentication, token reads/refreshes/persistence, profile lookup, and forgetting. HTTP workers share cached valid tokens. Expired tokens use the official library refresh source. Refresh rotation is saved before any new bearer token is dispatched. If saving a rotation fails, the new token remains pending in memory; the next operation retries storage before using it or performing another refresh. A process exit before that retry can require fresh consent. Independent managers/processes are not coordinated and must not concurrently refresh the same OAuth slot.

API GET/HEAD 401s get one serialized refresh and one replay, comparing the rejected access token and revision so concurrent workers reuse a completed refresh. A second 401, missing refresh token, revoked/invalid grant, client/configuration rejection, or insufficient-scope 403 requires explicit authorization. POSTs and other mutations are never replayed or redirected by auth. API errors requiring authorization publish `events.ReAuthRequired`; `Close` leaves the parent's bus open.

OAuth HTTP 408/429/5xx and recognized network outages return sanitized errors with `net.Error` Timeout/Temporary behavior. `IsTemporary` and `RetryDelay` work through ordinary HTTP error wrappers. OAuth `server_error` and `temporarily_unavailable` are retryable; invalid grants/configuration/consent remain terminal authorization errors. Retry-After seconds and HTTP dates are preserved, including delays beyond 24 hours; numeric overflow saturates at the maximum duration. Auth does not loop on these outages or open a browser.

## HTTP destinations and error handling

Bearer credentials are attached only to exact HTTPS API hosts `www.googleapis.com`, `drive.googleapis.com`, and `dataportability.googleapis.com`. Signed downloads are permitted without bearer credentials only for `storage.googleapis.com`, its subdomains, subdomains of `googleusercontent.com`, and `drive.usercontent.google.com`/subdomains of `usercontent.google.com`. All initial and redirect requests require HTTPS, no embedded credentials/fragment, port 443 or its implicit default, and no Host override. No general `google.com` account-cookie route is allowed.

Authorization, Proxy-Authorization, Cookie, Cookie2, and Referer headers are removed case-insensitively before dispatch; only API destinations then receive the manager's bearer token. Redirects repeat this policy, with at most ten requests in a redirect chain. Account cookies and signed URL referrers are never forwarded. Response bodies retain lifecycle cancellation until EOF/Close; body and transport diagnostics are sanitized.

Auth-owned errors/events omit client secrets, tokens, authorization codes, signed URLs, provider response bodies/descriptions, and keyring/opener diagnostics. One integration limitation follows from the required `*http.Client` API: Go's `Client.Do` wraps transport/redirect errors in `*url.Error` containing the caller's request URL. Callers must sanitize that wrapper before presenting or logging signed-download errors. Do not display raw `err.Error()` from `Do` for a URL with secret query parameters. Error classification through `errors.Is`/`errors.As`, `IsTemporary`, and `RetryDelay` remains available. The transport does not mutate the caller's request URL to hide this standard-library behavior.

## Forgetting legacy secrets

`Forget` clears cached authorization and deletes the currently configured OAuth client/scope/hint/directory slot. It also removes the exact legacy `credentials-<directoryHash>-<emailHash>` keyring entry containing the old password/TOTP, and exact `session-<emailHash>.enc` and `profile-<emailHash>.enc` entries in the data directory. Empty hints can clean the old default session/profile slot. Cleanup uses a confined filesystem root, removes symlinks themselves, refuses directories/special files, and never recursively removes unknown data. Errors from separate deletions are reported while other known cleanup proceeds.

The old `aes256-<directoryHash>` key is retained because it was shared with other account slots. Other hints, scope combinations, client IDs, unknown disposable profiles, and stale client configurations cannot be enumerated safely. Supply the original client/scopes/email/data directory to forget another known OAuth slot, or the original hint/data directory for legacy cleanup. Missing client/scopes do not prevent legacy-only cleanup; an unreadable configured client is reported when it prevents identifying the OAuth slot. Remote Google grants are not revoked by local `Forget`; users can revoke application access in their Google account settings.

## Bounds and validation receipt

Login's callback/network flow is bounded to five minutes, or the caller's earlier cancellation. Each exchange/refresh uses a 30-second deadline. Callback header/write/idle bounds are five seconds with an 8 KiB header limit; query validation caps encoded query size at 16 KiB. The authorization URL is limited to 32 KiB, client and token HTTP response bodies to 1 MiB, and keyring records to 64 KiB. Archive requests use the caller's deadlines rather than a global download timeout. The API GET/HEAD replay limit is one.

The system opener has no context parameter: cancellation stops the OAuth listener/wait, but an already started ordinary browser remains open and a blocked opener call may outlive it. OS keyring calls are synchronous and cannot be interrupted by this library; callers must unlock/configure the keyring. Already dispatched network requests cannot be recalled by `Forget`; `Close` cancels their contexts. Scope changes create different token slots. Email/profile slots do not prove actual consent identity.

Focused checks passed with fake keyring implementations and injected HTTP transports only:

```sh
go test -mod=readonly -p 2 -parallel 2 ./internal/auth
go vet -mod=readonly -p 2 ./internal/auth
go test -race -mod=readonly -p 2 -parallel 2 ./internal/auth
```

The race command ran through the shared heavy-command queue. The suite covers lazy/missing configuration, profile identity, desktop credentials/endpoints, scope isolation and incomplete grants, PKCE exchange, state/CSRF/path/Host/query rejection, one-shot concurrent callbacks, denial, cancellation/close, expired/revoked grants, refresh rotation and persistence failure, concurrent expiry/401 refresh, mutation non-replay, HTTP/cookie/bearer destination boundaries, sanitized temporary errors/Retry-After, and narrowly scoped legacy cleanup. Live Google account, real OS keyring, MFA/consent, and publishing-status acceptance remain user-owned; no broader integrated validation is claimed here.

## Official references

- [Google Desktop OAuth flow](https://developers.google.com/identity/protocols/oauth2/native-app)
- [Data Portability OAuth restrictions and consent identity limits](https://developers.google.com/data-portability/user-guide/configure-oauth)
- [Official Go OAuth library](https://pkg.go.dev/golang.org/x/oauth2)
- [Google desktop client parser](https://pkg.go.dev/golang.org/x/oauth2/google#ConfigFromJSON)
