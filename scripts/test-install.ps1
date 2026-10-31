# Exercise the Windows installer against the built archive without network access.
param([Parameter(Mandatory)][string]$Archive)
$ErrorActionPreference = "Stop"
$archivePath = (Resolve-Path $Archive).Path
$fixture = Join-Path ([IO.Path]::GetTempPath()) ("telia-installer-test-" + [guid]::NewGuid())
New-Item -ItemType Directory $fixture | Out-Null
try {
    $asset = "telia-windows-amd64.tar.gz"
    $checksum = (Get-FileHash $archivePath -Algorithm SHA256).Hash.ToLowerInvariant()
    $sums = Join-Path $fixture "SHA256SUMS"
    Set-Content -Encoding ASCII $sums "$checksum  $asset"
    function Invoke-WebRequest {
        param([string]$Uri, [string]$OutFile)
        if ($Uri.EndsWith("/SHA256SUMS")) { Copy-Item $sums $OutFile }
        elseif ($Uri.EndsWith("/$asset")) { Copy-Item $archivePath $OutFile }
        else { throw "Unexpected download: $Uri" }
    }
    function Invoke-RestMethod {
        param([string]$Uri)
        return @{ tag_name = "v0.0.0-test" }
    }
    $installer = Join-Path $PSScriptRoot "install.ps1"
    $install = Join-Path $fixture "installation with spaces"
    & $installer -InstallRoot $install
    $launcher = Join-Path $install "bin/telia.cmd"
    $output = & $launcher version
    if ($LASTEXITCODE -ne 0 -or $output -notlike 'telia * (*)') { throw "Launcher did not run compiler: $output" }
    if (-not (Test-Path "$install/v0.0.0-test/base/std")) { throw "Bundled libraries were not installed." }
    & $installer -Version "v0.0.1-test" -InstallRoot $install
    if (-not (Test-Path "$install/v0.0.0-test/telia.exe")) { throw "Upgrade removed the old version." }
    if ((Get-Content $launcher -Raw) -notmatch 'v0\.0\.1-test') { throw "Upgrade did not update the launcher." }

    $failed = $false
    try { & $installer -Version "v0.0.1-test" -InstallRoot $install } catch { $failed = $true }
    if (-not $failed) { throw "Duplicate installation should fail." }

    Set-Content -Encoding ASCII $sums "$('0' * 64)  $asset"
    $failed = $false
    try { & $installer -Version "v0.0.2-test" -InstallRoot $install } catch {
        if ($_.Exception.Message -notlike '*Checksum mismatch*') { throw }
        $failed = $true
    }
    if (-not $failed -or (Test-Path "$install/v0.0.2-test")) { throw "Checksum failure installed a bundle." }
    if (Get-ChildItem $install -Filter '.install-*') { throw "Temporary installation files were not cleaned up." }
    Write-Host "Windows installer checks passed"
} finally {
    Remove-Item -Recurse -Force $fixture
}
