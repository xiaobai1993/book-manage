.PHONY: help start-backend start-frontend start stop stop-backend stop-frontend clean install install-backend install-frontend check-ports kill-port-8080 kill-port-3000 status logs logs-backend logs-frontend dev-backend dev-frontend restart build-backend

# 变量定义
BACKEND_PORT := 8080
FRONTEND_PORT := 3000
BACKEND_DIR := .
FRONTEND_DIR := frontend
BACKEND_BINARY := ./book-manage
GO_CMD := go
NPM_CMD := npm

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help:
	@echo "图书管理系统 - Makefile 命令"
	@echo ""
	@echo "常用命令："
	@echo "  make start          - 启动所有服务（后端 + 前端）"
	@echo "  make start-backend  - 仅启动后端服务"
	@echo "  make start-frontend - 仅启动前端服务"
	@echo "  make stop           - 停止所有服务"
	@echo "  make stop-backend   - 停止后端服务"
	@echo "  make stop-frontend  - 停止前端服务"
	@echo "  make restart        - 重启所有服务"
	@echo "  make install        - 安装所有依赖"
	@echo "  make install-backend - 安装后端依赖"
	@echo "  make install-frontend - 安装前端依赖"
	@echo "  make build-backend  - 编译后端程序"
	@echo "  make clean          - 清理生成的文件"
	@echo "  make check-ports    - 检查端口占用情况"
	@echo ""

# 检查端口是否被占用
check-ports:
	@echo "检查端口占用情况..."
	@echo "检查后端端口 $(BACKEND_PORT)..."
	@lsof -ti:$(BACKEND_PORT) > /dev/null 2>&1 && echo "  ✓ 端口 $(BACKEND_PORT) 已被占用" || echo "  ✗ 端口 $(BACKEND_PORT) 未被占用"
	@echo "检查前端端口 $(FRONTEND_PORT)..."
	@lsof -ti:$(FRONTEND_PORT) > /dev/null 2>&1 && echo "  ✓ 端口 $(FRONTEND_PORT) 已被占用" || echo "  ✗ 端口 $(FRONTEND_PORT) 未被占用"

# 杀掉占用指定端口的进程
kill-port-8080:
	@echo "检查并清理端口 $(BACKEND_PORT)..."
	@if lsof -ti:$(BACKEND_PORT) > /dev/null 2>&1; then \
		echo "  发现端口 $(BACKEND_PORT) 被占用，正在清理..."; \
		lsof -ti:$(BACKEND_PORT) | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
		echo "  端口 $(BACKEND_PORT) 已清理"; \
	else \
		echo "  端口 $(BACKEND_PORT) 未被占用"; \
	fi

kill-port-3000:
	@echo "检查并清理端口 $(FRONTEND_PORT)..."
	@if lsof -ti:$(FRONTEND_PORT) > /dev/null 2>&1; then \
		echo "  发现端口 $(FRONTEND_PORT) 被占用，正在清理..."; \
		lsof -ti:$(FRONTEND_PORT) | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
		echo "  端口 $(FRONTEND_PORT) 已清理"; \
	else \
		echo "  端口 $(FRONTEND_PORT) 未被占用"; \
	fi

# 安装后端依赖
install-backend:
	@echo "安装后端依赖..."
	@cd $(BACKEND_DIR) && $(GO_CMD) mod download
	@echo "后端依赖安装完成"

# 安装前端依赖
install-frontend:
	@echo "安装前端依赖..."
	@cd $(FRONTEND_DIR) && $(NPM_CMD) install
	@echo "前端依赖安装完成"

# 安装所有依赖
install: install-backend install-frontend
	@echo "所有依赖安装完成"

# 编译后端
build-backend:
	@echo "编译后端程序..."
	@cd $(BACKEND_DIR) && $(GO_CMD) build -o $(BACKEND_BINARY) main.go
	@echo "后端程序编译完成"

# 停止后端服务
stop-backend: kill-port-8080
	@echo "后端服务已停止"

# 停止前端服务
stop-frontend: kill-port-3000
	@echo "前端服务已停止"

# 停止所有服务
stop: stop-backend stop-frontend
	@echo "所有服务已停止"

# 启动后端服务（后台运行）
start-backend: kill-port-8080
	@echo "启动后端服务..."
	@if [ ! -f $(BACKEND_BINARY) ]; then \
		echo "  后端程序不存在，正在编译..."; \
		$(MAKE) build-backend; \
	fi
	@cd $(BACKEND_DIR) && nohup $(BACKEND_BINARY) > backend.log 2>&1 & \
	echo "  后端进程已启动"; \
	sleep 2; \
	BACKEND_PID=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	if [ -n "$$BACKEND_PID" ]; then \
		echo "  ✓ 后端服务启动成功，端口: $(BACKEND_PORT)"; \
		echo "  日志文件: $(BACKEND_DIR)/backend.log"; \
		echo "  进程 PID: $$BACKEND_PID"; \
	else \
		echo "  ✗ 后端服务启动失败，请检查日志: $(BACKEND_DIR)/backend.log"; \
		exit 1; \
	fi

# 启动前端服务（后台运行）
start-frontend: kill-port-3000
	@echo "启动前端服务..."
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		echo "  前端依赖未安装，正在安装..."; \
		$(MAKE) install-frontend; \
	fi
	@cd $(FRONTEND_DIR) && nohup $(NPM_CMD) run dev > ../frontend.log 2>&1 & \
	echo "  前端进程已启动"; \
	sleep 3; \
	FRONTEND_PID=$$(lsof -ti:$(FRONTEND_PORT) 2>/dev/null); \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo "  ✓ 前端服务启动成功，端口: $(FRONTEND_PORT)"; \
		echo "  访问地址: http://localhost:$(FRONTEND_PORT)"; \
		echo "  日志文件: frontend.log"; \
		echo "  进程 PID: $$FRONTEND_PID"; \
	else \
		echo "  ✗ 前端服务启动失败，请检查日志: frontend.log"; \
		exit 1; \
	fi

# 启动所有服务
start: stop start-backend start-frontend
	@echo ""
	@echo "=========================================="
	@echo "  所有服务已启动"
	@echo "=========================================="
	@echo "  后端服务: http://localhost:$(BACKEND_PORT)"
	@echo "  前端服务: http://localhost:$(FRONTEND_PORT)"
	@echo ""
	@echo "  查看后端日志: tail -f backend.log"
	@echo "  查看前端日志: tail -f frontend.log"
	@echo ""
	@echo "  停止所有服务: make stop"
	@echo "=========================================="

# 重启所有服务
restart: stop start

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	@rm -f $(BACKEND_BINARY)
	@rm -f backend.log
	@rm -f frontend.log
	@rm -rf $(FRONTEND_DIR)/node_modules
	@rm -rf $(FRONTEND_DIR)/dist
	@echo "清理完成"

# 查看日志
logs-backend:
	@tail -f backend.log

logs-frontend:
	@tail -f frontend.log

logs: logs-backend

# 开发模式（前台运行，方便查看日志）
dev-backend: kill-port-8080
	@echo "启动后端服务（开发模式）..."
	@cd $(BACKEND_DIR) && $(GO_CMD) run main.go

dev-frontend: kill-port-3000
	@echo "启动前端服务（开发模式）..."
	@cd $(FRONTEND_DIR) && $(NPM_CMD) run dev

# 检查服务状态
status:
	@echo "检查服务状态..."
	@echo ""
	@echo "后端服务 (端口 $(BACKEND_PORT)):"
	@BACKEND_PID=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	if [ -n "$$BACKEND_PID" ]; then \
		echo "  ✓ 运行中 (PID: $$BACKEND_PID)"; \
	else \
		echo "  ✗ 未运行"; \
	fi
	@echo ""
	@echo "前端服务 (端口 $(FRONTEND_PORT)):"
	@FRONTEND_PID=$$(lsof -ti:$(FRONTEND_PORT) 2>/dev/null); \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo "  ✓ 运行中 (PID: $$FRONTEND_PID)"; \
	else \
		echo "  ✗ 未运行"; \
	fi
	@echo ""
