# Battle Bot 初始化脚本
Write-Host "🤖 Battle Bot 初始化脚本" -ForegroundColor Cyan
Write-Host ""

# 1. 检查配置文件
if (-not (Test-Path "config.yaml")) {
    Write-Host "📋 创建配置文件..." -ForegroundColor Yellow
    Copy-Item "config.yaml.example" "config.yaml"
    Write-Host "✅ 配置文件已创建: config.yaml" -ForegroundColor Green
    Write-Host "⚠️  请编辑 config.yaml 填写你的账号信息" -ForegroundColor Yellow
} else {
    Write-Host "✅ 配置文件已存在" -ForegroundColor Green
}
Write-Host ""

# 2. 创建必要的目录
Write-Host "📁 创建目录结构..." -ForegroundColor Yellow
$dirs = @(
    "internal\plaza\game",
    "internal\plaza\consts",
    "logs"
)

foreach ($dir in $dirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "  ✓ 创建: $dir" -ForegroundColor Gray
    }
}
Write-Host "✅ 目录结构创建完成" -ForegroundColor Green
Write-Host ""

# 3. 复制plaza协议代码
$sourcePath = "..\battle-tiles\internal\utils\plaza"
if (Test-Path $sourcePath) {
    Write-Host "📦 复制plaza协议代码..." -ForegroundColor Yellow
    
    # 复制plaza核心文件
    $plazaFiles = Get-ChildItem "$sourcePath\*.go" -ErrorAction SilentlyContinue
    if ($plazaFiles) {
        foreach ($file in $plazaFiles) {
            Copy-Item $file.FullName "internal\plaza\" -Force
            Write-Host "  ✓ 复制: $($file.Name)" -ForegroundColor Gray
        }
    }
    
    # 复制game相关文件
    $gamePath = "..\battle-tiles\internal\dal\vo\game"
    if (Test-Path $gamePath) {
        $gameFiles = Get-ChildItem "$gamePath\*.go" -ErrorAction SilentlyContinue
        if ($gameFiles) {
            foreach ($file in $gameFiles) {
                Copy-Item $file.FullName "internal\plaza\game\" -Force
                Write-Host "  ✓ 复制: game\$($file.Name)" -ForegroundColor Gray
            }
        }
    }
    
    # 复制consts文件
    $constsPath = "..\battle-tiles\internal\consts"
    if (Test-Path $constsPath) {
        $constsFiles = Get-ChildItem "$constsPath\*.go" -ErrorAction SilentlyContinue
        if ($constsFiles) {
            foreach ($file in $constsFiles) {
                Copy-Item $file.FullName "internal\plaza\consts\" -Force
                Write-Host "  ✓ 复制: consts\$($file.Name)" -ForegroundColor Gray
            }
        }
    }
    
    Write-Host "✅ 协议代码复制完成" -ForegroundColor Green
    Write-Host ""
    Write-Host "⚠️  注意: 需要手动修改import路径" -ForegroundColor Yellow
    Write-Host "   将 'battle-tiles/internal/...' 改为 'battle-bot/internal/plaza/...'" -ForegroundColor Gray
} else {
    Write-Host "⚠️  未找到 battle-tiles 项目，跳过协议代码复制" -ForegroundColor Yellow
    Write-Host "   请确保 battle-tiles 和 battle-bot 在同一目录下" -ForegroundColor Gray
}
Write-Host ""

# 4. 安装Go依赖
Write-Host "📦 安装Go依赖..." -ForegroundColor Yellow
go mod download
go mod tidy
Write-Host "✅ 依赖安装完成" -ForegroundColor Green
Write-Host ""

# 5. 完成提示
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "🎉 初始化完成！" -ForegroundColor Green
Write-Host ""
Write-Host "下一步:" -ForegroundColor Cyan
Write-Host "  1. 编辑 config.yaml 填写你的账号信息" -ForegroundColor White
Write-Host "  2. 如果复制了plaza代码，需要修改import路径" -ForegroundColor White
Write-Host "  3. 运行机器人:" -ForegroundColor White
Write-Host "     go run cmd/bot/main.go" -ForegroundColor Gray
Write-Host "     或" -ForegroundColor Gray
Write-Host "     make build && ./battle-bot.exe" -ForegroundColor Gray
Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
