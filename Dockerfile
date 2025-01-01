# 使用golang官方镜像作为构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache gcc musl-dev curl

# 安装 protoc 插件
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 安装 buf
RUN curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v1.48.0/buf-$(uname -s)-$(uname -m)" \
    -o "/usr/local/bin/buf" && \
    chmod +x "/usr/local/bin/buf"

# 复制proto文件和buf配置
COPY proto/ proto/
COPY buf.* ./

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 确保 protoc 插件在 PATH 中
ENV PATH="$PATH:$(go env GOPATH)/bin"

# 生成protobuf代码
RUN buf generate

# 复制其余源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 使用轻量级的alpine作为运行阶段的基础镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/main .
COPY init.sql .

# 创建 secret.img
RUN dd if=/dev/urandom of=secret.img bs=512 count=1024

# 初始化数据库
RUN sqlite3 mbs.sqlite ".read init.sql"

# 暴露9090端口（从代码中可以看到默认端口是9090）
EXPOSE 9090

# 运行应用
CMD ["./main"] 