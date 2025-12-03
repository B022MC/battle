# Battle Bot Makefile

.PHONY: help build run clean test install dev

# 默认目标
help:
	@echo "Battle Bot - 四川游戏家园机器人"
	@echo ""
	@echo "可用命令:"
	@echo "  make build    - 编译程序"
	@echo "  make run      - 运行程序"
	@echo "  make dev      - 开发模式运行（自动重载）"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理编译文件"
	@echo "  make install  - 安装依赖"
	@echo "  make setup    - 初始化项目（复制plaza代码）"

# 编译
build:
	@echo "编译 battle-bot..."
	go build -o battle-bot.exe cmd/bot/main.go
	@echo "编译完成: battle-bot.exe"

# 运行
run:
	@echo "启动 battle-bot..."
	go run cmd/bot/main.go

# 开发模式
dev:
	@echo "开发模式运行..."
	go run cmd/bot/main.go -config=config.yaml

# 测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理
clean:
	@echo "清理编译文件..."
	@if exist battle-bot.exe del battle-bot.exe
	@if exist *.log del *.log
	@echo "清理完成"

# 安装依赖
install:
	@echo "安装依赖..."
	go mod download
	go mod tidy
	@echo "依赖安装完成"

# 初始化项目（复制plaza代码）
setup:
	@echo "初始化项目..."
	@if not exist config.yaml (copy config.yaml.example config.yaml)
	@echo "请编辑 config.yaml 填写你的配置"
	@echo ""
	@echo "接下来需要手动复制plaza协议代码:"
	@echo "1. 从 battle-tiles/internal/utils/plaza/ 复制文件到 internal/plaza/"
	@echo "2. 从 battle-tiles/internal/dal/vo/game/ 复制文件到 internal/plaza/game/"
	@echo "3. 从 battle-tiles/internal/consts/ 复制文件到 internal/plaza/consts/"
	@echo ""
	@echo "或运行: make copy-plaza"

# 复制plaza代码（Windows）
copy-plaza:
	@echo "复制plaza协议代码..."
	@if not exist internal\plaza\game mkdir internal\plaza\game
	@if not exist internal\plaza\consts mkdir internal\plaza\consts
	xcopy /Y ..\battle-tiles\internal\utils\plaza\*.go internal\plaza\
	xcopy /Y ..\battle-tiles\internal\dal\vo\game\*.go internal\plaza\game\
	xcopy /Y ..\battle-tiles\internal\consts\*.go internal\plaza\consts\
	@echo "复制完成，请修改import路径"

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	go vet ./...
