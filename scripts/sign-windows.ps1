param([Parameter(Mandatory)][string]$Path, [switch]$Check)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'Windows and PowerShell 7 are required' }
foreach ($name in @('WINDOWS_CERTIFICATE_PFX_BASE64', 'WINDOWS_CERTIFICATE_PASSWORD', 'WINDOWS_TIMESTAMP_URL')) {
    if ([string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($name))) { throw "Missing required secret/configuration: $name" }
}
if ($env:WINDOWS_TIMESTAMP_URL -notmatch '^https://[^\s]+$') { throw 'WINDOWS_TIMESTAMP_URL must be an HTTPS RFC3161 endpoint' }
$signTool = Get-Command signtool.exe -ErrorAction SilentlyContinue
if (-not $signTool) { throw 'Install the Windows SDK signing tools and add signtool.exe to PATH' }
if ($Check) { Write-Host 'Windows signing prerequisites present (certificate not authenticated).'; return }
$target = (Resolve-Path -LiteralPath $Path).Path
$work = Join-Path ([IO.Path]::GetTempPath()) ('takedown-sign-' + [guid]::NewGuid())
New-Item -ItemType Directory -Path $work | Out-Null
try {
    # Restrict the temporary PFX directory to the current identity, before writing it.
    $acl = Get-Acl -LiteralPath $work
    $acl.SetAccessRuleProtection($true, $false)
    $rule = [Security.AccessControl.FileSystemAccessRule]::new(
        [Security.Principal.WindowsIdentity]::GetCurrent().User, 'FullControl',
        'ContainerInherit,ObjectInherit', 'None', 'Allow')
    $acl.AddAccessRule($rule)
    Set-Acl -LiteralPath $work -AclObject $acl
    $pfx = Join-Path $work 'signing.pfx'
    [IO.File]::WriteAllBytes($pfx, [Convert]::FromBase64String($env:WINDOWS_CERTIFICATE_PFX_BASE64))
    # SignTool's PFX password interface uses argv. Use an isolated runner, no tracing/transcripts.
    & $signTool.Source sign /fd SHA256 /td SHA256 /tr $env:WINDOWS_TIMESTAMP_URL /f $pfx /p $env:WINDOWS_CERTIFICATE_PASSWORD $target
    if ($LASTEXITCODE -ne 0) { throw "signtool sign failed: $LASTEXITCODE" }
    & $signTool.Source verify /pa /all /v $target
    if ($LASTEXITCODE -ne 0) { throw "signtool verify failed: $LASTEXITCODE" }
    $signature = Get-AuthenticodeSignature -LiteralPath $target
    if ($signature.Status -ne 'Valid' -or -not $signature.TimeStamperCertificate) {
        throw 'A valid Authenticode signature and timestamp are required'
    }
} finally {
    Remove-Item -LiteralPath $work -Recurse -Force
}
