# 个人博客系统后端

基于 Go 语言、Gin 框架和 GORM 库开发的个人博客系统后端 API。

## 功能说明

本项目实现了以下功能：

1. **用户认证**：用户注册、登录，JWT 认证
2. **文章管理**：创建、读取、更新、删除文章（CRUD）
3. **评论功能**：创建评论、获取文章评论列表
4. **权限控制**：只有文章作者可以修改/删除自己的文章
5. **错误处理**：统一的错误处理和日志记录

## 项目结构

```
gin-quickstart/
├── config/
│   └── database.go        # 数据库配置和初始化
├── models/
│   └── models.go          # 数据模型定义 (User, Post, Comment)
├── middleware/
│   ├── auth.go            # JWT 认证中间件
│   └── logger.go          # 日志记录中间件
├── handlers/
│   ├── auth.go            # 用户认证处理 (注册/登录)
│   ├── post.go            # 文章管理处理 (增删改查)
│   └── comment.go         # 评论功能处理
├── utils/
│   └── response.go        # 统一响应格式工具
├── main.go                # 主入口文件
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖版本锁定
└── README.md              # 项目说明文档
```

## 数据库设计

### users 表
- id：主键
- username：用户名（唯一）
- password：密码（bcrypt 加密）
- email：邮箱（唯一）
- created_at、updated_at、deleted_at：时间戳

### posts 表
- id：主键
- title：文章标题
- content：文章内容
- user_id：作者 ID（外键关联 users 表）
- created_at、updated_at、deleted_at：时间戳

### comments 表
- id：主键
- content：评论内容
- user_id：评论用户 ID（外键关联 users 表）
- post_id：文章 ID（外键关联 posts 表）
- created_at、updated_at、deleted_at：时间戳

## 运行环境

- Go 1.20 或更高版本
- MySQL 5.7 或更高版本

## 安装步骤

1. **克隆或下载项目**

2. **进入项目目录**
   ```bash
   cd gin-quickstart
   ```

3. **配置MySQL数据库**

   创建数据库：
   ```sql
   CREATE DATABASE grom CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

   修改 `config/database.go` 中的数据库连接信息（如果需要）：
   ```go
   dsn := "root:123456@tcp(localhost:3306)/grom?charset=utf8mb4&parseTime=True&loc=Local"
   ```

4. **安装依赖**
   ```bash
   go mod tidy
   ```

5. **运行项目**
   ```bash
   go run main.go
   ```

   服务器将在 `http://localhost:8080` 启动

## API 接口说明

### 公开接口（无需认证）

#### 1. 用户注册
```
POST /register
Content-Type: application/json

{
  "Username": "testuser",
  "Password": "password123",
  "Email": "test@example.com"
}
```

#### 2. 用户登录
```
POST /login
Content-Type: application/json

{
  "Username": "testuser",
  "Password": "password123"
}

响应：
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": 1
}
```

#### 3. 获取所有文章
```
GET /posts
```

#### 4. 获取单篇文章
```
GET /posts/:id
```

#### 5. 获取文章评论
```
GET /comments/post/:post_id
```

### 需要认证的接口（需要 JWT Token）

所有需要认证的接口都需要在请求头中添加：
```
Authorization: Bearer {token}
```

#### 6. 创建文章
```
POST /create-post
Authorization: Bearer {token}
Content-Type: application/json

{
  "Title": "文章标题",
  "Content": "文章内容"
}
```

#### 7. 更新文章（仅作者）
```
PUT /upposts/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "Title": "新标题",
  "Content": "新内容"
}
```

#### 8. 删除文章（仅作者）
```
DELETE /delposts/:id
Authorization: Bearer {token}
```

#### 9. 创建评论
```
POST /comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "Content": "评论内容",
  "PostID": 1
}
```

## 测试用例

### 使用 cURL 测试

#### 1. 注册用户
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"Username":"testuser","Password":"password123","Email":"test@example.com"}'
```

#### 2. 登录获取 Token
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"Username":"testuser","Password":"password123"}'
```

保存返回的 token 供后续使用。

#### 3. 创建文章
```bash
curl -X POST http://localhost:8080/create-post \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"Title":"我的第一篇博客","Content":"这是文章内容"}'
```

#### 4. 获取所有文章
```bash
curl http://localhost:8080/posts
```

#### 5. 获取单篇文章
```bash
curl http://localhost:8080/posts/1
```

#### 6. 更新文章
```bash
curl -X PUT http://localhost:8080/upposts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"Title":"更新的标题","Content":"更新的内容"}'
```

#### 7. 创建评论
```bash
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"Content":"很好的文章！","PostID":1}'
```

#### 8. 获取文章评论
```bash
curl http://localhost:8080/comments/post/1
```

#### 9. 删除文章
```bash
curl -X DELETE http://localhost:8080/delposts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 使用 Postman 测试

1. **设置环境变量**
   - `base_url`: `http://localhost:8080`
   - `token`: 登录后获取的 JWT token

2. **测试流程**
   - 注册用户 → 登录获取 token
   - 使用 token 创建文章
   - 查看文章列表和详情
   - 更新和删除文章
   - 创建和查看评论

## 错误处理

系统返回标准的 HTTP 状态码：

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或 token 无效
- `403 Forbidden`: 无权限操作
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

错误响应格式：
```json
{
  "error": "错误信息描述"
}
```

## 技术实现

- **密码加密**：使用 bcrypt 加密存储用户密码
- **JWT 认证**：使用 JWT 实现用户认证，token 有效期 24 小时
- **ORM**：使用 GORM 操作数据库，自动处理 SQL 注入
- **日志记录**：使用 log 包记录系统运行信息
- **数据库**：MySQL（支持生产环境）
- **模块化设计**：项目采用分层架构
  - **config**: 数据库配置和初始化
  - **models**: 数据模型定义
  - **middleware**: 中间件（认证、日志）
  - **handlers**: 业务逻辑处理
  - **utils**: 工具函数（响应格式化等）

## 项目文件

```
gin-quickstart/
├── config/
│   └── database.go        # 数据库配置和初始化
├── models/
│   └── models.go          # 数据模型定义 (User, Post, Comment)
├── middleware/
│   ├── auth.go            # JWT 认证中间件
│   └── logger.go          # 日志记录中间件
├── handlers/
│   ├── auth.go            # 用户认证处理 (注册/登录)
│   ├── post.go            # 文章管理处理 (增删改查)
│   └── comment.go         # 评论功能处理
├── utils/
│   └── response.go        # 统一响应格式工具
├── main.go                # 主入口文件
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖版本锁定
└── README.md              # 项目说明文档
```

## 依赖库

```
github.com/gin-gonic/gin          # Web 框架
gorm.io/gorm                      # ORM 库
gorm.io/driver/mysql              # MySQL 驱动
github.com/golang-jwt/jwt/v5      # JWT 认证
golang.org/x/crypto/bcrypt        # 密码加密
```

## 注意事项

1. JWT 密钥 `SECRET_KEY` 在生产环境中应修改为强密码并使用环境变量
2. 确保MySQL服务正在运行，并且已创建 `grom` 数据库
3. 根据实际情况修改数据库连接信息（用户名、密码、主机、端口）
4. Token 过期后需要重新登录获取新 token
5. GORM 会自动创建所需的表结构（users、posts、comments）

## 项目内容

✅ 项目初始化（go mod init）
✅ 数据库设计与模型定义（User、Post、Comment）
✅ 用户注册和登录（密码加密、JWT 认证）
✅ 文章管理功能（CRUD 操作）
✅ 评论功能（创建和读取）
✅ 权限控制（作者才能修改/删除文章）
✅ 错误处理和日志记录
✅ 项目模块化重构（分层架构）
✅ README 文档和测试用例

## 项目特色

1. **模块化设计**：采用分层架构，代码组织清晰，易于维护和扩展
2. **统一响应格式**：所有 API 返回统一的 JSON 格式
3. **完善的错误处理**：统一的错误处理机制和日志记录
4. **安全性**：密码 bcrypt 加密、JWT 认证、权限控制
5. **可扩展性**：模块化设计便于添加新功能