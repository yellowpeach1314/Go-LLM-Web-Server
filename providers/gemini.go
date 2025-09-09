package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// GeminiProvider Google Gemini提供商
type GeminiProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewGeminiProvider 创建Gemini提供商
func NewGeminiProvider(config ProviderConfig) *GeminiProvider {
	if config.APIURL == "" {
		config.APIURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	if config.Model == "" {
		config.Model = "gemini-2.5-flash"
	}

	return &GeminiProvider{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GeminiRequest Gemini API请求结构
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent Gemini内容结构
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart Gemini部分结构
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse Gemini API响应结构
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate Gemini候选结构
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

func (p *GeminiProvider) AskQuestion(question string) (string, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/models/%s:generateContent", p.config.APIURL, p.config.Model)

	// 构建请求体
	request := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: question},
				},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("响应格式错误或无内容")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (p *GeminiProvider) GetProviderName() string {
	return "Google Gemini"
}

func (p *GeminiProvider) CheckConnection() error {
	_, err := p.AskQuestion("Hello")
	if err != nil {
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			return fmt.Errorf("Gemini API认证失败，请检查API密钥")
		}
		return fmt.Errorf("Gemini API连接检查失败: %v", err)
	}

	log.Printf("Gemini API连接检查通过")
	return nil
}
