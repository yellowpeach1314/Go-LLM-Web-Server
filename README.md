# 🤖 现代化 LLM 问答系统 v2.0

这是一个完整的现代化后端应用，采用**模块化架构设计**和**用户认证系统**，展示了企业级Go应用的最佳实践。项目实现了关注点分离（Separation of Concerns）、依赖注入（Dependency Injection）和清洁架构（Clean Architecture）的设计原则。

## 🏗️ 项目架构

### 模块化设计

```
go-base-web-server/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序文件
├── internal/              # 内部模块（不对外暴露）
│   ├── auth/              # 🔐 认证模块
│   │   ├── handlers.go    # 认证处理器
│   │   ├── jwt.go         # JWT服务
│   │   ├── middleware.go  # 认证中间件
│   │   └── models.go      # 认证数据模型
│   ├── config/            # ⚙️ 配置模块
│   │   └── config.go      # 配置管理
│   ├── handlers/          # 🎯 处理器模块
│   │   └── handlers.go    # HTTP处理器
│   ├── llm/               # 🤖 LLM模块
│   │   └── client.go      # LLM客户端
│   ├── middleware/        # 🔧 中间件模块
│   │   └── middleware.go  # 通用中间件
│   └── storage/           # 💾 存储模块
│       ├── models.go      # 数据模型
│       ├── qa_storage.go  # QA记录存储
│       └── user_storage.go # 用户存储
├── providers/             # 🔌 LLM提供商
│   ├── interface.go       # 提供商接口
│   ├── openai.go         # OpenAI实现
│   ├── baidu.go          # 百度实现
│   ├── ali.go            # 阿里实现
│   ├── gemini.go         # Gemini实现
│   └── mock.go           # 模拟实现
├── .env.example          # 环境变量示例
├── go.mod               # Go模块文件
└── README.md           # 项目文档
```

### 核心组件关系

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   🎯 Handlers   │    │   🔐 Auth       │    │   💾 Storage    │
│                 │    │                 │    │                 │
│ • HTTP处理      │◄──►│ • JWT认证       │◄──►│ • 数据库操作     │
│ • 业务逻辑      │    │ • 用户管理       │    │ • 数据持久化     │
│ • 响应格式化     │    │ • 权限控制       │    │ • 事务管理       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   🤖 LLM        │
                    │                 │
                    │ • 多Provider支持 │
                    │ • 智能问答       │
                    │ • 连接管理       │
                    └─────────────────┘
```

## 🚀 新功能特性

### ✅ 用户认证系统
- **JWT认证**：基于Token的无状态认证
- **用户注册/登录**：完整的用户管理流程
- **权限控制**：支持匿名和认证用户访问
- **Token刷新**：自动续期机制

### ✅ 模块化架构
- **清洁架构**：按功能模块组织代码
- **依赖注入**：松耦合的组件设计
- **接口抽象**：易于测试和扩展
- **配置管理**：统一的配置处理

### ✅ 企业级特性
- **中间件支持**：日志、CORS、认证等
- **错误处理**：统一的错误响应格式
- **数据验证**：请求参数验证
- **安全加固**：密码加密、SQL注入防护

## 📋 API 接口

### 公开接口（无需认证）

| 方法 | 路径 | 描述 | 示例 |
|------|------|------|------|
| GET | `/` | API信息 | `curl http://localhost:8080/` |
| GET | `/api/health` | 健康检查 | `curl http://localhost:8080/api/health` |
| POST | `/api/auth/register` | 用户注册 | 见下方示例 |
| POST | `/api/auth/login` | 用户登录 | 见下方示例 |
| POST | `/api/auth/logout` | 用户登出 | `curl -X POST http://localhost:8080/api/auth/logout` |

### 可选认证接口（支持匿名访问）

| 方法 | 路径 | 描述 | 示例 |
|------|------|------|------|
| GET | `/api/ask` | 智能问答 | `curl "http://localhost:8080/api/ask?prompt=你好"` |
| GET | `/api/records` | 获取所有记录 | `curl http://localhost:8080/api/records` |
| GET | `/api/records/{id}` | 获取特定记录 | `curl http://localhost:8080/api/records/1` |

### 需要认证的接口

| 方法 | 路径 | 描述 | 示例 |
|------|------|------|------|
| GET | `/api/user/profile` | 获取用户资料 | 需要Bearer Token |
| POST | `/api/user/refresh-token` | 刷新Token | 需要Bearer Token |
| GET | `/api/user/records` | 获取用户记录 | 需要Bearer Token |
| GET | `/api/user/users` | 获取用户列表 | 需要Bearer Token |

## � 快速开始

### 1. 环境准备

确保已安装Go 1.21+：
```bash
go version
```

### 2. 克隆项目

```bash
git clone <repository-url>
cd go-base-web-server
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置环境

复制环境变量示例文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件：
```bash
# LLM配置
LLM_PROVIDER=openai
LLM_API_KEY=your_api_key_here
LLM_MODEL=gpt-3.5-turbo

# JWT配置
JWT_SECRET=your-super-secret-jwt-key-min-32-chars

# 数据库配置
DB_PATH=./qa_database.db

# 服务器配置
PORT=8080
```

### 5. 启动服务

```bash
go run cmd/main.go
```

### 6. 访问应用

打开浏览器访问：http://localhost:8080

## � 使用示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**响应：**
```json
{
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "api_key": "generated-api-key",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "status": "success"
}
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 认证请求

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/user/profile
```

### 智能问答（支持匿名）

```bash
# 匿名请求
curl "http://localhost:8080/api/ask?prompt=你好"

# 认证请求
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  "http://localhost:8080/api/ask?prompt=你好"
```

## 🔐 认证机制

### JWT Token格式

```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "exp": 1640995200,
  "iat": 1640908800,
  "iss": "go-base-web-server"
}
```

### 认证流程

1. **注册/登录** → 获取JWT Token
2. **请求头添加** → `Authorization: Bearer <token>`
3. **服务器验证** → 解析Token并验证用户
4. **上下文注入** → 将用户信息注入请求上下文

## 🛠️ 开发指南

### 添加新的API端点

1. 在 `internal/handlers/` 中添加处理函数
2. 在 `cmd/main.go` 中注册路由
3. 根据需要应用认证中间件

### 集成新的LLM Provider

1. 在 `providers/` 目录实现 `LLMProvider` 接口
2. 在 `internal/llm/client.go` 中添加Provider支持
3. 更新配置文件

### 扩展用户权限

1. 在 `internal/auth/models.go` 中扩展用户模型
2. 在 `internal/auth/middleware.go` 中添加权限检查
3. 更新数据库表结构

## 🧪 测试

### 功能测试

```bash
# 健康检查
curl http://localhost:8080/api/health

# 用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456"}'

# 智能问答
curl "http://localhost:8080/api/ask?prompt=测试问题"
```

### 压力测试

```bash
# 使用ab工具
ab -n 1000 -c 10 "http://localhost:8080/api/health"
```

## � 性能优化

### 数据库优化
- 连接池配置
- 索引优化
- 查询优化

### 缓存策略
- LLM响应缓存
- 用户会话缓存
- 静态资源缓存

### 并发处理
- Goroutine池
- 请求限流
- 超时控制

## 🚀 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 环境变量配置

```bash
# 生产环境
JWT_SECRET=your-production-secret-key
LLM_API_KEY=your-production-api-key
DB_PATH=/data/qa_database.db
PORT=8080
```

## 🔍 监控与日志

### 日志格式

```
2024/01/01 10:00:00 GET /api/ask 127.0.0.1
2024/01/01 10:00:01 收到问题: 你好
2024/01/01 10:00:01 认证用户提问: testuser (ID: 1)
2024/01/01 10:00:02 LLM响应成功，答案长度: 25
```

### 健康检查响应

```json
{
  "status": "ok",
  "timestamp": "2024-01-01 10:00:00",
  "services": {
    "database": "ok",
    "llm": "ok"
  },
  "llm_provider": {
    "provider": "OpenAI Provider",
    "status": "active"
  }
}
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## � 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP路由器
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT实现
- [Bcrypt](https://golang.org/x/crypto/bcrypt) - 密码加密
- [SQLite](https://www.sqlite.org/) - 数据库
- [Validator](https://github.com/go-playground/validator) - 数据验证

---

**项目状态**: ✅ 生产就绪  
**维护状态**: 🔄 积极维护  
**版本**: v2.0.0  
**最后更新**: 2024年1月