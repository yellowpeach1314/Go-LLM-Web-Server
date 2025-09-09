package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-base-web-server/providers"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// QAStorage QA存储接口
type QAStorage interface {
	SaveQuestion(question string, userID *int) (int, error)
	UpdateAnswer(id int, answer string) error
	GetRecord(id int) (interface{}, error)
	GetAllRecords() (interface{}, error)
	GetRecordsByUserID(userID int) (interface{}, error)
}

// LLMClient LLM客户端接口
type LLMClient interface {
	AskQuestion(question string) (string, error)
	CheckConnection() error
	GetProviderInfo() map[string]interface{}
	// 添加流式聊天方法
	ChatCompletionStream(ctx context.Context, question string) (<-chan *providers.ChatCompletionStreamResponse, <-chan error, error)
	SupportsStreaming() bool
}

// App 应用结构体，包含所有依赖
type App struct {
	qaStorage QAStorage
	llmClient LLMClient
}

// NewApp 创建新的应用实例
func NewApp(qaStorage QAStorage, llmClient LLMClient) *App {
	return &App{
		qaStorage: qaStorage,
		llmClient: llmClient,
	}
}

// HomeHandler API信息处理器
func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	apiInfo := map[string]interface{}{
		"name":        "LLM 问答系统 API",
		"version":     "2.0.0",
		"description": "现代化的前后端分离问答系统后端API - 支持用户认证",
		"endpoints": map[string]interface{}{
			"GET /":                        "API信息",
			"GET /api/health":              "健康检查",
			"POST /api/auth/register":      "用户注册",
			"POST /api/auth/login":         "用户登录",
			"POST /api/auth/logout":        "用户登出",
			"GET /api/ask":                 "提问接口 (参数: prompt)",
			"GET /api/ask/stream":          "流式提问接口 (参数: prompt) - SSE",
			"GET /api/records":             "获取所有问答记录",
			"GET /api/records/{id}":        "获取特定记录",
			"GET /api/user/profile":        "获取用户资料 (需要认证)",
			"POST /api/user/refresh-token": "刷新token (需要认证)",
			"GET /api/user/records":        "获取用户记录 (需要认证)",
			"GET /api/user/users":          "获取用户列表 (需要认证)",
		},
		"authentication": map[string]string{
			"type":   "Bearer Token (JWT)",
			"header": "Authorization: Bearer <token>",
		},
		"cors": "已启用跨域支持",
	}

	json.NewEncoder(w).Encode(apiInfo)
}

// AskHandler 提问处理器 - 核心业务流程（支持用户认证）
func (app *App) AskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	question := r.URL.Query().Get("prompt")
	if question == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "缺少prompt参数"})
		return
	}

	log.Printf("收到问题: %s", question)

	// 获取用户ID（如果已认证）
	var userID *int
	if user, ok := getUserFromContext(r); ok {
		if id := getUserID(user); id > 0 {
			userID = &id
			log.Printf("认证用户提问: ID %d", id)
		}
	} else {
		log.Printf("匿名用户提问")
	}

	// 1. 保存问题到数据库
	recordID, err := app.qaStorage.SaveQuestion(question, userID)
	if err != nil {
		log.Printf("保存问题失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "保存问题失败"})
		return
	}

	// 2. 调用LLM获取答案
	answer, err := app.llmClient.AskQuestion(question)
	if err != nil {
		log.Printf("LLM调用失败: %v", err)
		app.qaStorage.UpdateAnswer(recordID, "抱歉，AI服务暂时不可用")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "AI服务不可用"})
		return
	}

	// 3. 更新数据库中的答案
	err = app.qaStorage.UpdateAnswer(recordID, answer)
	if err != nil {
		log.Printf("更新答案失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "保存答案失败"})
		return
	}

	// 4. 返回完整的问答结果
	response := map[string]interface{}{
		"id":       recordID,
		"question": question,
		"answer":   answer,
		"user_id":  userID,
		"status":   "success",
	}

	log.Printf("问答完成，ID: %d", recordID)
	json.NewEncoder(w).Encode(response)
}

// GetRecordsHandler 获取所有记录处理器
func (app *App) GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	records, err := app.qaStorage.GetAllRecords()
	if err != nil {
		log.Printf("获取记录失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取记录失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "获取记录成功",
		"data":    records,
		"status":  "success",
	})
}

// GetRecordHandler 获取单个记录处理器
func (app *App) GetRecordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的ID格式"})
		return
	}

	record, err := app.qaStorage.GetRecord(id)
	if err != nil {
		log.Printf("获取记录失败: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "记录不存在"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "获取记录成功",
		"data":    record,
		"status":  "success",
	})
}

// GetUserRecordsHandler 获取当前用户的问答记录（需要认证）
func (app *App) GetUserRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从上下文获取用户信息
	user, ok := getUserFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "需要用户认证"})
		return
	}

	userID := getUserID(user)
	if userID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的用户ID"})
		return
	}

	records, err := app.qaStorage.GetRecordsByUserID(userID)
	if err != nil {
		log.Printf("获取用户记录失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取记录失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "获取用户记录成功",
		"data":    records,
		"status":  "success",
	})
}

// HealthHandler 健康检查处理器
func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取LLM Provider信息
	providerInfo := app.llmClient.GetProviderInfo()

	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"services": map[string]string{
			"database": "ok",
			"llm":      "ok",
		},
		"llm_provider": providerInfo,
	}

	json.NewEncoder(w).Encode(health)
}

// 辅助函数：从上下文获取用户信息
func getUserFromContext(r *http.Request) (interface{}, bool) {
	user := r.Context().Value("user")
	return user, user != nil
}

// 辅助函数：获取用户ID
func getUserID(user interface{}) int {
	if userMap, ok := user.(map[string]interface{}); ok {
		if id, ok := userMap["id"].(int); ok {
			return id
		}
	}
	return 0
}

// AskStreamHandler 流式提问处理器 - 支持SSE实时响应
func (app *App) AskStreamHandler(w http.ResponseWriter, r *http.Request) {
	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // 禁用nginx缓冲

	// 设置CORS头部（确保SSE请求的跨域支持）
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	question := r.URL.Query().Get("prompt")
	if question == "" {
		app.writeSSEError(w, "缺少prompt参数")
		return
	}

	log.Printf("收到流式问题: %s", question)

	// 检查是否支持流式聊天
	if !app.llmClient.SupportsStreaming() {
		app.writeSSEError(w, "当前LLM Provider不支持流式聊天")
		return
	}

	// 获取用户ID（如果已认证）
	var userID *int
	if user, ok := getUserFromContext(r); ok {
		if id := getUserID(user); id > 0 {
			userID = &id
			log.Printf("认证用户流式提问: ID %d", id)
		}
	} else {
		log.Printf("匿名用户流式提问")
	}

	// 1. 保存问题到数据库
	recordID, err := app.qaStorage.SaveQuestion(question, userID)
	if err != nil {
		log.Printf("保存问题失败: %v", err)
		app.writeSSEError(w, "保存问题失败")
		return
	}

	// 获取flusher用于立即发送数据
	flusher, ok := w.(http.Flusher)
	if !ok {
		app.writeSSEError(w, "不支持流式响应")
		return
	}

	// 发送开始事件
	app.writeSSEData(w, map[string]interface{}{
		"type":      "start",
		"record_id": recordID,
		"question":  question,
		"user_id":   userID,
	})
	log.Printf("发送开始事件，记录ID: %d", recordID)
	flusher.Flush() // 立即发送开始事件

	// 2. 调用LLM流式接口
	ctx := r.Context()
	responseChan, errorChan, err := app.llmClient.ChatCompletionStream(ctx, question)
	if err != nil {
		log.Printf("启动流式聊天失败: %v", err)
		app.writeSSEError(w, "启动流式聊天失败")
		return
	}

	var fullAnswer strings.Builder

	log.Printf("开始监听流式响应...")

	// 3. 处理流式响应
	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				// 流式响应结束
				finalAnswer := fullAnswer.String()
				log.Printf("流式响应结束，最终答案长度: %d", len(finalAnswer))

				// 更新数据库中的答案
				if err := app.qaStorage.UpdateAnswer(recordID, finalAnswer); err != nil {
					log.Printf("更新答案失败: %v", err)
				} else {
					log.Printf("答案更新成功，ID: %d", recordID)
				}

				// 发送结束事件
				app.writeSSEData(w, map[string]interface{}{
					"type":      "end",
					"record_id": recordID,
					"answer":    finalAnswer,
				})
				log.Printf("发送结束事件")
				flusher.Flush()

				log.Printf("流式问答完成，ID: %d", recordID)
				return
			}

			log.Printf("收到流式响应: %+v", resp)

			// 处理流式数据
			if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
				log.Printf("处理Delta数据: %+v", resp.Choices[0].Delta)
				if resp.Choices[0].Delta.Content != nil {
					if content, ok := resp.Choices[0].Delta.Content.(string); ok && content != "" {
						log.Printf("收到增量内容: %s", content)
						fullAnswer.WriteString(content)

						// 发送增量数据
						app.writeSSEData(w, map[string]interface{}{
							"type":    "delta",
							"content": content,
						})
						log.Printf("发送增量数据: %s", content)

						flusher.Flush()
					} else {
						log.Printf("Delta Content为空或类型不匹配: %v", resp.Choices[0].Delta.Content)
					}
				} else {
					log.Printf("Delta Content为nil")
				}
			} else {
				log.Printf("没有有效的Choices或Delta数据")
			}

		case err, ok := <-errorChan:
			if ok && err != nil {
				log.Printf("流式响应错误: %v", err)

				// 保存错误信息到数据库
				app.qaStorage.UpdateAnswer(recordID, "抱歉，AI服务出现错误")

				app.writeSSEError(w, fmt.Sprintf("流式响应错误: %v", err))
				return
			}

		case <-ctx.Done():
			log.Printf("客户端断开连接")
			return
		}
	}
}

// writeSSEData 写入SSE数据
func (app *App) writeSSEData(w http.ResponseWriter, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
}

// writeSSEError 写入SSE错误
func (app *App) writeSSEError(w http.ResponseWriter, message string) {
	errorData := map[string]interface{}{
		"type":  "error",
		"error": message,
	}
	app.writeSSEData(w, errorData)
}
