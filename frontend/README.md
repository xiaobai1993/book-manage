# 图书管理系统前端

基于 Vue3 + Element Plus 开发的图书管理系统前端应用。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vue Router** - 官方路由管理器
- **Pinia** - 状态管理
- **Element Plus** - Vue 3 组件库
- **Axios** - HTTP 客户端
- **Vite** - 构建工具
- **Sass** - CSS 预处理器

## 功能特性

- ✅ 用户注册、登录、密码找回
- ✅ 个人信息管理
- ✅ 图书搜索和查看
- ✅ 图书借阅和归还
- ✅ 借阅记录查询
- ✅ 管理员功能（图书增删改、全量借阅记录）
- ✅ 响应式设计，支持移动端和PC端
- ✅ 路由权限控制

## 安装和运行

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

服务将在 `http://localhost:3000` 启动

### 3. 构建生产版本

```bash
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── api/          # API接口封装
│   │   ├── user.js
│   │   ├── book.js
│   │   └── borrow.js
│   ├── assets/       # 静态资源
│   ├── components/   # 组件
│   ├── layouts/      # 布局组件
│   │   └── MainLayout.vue
│   ├── router/       # 路由配置
│   │   └── index.js
│   ├── stores/       # 状态管理
│   │   └── user.js
│   ├── styles/       # 样式文件
│   │   └── main.scss
│   ├── utils/        # 工具函数
│   │   └── request.js
│   ├── views/        # 页面组件
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   ├── Books.vue
│   │   └── ...
│   ├── App.vue
│   └── main.js
├── index.html
├── package.json
├── vite.config.js
└── README.md
```

## 页面说明

### 用户相关
- **登录页面** (`/login`) - 用户登录
- **注册页面** (`/register`) - 用户注册
- **找回密码** (`/forget-password`) - 密码找回
- **个人信息** (`/profile`) - 查看和修改个人信息

### 图书相关
- **图书列表** (`/books`) - 图书搜索和列表展示
- **图书详情** (`/book/:id`) - 图书详细信息
- **添加图书** (`/book-add`) - 管理员添加图书
- **编辑图书** (`/book-edit/:id`) - 管理员编辑图书

### 借阅相关
- **我的借阅** (`/my-borrows`) - 个人借阅记录
- **全部借阅记录** (`/all-borrows`) - 管理员查看全部借阅记录

## 配置说明

### API代理配置

在 `vite.config.js` 中配置了API代理，开发环境会自动代理到后端服务：

```javascript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### 环境变量

可以创建 `.env` 文件配置环境变量：

```env
VITE_API_BASE_URL=http://localhost:8080
```

## 注意事项

1. 确保后端服务已启动并运行在 `http://localhost:8080`
2. 所有API请求都需要携带token（登录后自动添加）
3. 管理员功能需要管理员权限才能访问
4. 响应式设计支持移动端和PC端，使用Element Plus的响应式布局

## 开发说明

### 添加新页面

1. 在 `src/views/` 目录下创建新的Vue组件
2. 在 `src/router/index.js` 中添加路由配置
3. 如需权限控制，在路由的 `meta` 中添加 `requiresAuth` 或 `requiresAdmin`

### 添加新API

1. 在 `src/api/` 目录下创建或修改对应的API文件
2. 使用 `request` 工具函数发送请求
3. 在组件中导入并使用API函数

## 许可证

MIT License
