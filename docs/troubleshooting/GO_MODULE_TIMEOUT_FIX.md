# Go 模块下载超时问题解决方案

## 🔴 错误信息

```
go: cloud.google.com/go@v0.118.3: Get "https://proxy.golang.org/cloud.google.com/go/@v/v0.118.3.mod": dial tcp 142.250.69.[...]: timeout
```

这是 Go 模块下载时的网络超时问题。

## ✅ 解决方案

### 方案一：使用国内 Go 代理（推荐）

#### 在 Dockerfile 中设置

如果使用 Docker 构建，在 Dockerfile 中添加：

```dockerfile
# 设置 Go 代理（在 RUN go mod download 之前）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

RUN go mod download
```

#### 在构建命令中设置

在 Makefile 或构建脚本中：

```makefile
go mod download:
	export GOPROXY=https://goproxy.cn,direct && \
	export GOSUMDB=sum.golang.google.cn && \
	go mod download
```

#### 在 Render 构建命令中设置

在 Render 的 Build Command 中：

```bash
export GOPROXY=https://goproxy.cn,direct && export GOSUMDB=sum.golang.google.cn && go mod download && go build -o book-manage
```

### 方案二：使用其他代理

如果 goproxy.cn 也不行，可以尝试：

```bash
# 阿里云代理
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 七牛云代理
export GOPROXY=https://goproxy.qiniu.com,direct

# 官方代理（如果网络好）
export GOPROXY=https://proxy.golang.org,direct
```

### 方案三：禁用校验和验证（临时方案）

**注意**：不推荐用于生产环境，仅用于临时解决网络问题。

```bash
export GOSUMDB=off
go mod download
```

### 方案四：增加超时时间

```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOTIMEOUT=300s  # 增加超时时间到 5 分钟
go mod download
```

---

## 🔧 针对不同场景的配置

### 场景 1：Docker 构建

在 Dockerfile 中：

```dockerfile
FROM golang:1.21

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o book-manage
```

### 场景 2：本地构建

在终端中：

```bash
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
go mod download
go build -o book-manage
```

### 场景 3：Render 部署

在 Render 的 Build Command 中：

```bash
export GOPROXY=https://goproxy.cn,direct && export GOSUMDB=sum.golang.google.cn && go mod download && go build -o book-manage
```

### 场景 4：Makefile

在 Makefile 中：

```makefile
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@export GOPROXY=https://goproxy.cn,direct && \
	export GOSUMDB=sum.golang.google.cn && \
	go mod download

.PHONY: build
build: deps
	go build -o book-manage
```

---

## 📝 常用 Go 代理地址

| 代理 | 地址 | 说明 |
|------|------|------|
| 七牛云 | `https://goproxy.qiniu.com,direct` | 国内，速度快 |
| 阿里云 | `https://mirrors.aliyun.com/goproxy/,direct` | 国内，稳定 |
| goproxy.cn | `https://goproxy.cn,direct` | 国内，推荐 |
| 官方 | `https://proxy.golang.org,direct` | 国外，可能慢 |

---

## 🔍 验证配置

设置代理后，验证是否生效：

```bash
go env GOPROXY
go env GOSUMDB
```

应该看到你设置的代理地址。

---

## 🆘 如果还是超时

1. **检查网络连接**：
   ```bash
   ping goproxy.cn
   ```

2. **尝试直接下载**：
   ```bash
   export GOPROXY=direct
   go mod download
   ```
   （这会直接从源下载，可能更慢但更稳定）

3. **使用 VPN 或代理**：
   如果网络环境限制，考虑使用 VPN

4. **分步下载**：
   ```bash
   go mod download -x  # 显示详细日志
   ```

---

## 💡 最佳实践

1. **生产环境**：使用稳定的代理（如 goproxy.cn）
2. **开发环境**：可以使用官方代理或直接连接
3. **CI/CD**：在构建脚本中明确设置代理
4. **Docker**：在 Dockerfile 中设置环境变量

---

## 📚 相关文档

- [Go Modules 官方文档](https://go.dev/ref/mod)
- [goproxy.cn 文档](https://goproxy.cn/)
- [Go 代理配置](https://golang.org/cmd/go/#hdr-Module_proxy_protocol)



