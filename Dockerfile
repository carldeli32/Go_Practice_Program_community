# ── 第一阶段：编译 Go ──
FROM golang:1.24-alpine AS go-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o community-server .

# ── 第二阶段：构建前端 ──
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ── 第三阶段：最小运行镜像 ──
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata

# Go 二进制（放 backend/ 子目录，保持代码中 ../frontend/dist 等相对路径有效）
WORKDIR /app/backend
COPY --from=go-builder /app/community-server .

# 前端构建产物
COPY --from=frontend-builder /app/dist /app/frontend/dist

# 上传目录
RUN mkdir -p /app/uploads/images /app/uploads/files

EXPOSE 8080
CMD ["./community-server"]
