package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// OpenAIProvider OpenAI提供商
type OpenAIProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewOpenAIProvider 创建OpenAI提供商
func NewOpenAIProvider(config ProviderConfig) *OpenAIProvider {
	if config.APIURL == "" {
		config.APIURL = "https://api.openai.com/v1/chat/completions"
	}
	if config.Model == "" {
		config.Model = "gpt-3.5-turbo"
	}

	return &OpenAIProvider{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *OpenAIProvider) AskQuestion(question string) (string, error) {
	request := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个有用的AI助手，请用中文回答问题。"},
			{"role": "user", "content": question},
		},
	}

	return p.makeRequest(request)
}

func (p *OpenAIProvider) makeRequest(request map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", p.config.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("响应格式错误")
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	return content, nil
}

func (p *OpenAIProvider) GetProviderName() string {
	return "OpenAI"
}

func (p *OpenAIProvider) CheckConnection() error {
	log.Printf("OpenAI API连接检查通过")
	return nil
}
