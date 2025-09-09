package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BaiduProvider 百度文心一言提供商
type BaiduProvider struct {
	config      ProviderConfig
	client      *http.Client
	accessToken string
}

// NewBaiduProvider 创建百度文心一言提供商
func NewBaiduProvider(config ProviderConfig) *BaiduProvider {
	if config.APIURL == "" {
		config.APIURL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant"
	}

	return &BaiduProvider{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *BaiduProvider) AskQuestion(question string) (string, error) {
	// 百度API需要先获取access_token
	if p.accessToken == "" {
		if err := p.getAccessToken(); err != nil {
			return "", err
		}
	}

	request := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": question},
		},
	}

	apiURL := fmt.Sprintf("%s?access_token=%s", p.config.APIURL, p.accessToken)
	jsonData, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

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

	if result, ok := response["result"].(string); ok {
		return result, nil
	}

	return "", fmt.Errorf("百度API响应格式错误")
}

func (p *BaiduProvider) getAccessToken() error {
	// 这里需要实现获取百度access_token的逻辑
	// 实际使用时需要配置client_id和client_secret
	p.accessToken = "mock_access_token"
	return nil
}

func (p *BaiduProvider) GetProviderName() string {
	return "百度文心一言"
}

func (p *BaiduProvider) CheckConnection() error {
	return p.getAccessToken()
}
