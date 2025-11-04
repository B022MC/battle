# === 构建阶段 ===
FROM golang:1.23.0-bullseye AS builder

WORKDIR /go/src/app
COPY . .

# 设置代理 + 构建
RUN go env -w GOPROXY="https://goproxy.cn,direct" && \
    go mod tidy && \
    make build

# === 运行阶段 ===
FROM alpine:latest

WORKDIR /work

# 拷贝编译好的可执行文件和配置文件
COPY --from=builder /go/src/app/bin/ /work/
COPY --from=builder /go/src/app/configs/ /work/configs

# 添加可执行路径
ENV PATH="/work:${PATH}"

CMD ["main"]
