# Bella智能问答Provider

Bella Provider是基于智能问答接口文档设计的Go语言实现，支持完整的聊天完成功能，包括工具调用、流式响应、图片输入和推理内容等高级特性。

## 特性

- ✅ 完整的聊天完成API支持
- ✅ 工具调用功能
- ✅ 流式响应
- ✅ 图片输入支持
- ✅ reasoning_content字段支持
- ✅ 多种消息类型（system, user, assistant, tool, developer）
- ✅ 自定义参数传递
- ✅ 错误处理和重试机制

## 快速开始

### 1. 创建Provider实例

```go
package main

import (
    "github.com/your-project/providers"
)

func main() {
    config := providers.ProviderConfig{
        Name:   "Bella",
        APIKey: "your-api-key-here",
        APIURL: "https://api.bella.com/v1/chat/completions",
        Model:  "gpt-4o",
    }
    
    provider := providers.NewBellaProvider(config)
}
```

### 2. 基础问答

```go
answer, err := provider.AskQuestion("你好，请介绍一下自己。")
if err != nil {
    log.Fatal(err)
}
fmt.Println(answer)
```

### 3. 聊天完成

```go
req := &providers.ChatCompletionRequest{
    Model: "gpt-4o",
    Messages: []providers.Message{
        {
            Role:    "system",
            Content: "你是一个有帮助的助手。",
        },
        {
            Role:    "user",
            Content: "请解释一下什么是人工智能？",
        },
    },
    Temperature: func() *float64 { f := 0.7; return &f }(),
    MaxTokens:   func() *int { i := 500; return &i }(),
}

resp, err := provider.ChatCompletion(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

// 访问响应内容
if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
    fmt.Println("回答:", resp.Choices[0].Message.Content)
    if resp.Choices[0].Message.ReasoningContent != "" {
        fmt.Println("推理过程:", resp.Choices[0].Message.ReasoningContent)
    }
}
```

## 高级功能

### 工具调用

```go
req := &providers.ChatCompletionRequest{
    Model: "gpt-4o",
    Messages: []providers.Message{
        {
            Role:    "user",
            Content: "北京今天的天气怎么样？",
        },
    },
    Tools: []providers.Tool{
        {
            Type: "function",
            Function: providers.Function{
                Name:        "get_weather",
                Description: "获取指定城市的天气信息",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "location": map[string]interface{}{
                            "type":        "string",
                            "description": "城市名称，如北京、上海",
                        },
                        "unit": map[string]interface{}{
                            "type":        "string",
                            "enum":        []string{"celsius", "fahrenheit"},
                            "description": "温度单位",
                        },
                    },
                    "required": []string{"location"},
                },
                Strict: true,
            },
        },
    },
    ToolChoice: "auto",
}

resp, err := provider.ChatCompletion(context.Background(), req)
// 处理工具调用响应...
```

### 图片输入

```go
req := provider.CreateImageRequest(
    "这张图片是什么内容？",
    "https://example.com/image.jpg",
    "high",
)

resp, err := provider.ChatCompletion(context.Background(), req)
// 处理图片分析结果...
```

### 流式响应

```go
req := &providers.ChatCompletionRequest{
    Model: "gpt-4o",
    Messages: []providers.Message{
        {
            Role:    "user",
            Content: "请写一首关于春天的诗。",
        },
    },
    Stream: true,
}

responseChan, errorChan := provider.ChatCompletionStream(context.Background(), req)

for {
    select {
    case resp, ok := <-responseChan:
        if !ok {
            return // 流式响应完成
        }
        if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
            if content, ok := resp.Choices[0].Delta.Content.(string); ok {
                fmt.Print(content)
            }
        }
    case err, ok := <-errorChan:
        if ok && err != nil {
            log.Printf("流式响应错误: %v", err)
            return
        }
    }
}
```

## 消息类型

### System消息
```go
{
    Role:    "system",
    Content: "你是一个有帮助的助手。",
}
```

### Developer消息
```go
{
    Role:    "developer",
    Content: "你是一个有帮助的助手，专注于回答用户的问题。",
}
```

### User消息
```go
{
    Role:    "user",
    Content: "你好，请介绍一下自己。",
}
```

### Assistant消息
```go
{
    Role:    "assistant",
    Content: "我是一个AI助手，可以回答问题和提供信息。",
}
```

### Tool消息
```go
{
    Role:       "tool",
    Content:    `{"temperature":32,"unit":"celsius","description":"晴朗"}`,
    ToolCallID: "call_abc123",
}
```

## 配置选项

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| Name | string | 是 | Provider名称 |
| APIKey | string | 是 | API密钥 |
| APIURL | string | 否 | API端点URL，默认为官方地址 |
| Model | string | 否 | 模型名称，默认为gpt-4o |

## 请求参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| model | string | 是 | 要使用的模型ID |
| messages | []Message | 是 | 包含对话历史的消息数组 |
| tools | []Tool | 否 | 模型可以调用的工具列表 |
| tool_choice | interface{} | 否 | 控制模型是否调用工具 |
| temperature | *float64 | 否 | 采样温度，默认为1 |
| top_p | *float64 | 否 | 核采样的概率质量，默认为1 |
| max_tokens | *int | 否 | 生成的最大令牌数 |
| stream | bool | 否 | 是否启用流式响应，默认为false |

## 响应格式

### 非流式响应
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-4o",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "你好！我能帮你什么忙吗？",
      "reasoning_content": "用户用中文问候，我应该用中文回复。"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

### 流式响应
```
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"你"},"finish_reason":null}]}
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"好"},"finish_reason":null}]}
data: [DONE]
```

## 错误处理

Provider会自动处理以下错误情况：
- 网络连接错误
- API认证错误
- 请求参数错误
- 响应解析错误

```go
resp, err := provider.ChatCompletion(context.Background(), req)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "401"):
        log.Println("API密钥无效")
    case strings.Contains(err.Error(), "429"):
        log.Println("请求频率过高")
    case strings.Contains(err.Error(), "500"):
        log.Println("服务器内部错误")
    default:
        log.Printf("其他错误: %v", err)
    }
}
```

## 最佳实践

1. **使用Context控制超时**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := provider.ChatCompletion(ctx, req)
```

2. **合理设置参数**
```go
req := &providers.ChatCompletionRequest{
    Model:       "gpt-4o",
    Messages:    messages,
    Temperature: func() *float64 { f := 0.7; return &f }(), // 平衡创造性和一致性
    MaxTokens:   func() *int { i := 1000; return &i }(),    // 限制响应长度
}
```

3. **处理流式响应**
```go
// 使用缓冲通道避免阻塞
responseChan, errorChan := provider.ChatCompletionStream(ctx, req)

// 设置超时处理
timeout := time.After(60 * time.Second)

for {
    select {
    case resp := <-responseChan:
        // 处理响应
    case err := <-errorChan:
        // 处理错误
    case <-timeout:
        log.Println("流式响应超时")
        return
    }
}
```

## 示例代码

完整的使用示例请参考 `bella_example.go` 文件。

## 注意事项

1. 确保API密钥的安全性，不要在代码中硬编码
2. 合理设置超时时间，避免长时间等待
3. 流式响应需要正确处理通道关闭
4. 工具调用需要实现相应的工具函数
5. 图片输入需要确保图片URL可访问

## 兼容性

- Go 1.18+
- 支持所有符合OpenAI Chat Completions API规范的服务
- 扩展支持reasoning_content字段