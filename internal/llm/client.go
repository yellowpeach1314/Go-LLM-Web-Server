package llm

import (
	"context"
	"fmt"
	"go-base-web-server/providers"
	"log"
	"strings"
)

// Client LLM客户端
type Client struct {
	provider     providers.LLMProvider
	chatProvider providers.ChatCompletionProvider // 添加流式provider
	config       Config                           // 添加配置字段
}

// Config LLM配置
type Config struct {
	Provider string
	APIKey   string
	APIURL   string
	Model    string
}

// NewClient 创建新的LLM客户端，支持多种Provider
func NewClient(config Config) *Client {
	providerType := strings.ToLower(config.Provider)
	if providerType == "" {
		providerType = "openai" // 默认使用OpenAI
	}

	providerConfig := providers.ProviderConfig{
		Name:   providerType,
		APIKey: config.APIKey,
		APIURL: config.APIURL,
		Model:  config.Model,
	}

	var provider providers.LLMProvider

	switch providerType {
	case "openai":
		provider = providers.NewOpenAIProvider(providerConfig)
	case "baidu", "wenxin":
		provider = providers.NewBaiduProvider(providerConfig)
	case "ali", "qwen", "tongyi":
		provider = providers.NewAliProvider(providerConfig)
	case "bella":
		provider = providers.NewBellaProvider(providerConfig)
	case "gemini", "google":
		provider = providers.NewGeminiProvider(providerConfig)
	case "mock":
		provider = providers.NewMockProvider()
	default:
		log.Printf("未知的Provider类型: %s，使用Mock Provider", providerType)
		provider = providers.NewMockProvider()
	}

	log.Printf("初始化LLM Provider: %s", provider.GetProviderName())

	// 检查是否支持流式聊天
	if chatProvider, ok := provider.(providers.ChatCompletionProvider); ok {
		log.Printf("Provider %s 支持流式聊天", provider.GetProviderName())
		return &Client{
			provider:     provider,
			chatProvider: chatProvider,
			config:       config, // 保存配置
		}
	}

	log.Printf("Provider %s 不支持流式聊天", provider.GetProviderName())
	return &Client{
		provider:     provider,
		chatProvider: nil,
		config:       config, // 保存配置
	}
}

// AskQuestion 向LLM提问
func (c *Client) AskQuestion(question string) (string, error) {
	if c.provider == nil {
		return "", fmt.Errorf("LLM Provider未初始化")
	}

	log.Printf("使用 %s 处理问题: %s", c.provider.GetProviderName(), question)

	answer, err := c.provider.AskQuestion(question)
	if err != nil {
		log.Printf("LLM调用失败: %v", err)
		return "", err
	}

	log.Printf("LLM响应成功，答案长度: %d", len(answer))
	return answer, nil
}

// CheckConnection 检查API连接
func (c *Client) CheckConnection() error {
	if c.provider == nil {
		return fmt.Errorf("LLM Provider未初始化")
	}

	return c.provider.CheckConnection()
}

// GetProviderInfo 获取当前Provider信息
func (c *Client) GetProviderInfo() map[string]interface{} {
	if c.provider == nil {
		return map[string]interface{}{
			"provider": "未初始化",
			"status":   "error",
		}
	}

	return map[string]interface{}{
		"provider": c.provider.GetProviderName(),
		"status":   "active",
	}
}

// ChatCompletionStream 流式聊天完成
func (c *Client) ChatCompletionStream(ctx context.Context, question string) (<-chan *providers.ChatCompletionStreamResponse, <-chan error, error) {
	if c.chatProvider == nil {
		return nil, nil, fmt.Errorf("当前Provider不支持流式聊天")
	}

	// 使用配置中的模型名称，如果为空则使用默认值
	model := c.config.Model
	if model == "" {
		model = "gpt-3.5-turbo" // 默认模型
	}

	req := &providers.ChatCompletionRequest{
		Model: model, // 使用配置中的模型
		Messages: []providers.Message{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手，请用中文回答问题。",
			},
			{
				Role:    "user",
				Content: question,
			},
		},
		Stream: true,
	}

	log.Printf("使用 %s 处理流式问题 (模型: %s): %s", c.provider.GetProviderName(), model, question)

	responseChan, errorChan := c.chatProvider.ChatCompletionStream(ctx, req)
	return responseChan, errorChan, nil
}

// SupportsStreaming 检查是否支持流式聊天
func (c *Client) SupportsStreaming() bool {
	return c.chatProvider != nil
}
