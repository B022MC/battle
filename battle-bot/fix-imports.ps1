# Fix Import Paths
# 批量替换所有Go文件的import路径

Write-Host "🔧 开始修复import路径..." -ForegroundColor Cyan
Write-Host ""

$replacements = @{
    "battle-tiles/internal/utils/plaza" = "battle-bot/internal/plaza"
    "battle-tiles/internal/dal/vo/game" = "battle-bot/internal/plaza/game"
    "battle-tiles/internal/consts" = "battle-bot/internal/plaza/consts"
}

$files = Get-ChildItem -Path "internal\plaza" -Recurse -Filter "*.go"
$totalFiles = $files.Count
$modifiedFiles = 0

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    $originalContent = $content
    $modified = $false
    
    foreach ($old in $replacements.Keys) {
        $new = $replacements[$old]
        if ($content -match [regex]::Escape($old)) {
            $content = $content -replace [regex]::Escape($old), $new
            $modified = $true
        }
    }
    
    if ($modified) {
        Set-Content $file.FullName -Value $content -Encoding UTF8 -NoNewline
        Write-Host "  ✓ $($file.Name)" -ForegroundColor Green
        $modifiedFiles++
    }
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "✅ 完成！" -ForegroundColor Green
Write-Host "   检查的文件: $totalFiles" -ForegroundColor White
Write-Host "   修改的文件: $modifiedFiles" -ForegroundColor White
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host ""
Write-Host "Next step: go mod tidy" -ForegroundColor Yellow
