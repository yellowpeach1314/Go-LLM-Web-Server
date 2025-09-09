package providers

import (
	"context"
	"fmt"
	"log"
)

// BellaExample Bella provider使用示例
func BellaExample() {
	// 创建Bella provider配置
	config := ProviderConfig{
		Name:   "Bella",
		APIKey: "your-api-key-here",
		APIURL: "https://api.bella.com/v1/chat/completions",
		Model:  "gpt-4o",
	}

	// 创建provider实例
	provider := NewBellaProvider(config)

	// 示例1: 基础问答
	fmt.Println("=== 基础问答示例 ===")
	answer, err := provider.AskQuestion("你好，请介绍一下自己。")
	if err != nil {
		log.Printf("基础问答失败: %v", err)
	} else {
		fmt.Printf("回答: %s\n", answer)
	}

	// 示例2: 聊天完成
	fmt.Println("\n=== 聊天完成示例 ===")
	chatReq := &ChatCompletionRequest{
		Model: config.Model,
		Messages: []Message{
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

	chatResp, err := provider.ChatCompletion(context.Background(), chatReq)
	if err != nil {
		log.Printf("聊天完成失败: %v", err)
	} else {
		if len(chatResp.Choices) > 0 && chatResp.Choices[0].Message != nil {
			fmt.Printf("回答: %v\n", chatResp.Choices[0].Message.Content)
			if chatResp.Choices[0].Message.ReasoningContent != "" {
				fmt.Printf("推理过程: %s\n", chatResp.Choices[0].Message.ReasoningContent)
			}
		}
	}

	// 示例3: 工具调用
	fmt.Println("\n=== 工具调用示例 ===")
	toolReq := &ChatCompletionRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: "北京今天的天气怎么样？",
			},
		},
		Tools: []Tool{
			{
				Type: "function",
				Function: Function{
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

	toolResp, err := provider.ChatCompletion(context.Background(), toolReq)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		if len(toolResp.Choices) > 0 && toolResp.Choices[0].Message != nil {
			if len(toolResp.Choices[0].Message.ToolCalls) > 0 {
				fmt.Printf("工具调用: %+v\n", toolResp.Choices[0].Message.ToolCalls[0])
			}
		}
	}

	// 示例4: 图片输入
	fmt.Println("\n=== 图片输入示例 ===")
	imageReq := provider.CreateImageRequest(
		"这张图片是什么内容？",
		"https://example.com/image.jpg",
		"high",
	)

	imageResp, err := provider.ChatCompletion(context.Background(), imageReq)
	if err != nil {
		log.Printf("图片输入失败: %v", err)
	} else {
		if len(imageResp.Choices) > 0 && imageResp.Choices[0].Message != nil {
			fmt.Printf("图片分析结果: %v\n", imageResp.Choices[0].Message.Content)
		}
	}

	// 示例5: 流式响应
	fmt.Println("\n=== 流式响应示例 ===")
	streamReq := &ChatCompletionRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: "请写一首关于春天的诗。",
			},
		},
		Stream: true,
	}

	responseChan, errorChan := provider.ChatCompletionStream(context.Background(), streamReq)

	fmt.Print("流式回答: ")
	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				fmt.Println("\n流式响应完成")
				return
			}
			if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
				if resp.Choices[0].Delta.Content != nil {
					if content, ok := resp.Choices[0].Delta.Content.(string); ok {
						fmt.Print(content)
					}
				}
				if resp.Choices[0].Delta.ReasoningContent != "" {
					fmt.Printf("[推理: %s]", resp.Choices[0].Delta.ReasoningContent)
				}
			}
		case err, ok := <-errorChan:
			if ok && err != nil {
				log.Printf("流式响应错误: %v", err)
				return
			}
		}
	}
}
