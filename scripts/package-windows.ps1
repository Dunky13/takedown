param(
    [ValidatePattern('^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$')][string]$Version = '0.0.0',
    [ValidateSet('local', 'release')][string]$Mode = 'local',
    [string]$Binary = '',
    [string]$OutputDirectory = ''
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'Run natively on Windows with PowerShell 7' }
$root = Split-Path $PSScriptRoot -Parent
Push-Location $root
$work = $null
try {
    foreach ($component in $Version.Split('.')) {
        if ([decimal]$component -gt 65535) { throw 'Windows package version components must be at most 65535' }
    }
    if ((& go env GOVERSION) -ne 'go1.26.10' -or $LASTEXITCODE -ne 0) { throw 'Go 1.26.10 is required' }
    if ((& go list -m -mod=readonly -f '{{.Version}}' fyne.io/fyne/v2) -ne 'v2.8.2' -or $LASTEXITCODE -ne 0) { throw 'Fyne v2.8.2 is required' }
    $arch = & go env GOHOSTARCH
    if ($arch -ne 'amd64' -or $LASTEXITCODE -ne 0) { throw 'The Windows installer currently supports native amd64 only' }
    $nsis = Get-Command makensis.exe -ErrorAction SilentlyContinue
    if (-not $nsis) { throw 'Install NSIS 3.11 and add makensis.exe to PATH' }
    $nsisVersion = & $nsis.Source /VERSION
    if ($LASTEXITCODE -ne 0 -or $nsisVersion -notmatch '^v3\.') { throw 'NSIS 3.x Unicode is required (CI uses 3.11)' }
    if ($Mode -eq 'release') { & "$PSScriptRoot/sign-windows.ps1" -Path unused -Check }
    if (-not $OutputDirectory) { $OutputDirectory = "dist/packages/windows-$arch-$Mode" }
    $output = [IO.Path]::GetFullPath($OutputDirectory)
    if (Test-Path -LiteralPath $output) { throw 'OutputDirectory already exists; choose a fresh directory' }
    $parent = Split-Path $output -Parent
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    $work = Join-Path $parent ('.package-' + [guid]::NewGuid())
    $stage = Join-Path $work 'payload'
    $artifacts = Join-Path $work 'output'
    New-Item -ItemType Directory -Path $stage, $artifacts | Out-Null
    $exe = Join-Path $stage 'takedown.exe'
    if ($Binary) {
        Copy-Item -LiteralPath $Binary -Destination $exe
    } else {
        if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) { throw 'Install a native MinGW-w64 gcc C compiler and add it to PATH' }
        $previousCgo = $env:CGO_ENABLED
        $previousOS = $env:GOOS
        $previousArch = $env:GOARCH
        try {
            $env:CGO_ENABLED = '1'; $env:GOOS = 'windows'; $env:GOARCH = $arch
            & go build -mod=readonly -p 2 -trimpath -ldflags "-X github.com/Dunky13/takedown/internal/cli.Version=$Version" -o $exe ./cmd/takedown
            if ($LASTEXITCODE -ne 0) { throw "Native build failed: $LASTEXITCODE" }
        } finally { $env:CGO_ENABLED = $previousCgo; $env:GOOS = $previousOS; $env:GOARCH = $previousArch }
    }
    $info = & go version -m $exe
    if ($LASTEXITCODE -ne 0) { throw 'Could not read binary build metadata' }
    $joined = $info -join "`n"
    foreach ($pattern in @(': go1\.26\.10\r?$', 'dep\s+fyne.io/fyne/v2\s+v2.8.2\s', 'CGO_ENABLED=1', 'GOOS=windows', 'GOARCH=amd64')) {
        if ($joined -notmatch "(?m)$pattern") { throw "Binary is not the required native desktop build: $pattern" }
    }
    if ($joined -match '-tags=.*(headless|ci)') { throw 'Headless/ci tagged binaries cannot be packaged as desktop apps' }
    foreach ($license in @('LICENSE', 'NOTICE', 'THIRD_PARTY_NOTICES.md')) {
        $item = Get-Item -LiteralPath $license -Force
        if ($item.PSIsContainer -or ($item.Attributes -band [IO.FileAttributes]::ReparsePoint)) { throw "License payload must be a regular file: $license" }
    }
    Copy-Item -LiteralPath 'README.md', 'PRIVACY.md', 'PRIVACY.html', 'LICENSE', 'NOTICE', 'THIRD_PARTY_NOTICES.md' -Destination $stage
    foreach ($license in @('LICENSE.md', 'LICENSE.txt', 'COPYING')) {
        if (Test-Path -LiteralPath $license) { Copy-Item -LiteralPath $license -Destination $stage }
    }
    $metadata = @("Package version: $Version", "Mode: $Mode") + $info[1..($info.Count - 1)]
    $metadata | Set-Content -LiteralPath (Join-Path $stage 'BUILD-INFO.txt') -Encoding utf8
    if ($Mode -eq 'release') { & "$PSScriptRoot/sign-windows.ps1" -Path $exe }
    $base = "takedown-$Version-windows-$arch-$Mode"
    $installer = Join-Path $artifacts "$base-setup.exe"
    $nsisArgs = @('/V2', "/DVERSION=$Version", "/DPAYLOAD=$stage", "/DOUTPUT=$installer")
    if ($Mode -eq 'release') { $nsisArgs += "/DSIGN_SCRIPT=$PSScriptRoot/sign-windows.ps1" }
    & $nsis.Source @nsisArgs "$root/build/windows/installer.nsi"
    if ($LASTEXITCODE -ne 0) { throw "NSIS installer build failed: $LASTEXITCODE" }
    if ($Mode -eq 'release') { & "$PSScriptRoot/sign-windows.ps1" -Path $installer }
    Compress-Archive -Path "$stage/*" -DestinationPath (Join-Path $artifacts "$base.zip") -CompressionLevel Optimal
    $sums = foreach ($file in (Get-ChildItem -LiteralPath $artifacts -File | Sort-Object Name)) {
        '{0}  {1}' -f (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant(), $file.Name
    }
    $sums | Set-Content -LiteralPath (Join-Path $artifacts 'SHA256SUMS') -Encoding ascii
    Move-Item -LiteralPath $artifacts -Destination $output
    Write-Host "Artifacts: $output"
} finally {
    if ($work -and (Test-Path -LiteralPath $work)) { Remove-Item -LiteralPath $work -Recurse -Force }
    Pop-Location
}
