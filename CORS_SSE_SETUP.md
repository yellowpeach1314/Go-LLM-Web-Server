# CORS和SSE配置说明

## 问题描述
前端运行在 `http://localhost:5173`，后端运行在 `http://localhost:8080`，由于端口不同导致跨域请求被浏览器阻止，特别是SSE（Server-Sent Events）请求。

## 解决方案

### 1. 更新CORS中间件
- 文件：`internal/middleware/middleware.go`
- 支持多个前端域名（包括localhost:5173）
- 添加完整的CORS头部支持
- 支持预检请求处理

### 2. 优化SSE处理器
- 文件：`internal/handlers/handlers.go`
- 在SSE响应中添加专门的CORS头部
- 支持动态Origin处理
- 添加nginx缓冲禁用

### 3. 更新测试页面
- 文件：`stream_chat_demo.html`
- 添加CORS测试功能
- 改进错误处理和状态显示
- 支持SSE连接状态监控

## 测试步骤

### 1. 启动后端服务器


### 2. 测试CORS配置


### 3. 使用HTML测试页面
1. 打开 `stream_chat_demo.html`
2. 点击"测试CORS"按钮
3. 点击"测试连接"按钮
4. 尝试发送流式聊天请求

### 4. 前端集成测试
在前端项目中使用以下代码测试：



## 支持的前端域名
- `http://localhost:3000`
- `http://localhost:5173` (Vite默认)
- `http://localhost:8080`
- `http://127.0.0.1:3000`
- `http://127.0.0.1:5173`
- `http://127.0.0.1:8080`

## 故障排除

### 1. 仍然有CORS错误
- 检查前端运行的确切端口
- 确认Origin头部是否正确发送
- 查看浏览器开发者工具的Network标签

### 2. SSE连接失败
- 检查浏览器是否支持EventSource
- 确认网络连接稳定
- 查看服务器日志中的错误信息

### 3. 预检请求失败
- 确认OPTIONS请求返回200状态码
- 检查Access-Control-Allow-Methods头部
- 验证Access-Control-Allow-Headers包含所需头部

## 配置详情

### CORS头部说明
- `Access-Control-Allow-Origin`: 允许的源域名
- `Access-Control-Allow-Methods`: 允许的HTTP方法
- `Access-Control-Allow-Headers`: 允许的请求头部
- `Access-Control-Allow-Credentials`: 是否允许发送凭据
- `Access-Control-Expose-Headers`: 暴露给前端的响应头部
- `Access-Control-Max-Age`: 预检请求缓存时间

### SSE特殊头部
- `Content-Type: text/event-stream`: SSE内容类型
- `Cache-Control: no-cache`: 禁用缓存
- `Connection: keep-alive`: 保持连接
- `X-Accel-Buffering: no`: 禁用nginx缓冲