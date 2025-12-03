# Copy files with proper UTF-8 encoding
param(
    [string]$SourceDir = "..\battle-tiles\internal\utils\plaza",
    [string]$DestDir = "internal\plaza"
)

Write-Host "Copying files from $SourceDir to $DestDir..." -ForegroundColor Cyan

$files = Get-ChildItem -Path $SourceDir -Filter "*.go"
$count = 0

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    
    # Replace imports
    $content = $content -replace 'battle-tiles/internal/utils/plaza', 'battle-bot/internal/plaza'
    $content = $content -replace 'battle-tiles/internal/dal/vo/game', 'battle-bot/internal/plaza/game'
    $content = $content -replace 'battle-tiles/internal/consts', 'battle-bot/internal/plaza/consts'
    
    $destFile = Join-Path $DestDir $file.Name
    [System.IO.File]::WriteAllText($destFile, $content, [System.Text.UTF8Encoding]::new($false))
    
    Write-Host "  OK: $($file.Name)" -ForegroundColor Green
    $count++
}

Write-Host "Copied $count files" -ForegroundColor Yellow
