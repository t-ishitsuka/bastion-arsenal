# Arsenal installer for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/t-ishitsuka/bastion-arsenal/main/scripts/install.ps1 | iex

$ErrorActionPreference = 'Stop'

# Configuration
$Repo = "t-ishitsuka/bastion-arsenal"
$BinaryName = "bastion-arsenal"
$InstallDir = if ($env:ARSENAL_INSTALL_DIR) { $env:ARSENAL_INSTALL_DIR } else { "$env:USERPROFILE\.local\bin" }

# Functions
function Write-Info {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Err {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Write-Warn {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Err "サポートされていないアーキテクチャ: $arch"
            exit 1
        }
    }
}

function Get-LatestRelease {
    Write-Info "最新リリース情報を取得中..."

    $apiUrl = "https://api.github.com/repos/$Repo/releases/latest"

    try {
        $response = Invoke-RestMethod -Uri $apiUrl -Method Get
        $version = $response.tag_name

        Write-Success "最新バージョン: $version"
        return $version
    }
    catch {
        Write-Err "リリース情報の取得に失敗しました: $_"
        exit 1
    }
}

function Download-Archive {
    param(
        [string]$Version,
        [string]$Arch
    )

    $archiveName = "$BinaryName-$Version-windows-$Arch.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$archiveName"

    Write-Info "ダウンロード中: $archiveName"

    $tmpDir = Join-Path $env:TEMP "arsenal-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    $archivePath = Join-Path $tmpDir $archiveName

    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing
        Write-Success "ダウンロード完了"
        return @{
            ArchivePath = $archivePath
            TmpDir = $tmpDir
        }
    }
    catch {
        Write-Err "ダウンロードに失敗しました: $_"
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }
}

function Extract-AndInstall {
    param(
        [string]$ArchivePath,
        [string]$TmpDir,
        [string]$Arch
    )

    Write-Info "展開中..."

    try {
        Expand-Archive -Path $ArchivePath -DestinationPath $TmpDir -Force

        $binaryFile = Join-Path $TmpDir "$BinaryName-windows-$Arch.exe"

        if (-not (Test-Path $binaryFile)) {
            Write-Err "バイナリファイルが見つかりません: $binaryFile"
            exit 1
        }

        # Create install directory if it doesn't exist
        if (-not (Test-Path $InstallDir)) {
            Write-Info "インストールディレクトリを作成中: $InstallDir"
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        $targetPath = Join-Path $InstallDir "$BinaryName.exe"

        # Backup existing binary if it exists
        if (Test-Path $targetPath) {
            Write-Info "既存のバイナリをバックアップ中..."
            $backupPath = Join-Path $InstallDir "$BinaryName.exe.backup"
            Move-Item -Path $targetPath -Destination $backupPath -Force
        }

        Write-Info "インストール中: $targetPath"
        Copy-Item -Path $binaryFile -Destination $targetPath -Force

        Write-Success "インストールが完了しました"
    }
    catch {
        Write-Err "インストールに失敗しました: $_"
        exit 1
    }
    finally {
        # Cleanup
        Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($currentPath -notlike "*$InstallDir*") {
        Write-Info "PATH に追加中: $InstallDir"

        $newPath = "$InstallDir;$currentPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")

        # Update current session PATH
        $env:Path = "$InstallDir;$env:Path"

        Write-Success "PATH に追加しました"
        Write-Warn "新しいターミナルセッションで有効になります"
    }
    else {
        Write-Success "PATH には既に追加されています"
    }
}

function Verify-Installation {
    Write-Info "インストールを確認中..."

    $targetPath = Join-Path $InstallDir "$BinaryName.exe"

    if (-not (Test-Path $targetPath)) {
        Write-Err "インストールに失敗しました: $targetPath が見つかりません"
        exit 1
    }

    # Refresh PATH in current session
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "User") + ";" + [Environment]::GetEnvironmentVariable("Path", "Machine")

    try {
        $version = & $targetPath version 2>&1 | Select-Object -First 1
        Write-Success "インストール済みバージョン: $version"
    }
    catch {
        Write-Warn "バージョン確認に失敗しましたが、インストールは完了しています"
    }
}

function Print-NextSteps {
    Write-Host ""
    Write-Host "次のステップ:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  1. 新しいターミナルを開くか、現在のセッションを再起動"
    Write-Host ""
    Write-Host "  2. シェルの設定に Arsenal を追加:"
    Write-Host "     & $BinaryName init-shell powershell | Invoke-Expression"
    Write-Host ""
    Write-Host "  3. ツールをインストール:"
    Write-Host "     $BinaryName install node 20.10.0"
    Write-Host ""
    Write-Host "  4. .toolversions から同期:"
    Write-Host "     $BinaryName sync"
    Write-Host ""
}

# Main
Write-Info "Arsenal インストーラー"
Write-Host ""

$arch = Get-Architecture
Write-Info "検出されたアーキテクチャ: $arch"

$version = Get-LatestRelease
$download = Download-Archive -Version $version -Arch $arch
Extract-AndInstall -ArchivePath $download.ArchivePath -TmpDir $download.TmpDir -Arch $arch
Add-ToPath
Verify-Installation
Print-NextSteps

Write-Success "インストールが完了しました！"
