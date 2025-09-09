package providers

import (
	"fmt"
	"log"
)

// MockProvider 模拟Provider，用于演示和测试
type MockProvider struct{}

// NewMockProvider 创建Mock Provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) AskQuestion(question string) (string, error) {
	responses := map[string]string{
		"你好":      "你好！我是AI助手，很高兴为您服务！有什么可以帮助您的吗？",
		"再见":      "再见！祝您有美好的一天！",
		"谢谢":      "不客气！很高兴能够帮助到您。",
		"今天天气怎么样": "抱歉，我无法获取实时天气信息。建议您查看天气预报应用或网站。",
		"你是谁":     "我是一个AI助手，专门为您提供问答服务。",
	}

	if answer, exists := responses[question]; exists {
		return answer, nil
	}

	return fmt.Sprintf("感谢您的问题：「%s」。这是一个模拟回答，因为当前运行在演示模式下。要获得真实的AI回答，请配置相应的LLM Provider。", question), nil
}

func (p *MockProvider) GetProviderName() string {
	return "Mock Provider (演示模式)"
}

func (p *MockProvider) CheckConnection() error {
	log.Println("Mock Provider连接检查通过")
	return nil
}
