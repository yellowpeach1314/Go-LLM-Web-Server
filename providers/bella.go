package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// BellaProvider Bella智能问答提供商
type BellaProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewBellaProvider 创建Bella提供商
func NewBellaProvider(config ProviderConfig) *BellaProvider {
	if config.APIURL == "" {
		config.APIURL = "https://api.bella.com/v1/chat/completions"
	}
	if config.Model == "" {
		config.Model = "gpt-4o"
	}

	return &BellaProvider{
		config: config,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// AskQuestion 实现基础的问答接口
func (p *BellaProvider) AskQuestion(question string) (string, error) {
	req := &ChatCompletionRequest{
		Model: p.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手，请用中文回答问题。",
			},
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	resp, err := p.ChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("没有返回任何选择")
	}

	if resp.Choices[0].Message != nil && resp.Choices[0].Message.Content != nil {
		if content, ok := resp.Choices[0].Message.Content.(string); ok {
			return content, nil
		}
	}

	return "", fmt.Errorf("响应格式错误")
}

// ChatCompletion 实现聊天完成功能
func (p *BellaProvider) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 确保不是流式请求
	req.Stream = false

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// ChatCompletionStream 实现流式聊天完成功能
func (p *BellaProvider) ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (<-chan *ChatCompletionStreamResponse, <-chan error) {
	responseChan := make(chan *ChatCompletionStreamResponse, 100)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		// 设置为流式请求
		req.Stream = true

		jsonData, err := json.Marshal(req)
		if err != nil {
			errorChan <- fmt.Errorf("序列化请求失败: %v", err)
			return
		}

		log.Printf("Bella流式请求数据: %s", string(jsonData))

		httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.APIURL, bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("创建HTTP请求失败: %v", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
		httpReq.Header.Set("Accept", "text/event-stream")

		log.Printf("发送Bella流式请求到: %s", p.config.APIURL)

		resp, err := p.client.Do(httpReq)
		if err != nil {
			errorChan <- fmt.Errorf("发送HTTP请求失败: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Bella响应状态码: %d", resp.StatusCode)
		log.Printf("Bella响应头: %v", resp.Header)

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Bella API错误响应: %s", string(body))
			errorChan <- fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
			return
		}

		log.Printf("开始读取Bella流式响应...")
		scanner := bufio.NewScanner(resp.Body)
		lineCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineCount++
			log.Printf("Bella流式响应第%d行: %s", lineCount, line)

			// 跳过空行和注释行
			if line == "" || strings.HasPrefix(line, ":") {
				log.Printf("跳过空行或注释行")
				continue
			}

			// 处理SSE数据
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				log.Printf("处理SSE数据: %s", data)

				// 检查是否为结束标记
				if data == "[DONE]" {
					log.Printf("收到结束标记[DONE]")
					return
				}

				var streamResp ChatCompletionStreamResponse
				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					log.Printf("解析流式响应失败: %v, 数据: %s", err, data)
					continue
				}

				log.Printf("成功解析流式响应: %+v", streamResp)

				select {
				case responseChan <- &streamResp:
					log.Printf("成功发送流式响应到通道")
				case <-ctx.Done():
					log.Printf("上下文取消，退出流式响应")
					return
				}
			} else {
				log.Printf("非SSE数据行: %s", line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("扫描器错误: %v", err)
			errorChan <- fmt.Errorf("读取流式响应失败: %v", err)
		}

		log.Printf("Bella流式响应读取完成，总共处理了%d行", lineCount)
	}()

	return responseChan, errorChan
}

// GetProviderName 获取提供商名称
func (p *BellaProvider) GetProviderName() string {
	return "Bella"
}

// CheckConnection 检查连接
func (p *BellaProvider) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &ChatCompletionRequest{
		Model: p.config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		MaxTokens: func() *int { i := 1; return &i }(),
	}

	_, err := p.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("Bella API连接检查失败: %v", err)
	}

	log.Printf("Bella API连接检查通过")
	return nil
}

// CreateToolCallRequest 创建工具调用请求的辅助方法
func (p *BellaProvider) CreateToolCallRequest(messages []Message, tools []Tool, toolChoice interface{}) *ChatCompletionRequest {
	req := &ChatCompletionRequest{
		Model:      p.config.Model,
		Messages:   messages,
		Tools:      tools,
		ToolChoice: toolChoice,
	}
	return req
}

// CreateImageRequest 创建图片输入请求的辅助方法
func (p *BellaProvider) CreateImageRequest(text string, imageURL string, detail string) *ChatCompletionRequest {
	content := []ContentPart{
		{
			Type: "text",
			Text: text,
		},
		{
			Type: "image_url",
			ImageURL: &ImageURL{
				URL:    imageURL,
				Detail: detail,
			},
		},
	}

	req := &ChatCompletionRequest{
		Model: p.config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	return req
}
