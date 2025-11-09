# Supabase 获取数据库连接字符串详细步骤

## 📍 方法一：通过 Settings → Database（推荐）

### 步骤 1：进入项目设置
1. 登录 [Supabase Dashboard](https://app.supabase.com)
2. 选择你的项目（book-manage 或你创建的项目名）
3. 点击左侧菜单的 **Settings**（设置图标，齿轮⚙️）

### 步骤 2：进入 Database 设置
1. 在 Settings 页面，点击左侧子菜单中的 **Database**
2. 你会看到 "Database Settings" 页面

### 步骤 3：找到 Connection string
在 "Database Settings" 页面中，向下滚动，你会看到：

**Connection string** 部分，包含几个标签页：
- **URI** ← 选择这个！
- **JDBC**
- **Golang**
- **Python**
- **Node.js**

### 步骤 4：复制 URI 连接字符串
1. 点击 **URI** 标签页
2. 你会看到一个连接字符串，格式类似：
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres
   ```
3. 点击连接字符串右侧的 **复制按钮**（📋 图标）
4. **重要**：这个字符串中的 `[YOUR-PASSWORD]` 需要替换为你创建项目时设置的数据库密码

### 步骤 5：替换密码
连接字符串中的 `[YOUR-PASSWORD]` 需要替换为实际密码：

**示例**：
- 原始：`postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres`
- 替换后：`postgresql://postgres:你的实际密码@db.xxx.supabase.co:5432/postgres`

---

## 📍 方法二：通过 Connection Pooling（如果方法一找不到）

### 步骤 1：进入 Database 设置
1. Settings → Database
2. 找到 **Connection Pooling** 部分

### 步骤 2：使用 Session 模式
1. 在 Connection Pooling 中，选择 **Session** 模式
2. 复制 **Connection string**（格式类似方法一）

---

## 📍 方法三：手动构建连接字符串

如果以上方法都找不到，你可以手动构建：

### 步骤 1：获取数据库信息
在 Settings → Database 页面，找到：
- **Host**: `db.xxx.supabase.co`（类似格式）
- **Port**: `5432`
- **Database name**: `postgres`
- **User**: `postgres`
- **Password**: 你创建项目时设置的密码

### 步骤 2：构建连接字符串
格式：
```
postgresql://postgres:你的密码@db.xxx.supabase.co:5432/postgres
```

**示例**：
```
postgresql://postgres:mypassword123@db.abcdefghijk.supabase.co:5432/postgres
```

---

## 🔍 如果还是找不到？

### 检查点：
1. ✅ 确认你已经创建了 Supabase 项目
2. ✅ 确认项目已经完成初始化（等待 2-3 分钟）
3. ✅ 确认你在正确的项目页面
4. ✅ 尝试刷新页面（F5 或 Cmd+R）

### 截图位置参考：
- **Settings** 在左侧菜单最下方（齿轮图标）
- **Database** 在 Settings 的子菜单中
- **Connection string** 在 Database 页面中间位置

---

## 📝 使用连接字符串

获取到连接字符串后，在 Render 的环境变量中设置：

```
DATABASE_URL=postgresql://postgres:你的密码@db.xxx.supabase.co:5432/postgres
```

**注意**：
- 不要包含方括号 `[]`
- 密码中的特殊字符可能需要 URL 编码
- 确保整个字符串在一行内，没有换行

---

## 🆘 仍然有问题？

如果按照以上步骤还是找不到，请告诉我：
1. 你在 Supabase Dashboard 的哪个页面？
2. Settings 菜单中能看到哪些选项？
3. 能否截图或描述你看到的界面？

我可以根据你的具体情况提供更精确的指导。

