# 🚀 Go后端开发快速入门指南

## 🎯 立即开始（5分钟设置）

### 1. 环境检查


### 2. 项目启动


### 3. 验证安装
打开新终端，测试API：


如果看到JSON响应，恭喜！你已经成功启动了Go后端服务器！

---

## 📅 第一天学习计划（2-3小时）

### 上午：理解项目结构（1小时）
1. **阅读项目文档**（20分钟）
   - 打开 `README.md`，了解项目功能
   - 查看 `ARCHITECTURE.md`，理解架构设计

2. **探索项目结构**（20分钟）
   ```bash
   # 查看项目目录结构
   tree . -I 'node_modules|.git'
   
   # 或者使用ls命令
   ls -la
   ls -la internal/
   ls -la providers/
   ```

3. **理解Go模块系统**（20分钟）
   - 查看 `go.mod` 文件，了解项目依赖
   - 运行 `go list -m all` 查看所有依赖

### 下午：运行和测试（1-2小时）

1. **启动服务器并观察日志**（30分钟）
   ```bash
   go run cmd/main.go
   ```
   观察启动日志，理解：
   - 数据库初始化过程
   - 路由注册信息
   - 服务器启动信息

2. **测试所有API接口**（30分钟）
   ```bash
   # 基础接口
   curl http://localhost:8080/
   curl http://localhost:8080/api/health
   
   # 用户注册
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
   
   # 用户登录
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"password123"}'
   
   # 智能问答
   curl "http://localhost:8080/api/ask?prompt=Go语言有什么特点"
   
   # 查看记录
   curl http://localhost:8080/api/records
   ```

3. **查看数据库**（30分钟）
   ```bash
   # 安装sqlite3（如果没有）
   # macOS: brew install sqlite3
   # Ubuntu: sudo apt-get install sqlite3
   
   # 查看数据库
   sqlite3 qa_database.db
   .tables
   .schema users
   .schema qa_records
   SELECT * FROM users;
   SELECT * FROM qa_records LIMIT 5;
   .quit
   ```

---

## 📚 第一周学习计划

### 第1-2天：项目结构和基础概念
- [ ] 完成快速入门
- [ ] 阅读所有文档文件
- [ ] 理解Go模块和包管理
- [ ] 熟悉项目目录结构

### 第3-4天：HTTP服务器和路由
- [ ] 研究 `cmd/main.go` 文件
- [ ] 理解路由注册过程
- [ ] 学习 `internal/handlers/handlers.go`
- [ ] 测试所有API接口

### 第5-6天：数据库操作
- [ ] 学习 `internal/storage/` 目录下的文件
- [ ] 理解数据模型设计
- [ ] 练习SQL查询
- [ ] 了解数据库连接管理

### 第7天：认证系统
- [ ] 研究 `internal/auth/` 目录
- [ ] 理解JWT工作原理
- [ ] 测试用户注册登录流程
- [ ] 学习中间件概念

---

## 🔍 代码阅读顺序建议

### 第一轮：整体理解
1. `cmd/main.go` - 程序入口，了解整体流程
2. `internal/config/config.go` - 配置管理
3. `internal/handlers/handlers.go` - 核心业务逻辑

### 第二轮：深入模块
1. `internal/storage/models.go` - 数据模型
2. `internal/storage/qa_storage.go` - 数据库操作
3. `internal/auth/models.go` - 用户模型
4. `internal/auth/handlers.go` - 认证逻辑

### 第三轮：高级特性
1. `internal/auth/jwt.go` - JWT实现
2. `internal/auth/middleware.go` - 认证中间件
3. `internal/llm/client.go` - LLM集成
4. `providers/interface.go` - 接口设计

---

## 🛠️ 实践练习建议

### 初学者练习
1. **修改响应消息**
   - 在 `handlers.go` 中修改API响应消息
   - 重启服务器，测试变化

2. **添加新的API接口**
   - 在 `handlers.go` 中添加一个简单的 `/api/time` 接口
   - 返回当前时间的JSON响应

3. **修改数据库表结构**
   - 在 `models.go` 中为用户表添加新字段
   - 重新启动服务器，观察表结构变化

### 进阶练习
1. **实现用户资料更新功能**
2. **添加问答记录的删除功能**
3. **实现简单的用户权限控制**

---

## ❓ 常见问题解答

### Q: 启动时提示"端口被占用"怎么办？
A: 修改 `.env` 文件中的 `PORT=8080` 为其他端口，如 `PORT=8081`

### Q: 如何查看详细的错误日志？
A: 在启动服务器的终端中可以看到所有日志输出

### Q: 数据库文件在哪里？
A: 默认在项目根目录的 `qa_database.db` 文件

### Q: 如何重置数据库？
A: 删除 `qa_database.db` 文件，重启服务器会自动创建新的数据库

### Q: JWT Token在哪里查看？
A: 用户登录成功后，响应中的 `token` 字段就是JWT Token

### Q: 如何配置真实的LLM服务？
A: 修改 `.env` 文件中的 `LLM_PROVIDER` 和 `LLM_API_KEY`

---

## 📖 推荐学习资源

### Go语言基础
- [Go官方教程](https://tour.golang.org/welcome/1)
- [Go by Example](https://gobyexample.com/)

### Web开发
- [Gorilla Mux文档](https://github.com/gorilla/mux)
- [Go Web编程教程](https://github.com/astaxie/build-web-application-with-golang)

### 数据库
- [SQLite教程](https://www.runoob.com/sqlite/sqlite-tutorial.html)
- [Go数据库编程](https://go-database-sql.org/)

---

## 🎯 学习目标检查

完成第一周学习后，你应该能够：
- [ ] 独立启动和停止Go服务器
- [ ] 理解项目的基本架构
- [ ] 使用curl测试所有API接口
- [ ] 查看和理解数据库结构
- [ ] 修改简单的代码并看到效果
- [ ] 理解HTTP请求处理流程
- [ ] 了解JWT认证的基本原理

---

## 🚀 下一步

完成快速入门后，请阅读详细的 `GO_LEARNING_GUIDE.md` 文件，按照7个阶段系统地学习Go后端开发。

记住：学习编程最重要的是动手实践，不要只看不做！每个概念都要通过代码来验证和理解。

祝你学习愉快！🎉