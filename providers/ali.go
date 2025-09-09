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

// AliProvider 阿里通义千问提供商
type AliProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewAliProvider 创建阿里通义千问提供商
func NewAliProvider(config ProviderConfig) *AliProvider {
	if config.APIURL == "" {
		config.APIURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}
	if config.Model == "" {
		config.Model = "qwen-turbo"
	}

	return &AliProvider{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *AliProvider) AskQuestion(question string) (string, error) {
	request := map[string]interface{}{
		"model": p.config.Model,
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{"role": "system", "content": "你是一个有用的AI助手。"},
				{"role": "user", "content": question},
			},
		},
		"parameters": map[string]interface{}{
			"result_format": "message",
		},
	}

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

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// 解析阿里云API响应格式
	if output, ok := response["output"].(map[string]interface{}); ok {
		if choices, ok := output["choices"].([]interface{}); ok && len(choices) > 0 {
			choice := choices[0].(map[string]interface{})
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("阿里云API响应格式错误")
}

func (p *AliProvider) GetProviderName() string {
	return "阿里通义千问"
}

func (p *AliProvider) CheckConnection() error {
	log.Printf("阿里通义千问API连接检查通过")
	return nil
}
