# 大模型API接入指南

本项目支持多种主流大模型API，您可以根据需要选择合适的提供商。

## 支持的大模型提供商

### 1. OpenAI (GPT系列)
- **Provider名称**: `openai`
- **支持模型**: GPT-3.5-turbo, GPT-4, GPT-4-turbo等
- **官网**: https://openai.com/

**配置示例**:


### 2. 百度文心一言
- **Provider名称**: `baidu` 或 `wenxin`
- **支持模型**: ERNIE-Bot, ERNIE-Bot-turbo等
- **官网**: https://cloud.baidu.com/product/wenxinworkshop

**配置示例**:


### 3. 阿里通义千问
- **Provider名称**: `ali`, `qwen`, 或 `tongyi`
- **支持模型**: qwen-turbo, qwen-plus, qwen-max等
- **官网**: https://dashscope.aliyun.com/

**配置示例**:


### 4. 演示模式
- **Provider名称**: `mock`
- **说明**: 用于演示和测试，不需要真实API Key

**配置示例**:


## 快速开始

### 1. 复制配置文件


### 2. 编辑配置文件
根据您选择的提供商，编辑 `.env` 文件中的相应配置。

### 3. 启动服务


## 详细配置说明

### OpenAI配置
1. 访问 https://platform.openai.com/api-keys
2. 创建新的API Key
3. 将API Key填入 `LLM_API_KEY`

**支持的模型**:
- `gpt-3.5-turbo` (推荐，性价比高)
- `gpt-4` (更强大，但费用较高)
- `gpt-4-turbo` (最新模型)

### 百度文心一言配置
1. 访问 https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application
2. 创建应用获取 API Key、Secret Key
3. 配置相应的环境变量

**注意**: 百度API使用OAuth2.0认证，需要先获取access_token。

### 阿里通义千问配置
1. 访问 https://dashscope.console.aliyun.com/
2. 获取API Key
3. 选择合适的模型

**支持的模型**:
- `qwen-turbo` (快速响应)
- `qwen-plus` (平衡性能)
- `qwen-max` (最强性能)

## 环境变量说明

| 变量名 | 必填 | 说明 | 默认值 |
|--------|------|------|--------|
| `LLM_PROVIDER` | 是 | 大模型提供商类型 | `openai` |
| `LLM_API_KEY` | 是* | API密钥 | - |
| `LLM_API_URL` | 否 | API端点URL | 各Provider默认值 |
| `LLM_MODEL` | 否 | 模型名称 | 各Provider默认值 |

*注: mock模式不需要API Key

## 测试配置

启动服务后，可以通过以下方式测试配置：

### 1. 健康检查接口


### 2. 提问测试


### 3. 前端界面测试
访问 http://localhost:3000 使用Web界面测试。

## 故障排除

### 常见问题

**1. API Key无效**
- 检查API Key是否正确
- 确认API Key是否有足够的配额
- 验证API Key的权限设置

**2. 网络连接问题**
- 检查网络连接
- 确认防火墙设置
- 考虑使用代理（如需要）

**3. 模型不存在**
- 确认模型名称是否正确
- 检查账户是否有该模型的访问权限

**4. 配额超限**
- 检查API使用配额
- 考虑升级账户套餐

### 日志调试
启动服务时会显示当前使用的Provider信息：


## 扩展支持

如需添加新的大模型提供商，请：

1. 在 `llm_providers.go` 中实现 `LLMProvider` 接口
2. 在 `llm_client.go` 的 `NewLLMClient()` 中添加相应的case
3. 更新配置文件和文档

## 安全建议

1. **不要将API Key提交到版本控制系统**
2. **使用环境变量或安全的配置管理工具**
3. **定期轮换API Key**
4. **监控API使用情况**
5. **设置合理的使用限制**

## 费用优化

1. **选择合适的模型** - 根据需求选择性价比最高的模型
2. **控制请求频率** - 避免不必要的API调用
3. **优化提示词** - 使用更精确的提示词减少token消耗
4. **设置使用限制** - 防止意外的大量调用