param([Parameter(Mandatory)][string]$Installer)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
if (-not $IsWindows -or $env:GITHUB_ACTIONS -ne 'true' -or -not $env:RUNNER_TEMP) {
    throw 'This installation regression is restricted to a disposable Windows Actions runner'
}
$installed = Join-Path $env:LOCALAPPDATA 'Programs/takedown'
$key = 'HKCU:/Software/Microsoft/Windows/CurrentVersion/Uninstall/takedown'
$shortcut = Join-Path ([Environment]::GetFolderPath('Programs')) 'takedown/takedown.lnk'
if ((Test-Path -LiteralPath $installed) -or (Test-Path $key) -or (Test-Path -LiteralPath $shortcut)) {
    throw 'Refusing to replace existing application files or registration'
}
$setup = (Resolve-Path -LiteralPath $Installer).Path
$work = Join-Path $env:RUNNER_TEMP ('takedown-uninstall-test-' + [guid]::NewGuid())
New-Item -ItemType Directory -Path $work | Out-Null
$handle = $null
try {
    $process = Start-Process -FilePath $setup -ArgumentList '/S' -Wait -PassThru
    if ($process.ExitCode -ne 0) { throw 'Installer failed' }
    & "$installed/takedown.exe" --version
    if ($LASTEXITCODE -ne 0) { throw 'Installed executable failed' }
    $registration = Get-ItemProperty $key
    if ($registration.InstallLocation -ne $installed) { throw 'Incorrect per-user InstallLocation' }
    $paths = @("$installed/takedown.exe", "$installed/uninstall.exe", $shortcut)
    foreach ($name in @('LICENSE', 'NOTICE', 'THIRD_PARTY_NOTICES.md')) {
        $path = Join-Path $installed $name
        $source = Join-Path (Join-Path $PSScriptRoot '..') $name
        if ((Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash -ne (Get-FileHash -LiteralPath $source -Algorithm SHA256).Hash) {
            throw "Installed license payload differs from source: $name"
        }
        $paths += $path
    }
    $hashes = @{}
    foreach ($path in $paths) { $hashes[$path] = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash }
    # Execute a copied worker with _?= LAST and unquoted, as NSIS documents.
    # This suppresses the asynchronous self-copy launcher; ExitCode is the worker's.
    $worker = Join-Path $work 'uninstall-worker.exe'
    Copy-Item -LiteralPath "$installed/uninstall.exe" -Destination $worker
    $handle = [IO.File]::Open("$installed/takedown.exe", 'Open', 'Read', 'Read')
    $process = Start-Process -FilePath $worker -ArgumentList "/S _?=$installed" -Wait -PassThru
    if ($process.ExitCode -ne 10) { throw "Locked payload must fail with worker exit 10, got $($process.ExitCode)" }
    foreach ($path in $paths) {
        if ((Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash -ne $hashes[$path]) { throw "Locked uninstall changed recovery file: $path" }
    }
    $after = Get-ItemProperty $key
    foreach ($field in @('DisplayName', 'DisplayVersion', 'InstallLocation', 'UninstallString')) {
        if ($after.$field -ne $registration.$field) { throw "Locked uninstall changed registration: $field" }
    }
    $handle.Dispose(); $handle = $null
    # Retry through the retained uninstaller, again using its actual copied worker.
    Copy-Item -LiteralPath "$installed/uninstall.exe" -Destination $worker -Force
    $process = Start-Process -FilePath $worker -ArgumentList "/S _?=$installed" -Wait -PassThru
    if ($process.ExitCode -ne 0) { throw "Unlocked worker failed: $($process.ExitCode)" }
    foreach ($path in $paths) { if (Test-Path -LiteralPath $path) { throw "Owned file survived completed uninstall: $path" } }
    if (Test-Path $key) { throw 'Completed uninstall retained registry entry' }
    Write-Host 'PASS: locked payload preserved recovery metadata; unlocked worker removed owned files/registration'
} finally {
    if ($handle) { $handle.Dispose() }
    # Failure leaves a usable registered installation for diagnosis; only worker scratch is removed.
    Remove-Item -LiteralPath $work -Recurse -Force
}
