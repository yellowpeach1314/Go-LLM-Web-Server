package providers

import "context"

// LLMProvider 大模型提供商接口
type LLMProvider interface {
	AskQuestion(question string) (string, error)
	GetProviderName() string
	CheckConnection() error
}

// ChatCompletionProvider 聊天完成提供商接口
type ChatCompletionProvider interface {
	LLMProvider
	ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (<-chan *ChatCompletionStreamResponse, <-chan error)
}

// ProviderConfig 提供商配置
type ProviderConfig struct {
	Name   string
	APIKey string
	APIURL string
	Model  string
}
