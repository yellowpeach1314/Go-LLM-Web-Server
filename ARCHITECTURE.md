# 项目架构文档

## 🏗️ 模块化重构说明

本项目已从单体架构重构为模块化架构，采用清洁架构原则，提高了代码的可维护性和可扩展性。

## 📁 新项目结构

```
go-base-web-server/
├── cmd/                    # 应用程序入口
│   └── main.go            # 新的主程序文件
├── internal/              # 内部模块（Go 1.4+ 约定）
│   ├── auth/              # 🔐 认证模块
│   │   ├── handlers.go    # 用户注册、登录、权限管理
│   │   ├── jwt.go         # JWT token生成和验证
│   │   ├── middleware.go  # 认证中间件
│   │   └── models.go      # 认证相关数据模型
│   ├── config/            # ⚙️ 配置模块
│   │   └── config.go      # 统一配置管理
│   ├── handlers/          # 🎯 HTTP处理器模块
│   │   └── handlers.go    # 业务逻辑处理器
│   ├── llm/               # 🤖 LLM客户端模块
│   │   └── client.go      # LLM服务抽象层
│   ├── middleware/        # 🔧 通用中间件模块
│   │   └── middleware.go  # 日志、CORS等中间件
│   └── storage/           # 💾 数据存储模块
│       ├── models.go      # 数据模型定义
│       ├── qa_storage.go  # QA记录存储逻辑
│       └── user_storage.go # 用户数据存储逻辑
├── providers/             # 🔌 LLM提供商（保持不变）
├── 旧文件/                # 原始文件（可删除）
│   ├── main.go           # 旧主文件
│   ├── handlers.go       # 旧处理器
│   ├── storage.go        # 旧存储
│   ├── auth_*.go         # 旧认证文件
│   └── ...
└── 配置文件
    ├── .env.example      # 更新的环境变量示例
    ├── go.mod           # Go模块文件
    └── README.md        # 更新的项目文档
```



## 🔄 重构对比

### 重构前（单体架构）
- 所有功能在根目录的单个文件中
- 紧耦合的组件设计
- 难以测试和维护
- 缺乏用户认证系统

### 重构后（模块化架构）
- 按功能模块组织代码
- 松耦合的接口设计
- 易于测试和扩展
- 完整的用户认证系统

## 🚀 如何运行新版本

### 方法1：使用新的主文件
```bash
go run cmd/main.go
```

### 方法2：构建可执行文件
```bash
go build -o app cmd/main.go
./app
```


## 🔧 配置说明

新版本需要配置JWT密钥：

```bash
# 复制环境变量文件
cp .env.example .env

# 编辑配置
vim .env
```

必需配置：
- `JWT_SECRET`: JWT签名密钥（生产环境必须设置）
- `LLM_API_KEY`: LLM服务API密钥（可选，无则使用模拟模式）

## 🆕 新增功能

1. **用户认证系统**
   - 用户注册和登录
   - JWT token认证
   - 权限控制

2. **模块化架构**
   - 清洁架构设计
   - 依赖注入
   - 接口抽象

3. **企业级特性**
   - 统一配置管理
   - 结构化错误处理
   - 请求参数验证

## 🔄 迁移指南

如果你有基于旧版本的代码：

1. **API兼容性**: 大部分API保持兼容，新增认证相关接口
2. **数据库**: 自动添加用户表，现有数据不受影响
3. **配置**: 需要添加JWT_SECRET配置项

## 🧪 测试新功能

```bash
# 1. 启动服务
go run cmd/main.go

# 2. 测试用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456"}'

# 3. 测试智能问答（匿名）
curl "http://localhost:8080/api/ask?prompt=你好"

# 4. 测试认证问答
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/ask?prompt=你好"
```



## 📚 开发建议

1. **添加新功能**: 在对应的模块目录中添加
2. **修改配置**: 在 `internal/config/` 中统一管理
3. **添加中间件**: 在 `internal/middleware/` 中实现
4. **扩展认证**: 在 `internal/auth/` 中添加新的认证逻辑

## 🔮 未来规划

- [ ] 添加Redis缓存支持
- [ ] 实现API限流
- [ ] 添加Prometheus监控
- [ ] 支持数据库迁移
- [ ] 实现角色权限系统
- [ ] 添加WebSocket支持