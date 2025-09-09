# LLM Providers 包

这个包包含了所有大模型API的Provider实现，采用模块化设计，每个Provider都在独立的文件中。

## 文件结构

```
providers/
├── interface.go          # 接口定义
├── chat_types.go         # 聊天完成相关数据结构
├── openai.go            # OpenAI Provider
├── baidu.go             # 百度文心一言Provider
├── ali.go               # 阿里通义千问Provider
├── bella.go             # Bella智能问答Provider
├── mock.go              # Mock Provider (测试用)
├── bella_example.go     # Bella Provider使用示例
├── BELLA_PROVIDER.md    # Bella Provider详细文档
└── README.md            # 本文档
```

## 接口定义

所有Provider都必须实现`LLMProvider`接口：

```go
type LLMProvider interface {
    AskQuestion(question string) (string, error)
    GetProviderName() string
    CheckConnection() error
}
```

扩展的聊天完成Provider还需要实现`ChatCompletionProvider`接口：

```go
type ChatCompletionProvider interface {
    LLMProvider
    ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
    ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (<-chan *ChatCompletionStreamResponse, <-chan error)
}
```

## 支持的Provider

### 1. OpenAI Provider (`openai.go`)
- 支持GPT-3.5-turbo、GPT-4等模型
- 使用标准的OpenAI API格式
- 默认模型：`gpt-3.5-turbo`

### 2. 百度文心一言Provider (`baidu.go`)
- 支持ERNIE-Bot系列模型
- 使用百度千帆平台API
- 需要OAuth2.0认证

### 3. 阿里通义千问Provider (`ali.go`)
- 支持qwen-turbo、qwen-plus、qwen-max等模型
- 使用阿里云灵积平台API
- 默认模型：`qwen-turbo`

### 4. Bella智能问答Provider (`bella.go`) ⭐ 新增
- **完整的聊天完成API支持**
- **工具调用功能**
- **流式响应**
- **图片输入支持**
- **reasoning_content字段支持**
- **多种消息类型（system, user, assistant, tool, developer）**
- 默认模型：`gpt-4o`
- 详细文档：[BELLA_PROVIDER.md](./BELLA_PROVIDER.md)

### 5. Mock Provider (`mock.go`)
- 用于演示和测试
- 不需要真实的API Key
- 返回预设的模拟回答

## Bella Provider 快速使用

```go
// 创建Bella provider
config := providers.ProviderConfig{
    Name:   "Bella",
    APIKey: "your-api-key-here",
    APIURL: "https://api.bella.com/v1/chat/completions",
    Model:  "gpt-4o",
}
provider := providers.NewBellaProvider(config)

// 基础问答
answer, err := provider.AskQuestion("你好")

// 聊天完成
req := &providers.ChatCompletionRequest{
    Model: "gpt-4o",
    Messages: []providers.Message{
        {Role: "user", Content: "请介绍一下人工智能"},
    },
}
resp, err := provider.ChatCompletion(context.Background(), req)

// 流式响应
responseChan, errorChan := provider.ChatCompletionStream(context.Background(), req)
```

## 添加新Provider

要添加新的Provider，请按以下步骤：

1. **创建新的Provider文件**（如`tencent.go`）
2. **实现LLMProvider接口**
3. **在`llm_client.go`中添加对应的case**
4. **更新配置文档和示例**

### 示例：添加腾讯混元Provider

```go
package providers

type TencentProvider struct {
    config ProviderConfig
    client *http.Client
}

func NewTencentProvider(config ProviderConfig) *TencentProvider {
    return &TencentProvider{
        config: config,
        client: &http.Client{Timeout: 30 * time.Second},
    }
}

func (p *TencentProvider) AskQuestion(question string) (string, error) {
    // 实现腾讯混元API调用
}

func (p *TencentProvider) GetProviderName() string {
    return "Tencent"
}

func (p *TencentProvider) CheckConnection() error {
    // 实现连接检查
}
```

然后在`llm_client.go`中添加：

```go
case "tencent":
    return NewTencentProvider(config), nil
```

## 配置说明

每个Provider使用统一的`ProviderConfig`结构体：

```go
type ProviderConfig struct {
    Name   string  // Provider名称
    APIKey string  // API密钥
    APIURL string  // API端点URL
    Model  string  // 模型名称
}
```

## 最佳实践

1. **错误处理**：所有Provider都应该有完善的错误处理
2. **超时设置**：设置合理的HTTP客户端超时时间
3. **日志记录**：记录关键操作和错误信息
4. **配置验证**：在构造函数中验证必要的配置参数
5. **响应解析**：正确解析各厂商不同的API响应格式

## 测试

每个Provider都应该包含相应的测试：

```go
func TestBellaProvider_AskQuestion(t *testing.T) {
    config := ProviderConfig{
        Name:   "Bella",
        APIKey: "test-key",
        Model:  "gpt-4o",
    }
    provider := NewBellaProvider(config)
    
    // 测试基础功能
    answer, err := provider.AskQuestion("Hello")
    assert.NoError(t, err)
    assert.NotEmpty(t, answer)
}
```

## 扩展性

这个架构设计具有良好的扩展性：

- **易于添加新Provider**：只需实现接口即可
- **统一的调用方式**：上层代码无需修改
- **独立的文件管理**：每个Provider独立维护
- **配置灵活**：支持不同Provider的特殊配置需求
- **功能分层**：基础Provider和扩展Provider分离

## 特殊功能

### Bella Provider 高级特性

Bella Provider提供了最完整的聊天完成功能：

1. **工具调用**：支持函数调用和工具集成
2. **流式响应**：实时获取生成内容
3. **多模态输入**：支持文本和图片输入
4. **推理内容**：获取模型的思考过程
5. **自定义参数**：支持传递任意额外参数

详细使用方法请参考：
- [BELLA_PROVIDER.md](./BELLA_PROVIDER.md) - 完整文档
- [bella_example.go](./bella_example.go) - 使用示例