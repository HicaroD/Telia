# Install the complete Windows x64 bundle without administrator privileges.
param(
    [string]$Version = "latest",
    [string]$InstallRoot = (Join-Path $env:LOCALAPPDATA "Telia")
)
$ErrorActionPreference = "Stop"

if (-not [Console]::IsOutputRedirected) {
    Write-Host @'
_________ _______  _       _________ _______
\__   __/(  ____ \( \      \__   __/(  ___  )
   ) (   | (    \/| (         ) (   | (   ) |
   | |   | (__    | |         | |   | (___) |
   | |   |  __)   | |         | |   |  ___  |
   | |   | (      | |         | |   | (   ) |
   | |   | (____/\| (____/\___) (___| )   ( |
   )_(   (_______/(_______/\_______/|/     \|

          THE TELIA PROGRAMMING LANGUAGE
'@
}

if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString() -ne "X64") {
    throw "Only Windows x64 releases are currently available."
}
if (-not (Get-Command tar.exe -ErrorAction SilentlyContinue)) { throw "tar.exe is required (included with modern Windows)." }
$repo = "https://github.com/HicaroD/Telia"
if ($Version -eq "latest") {
    $release = Invoke-RestMethod "https://api.github.com/repos/HicaroD/Telia/releases/latest"
    $Version = $release.tag_name
}
if ($Version -notmatch '^v[0-9][A-Za-z0-9.+-]*$') { throw "Expected a version tag such as v0.1.0." }
if (-not [IO.Path]::IsPathRooted($InstallRoot)) { throw "InstallRoot must be an absolute path." }
$destination = Join-Path $InstallRoot $Version
if (Test-Path $destination) { throw "$destination already exists; choose another version or InstallRoot." }
$bin = Join-Path $InstallRoot "bin"
$launcher = Join-Path $bin "telia.cmd"
if ((Test-Path $launcher) -and (Get-Content $launcher -First 1) -ne '@rem Telia managed launcher') {
    throw "Refusing to replace an unrelated launcher: $launcher"
}
New-Item -ItemType Directory -Force $InstallRoot | Out-Null
$work = Join-Path $InstallRoot (".install-" + [guid]::NewGuid().ToString())
New-Item -ItemType Directory $work | Out-Null
try {
    $asset = "telia-windows-amd64.tar.gz"
    $base = "$repo/releases/download/$Version"
    Write-Host "Installing Telia $Version for windows-amd64..."
    Invoke-WebRequest "$base/$asset" -OutFile (Join-Path $work $asset)
    Invoke-WebRequest "$base/SHA256SUMS" -OutFile (Join-Path $work "SHA256SUMS")
    $entries = @(Get-Content (Join-Path $work "SHA256SUMS") | Where-Object { $_ -match "^[a-fA-F0-9]{64}  $([regex]::Escape($asset))$" })
    if ($entries.Count -ne 1) { throw "Missing or invalid checksum for $asset" }
    $expected = $entries[0].Substring(0, 64)
    $actual = (Get-FileHash (Join-Path $work $asset) -Algorithm SHA256).Hash
    if ($expected -ne $actual) { throw "Checksum mismatch; nothing was installed." }
    $bundle = Join-Path $work "bundle"
    New-Item -ItemType Directory $bundle | Out-Null
    & tar.exe -xzf (Join-Path $work $asset) -C $bundle
    if ($LASTEXITCODE -ne 0) { throw "Unable to extract release archive." }
    if (-not (Test-Path "$bundle/telia.exe") -or -not (Test-Path "$bundle/base/std")) { throw "Incomplete release bundle." }
    Move-Item $bundle $destination
    New-Item -ItemType Directory -Force $bin | Out-Null
    Set-Content -Encoding ASCII $launcher "@rem Telia managed launcher`r`n@echo off`r`n`"$destination\telia.exe`" %*`r`nexit /b %errorlevel%"
    Write-Host "Installed in $destination"
    Write-Host "Add this directory to your user PATH: $bin"
    Write-Host "Then run: telia version"
    if (-not ((Get-Command cc,gcc,clang -ErrorAction SilentlyContinue) | Select-Object -First 1)) {
        Write-Warning "A C compiler is required. Install MSYS2 UCRT64 GCC and add C:\msys64\ucrt64\bin to PATH."
    }
} finally {
    Remove-Item -Recurse -Force $work
}
