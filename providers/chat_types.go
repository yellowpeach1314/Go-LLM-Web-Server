package providers

// Message 消息对象
type Message struct {
	Role             string      `json:"role"`
	Content          interface{} `json:"content,omitempty"`
	ReasoningContent string      `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID       string      `json:"tool_call_id,omitempty"`
	Name             string      `json:"name,omitempty"`
}

// ContentPart 消息内容部分
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL 图片URL
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
	Index    int          `json:"index,omitempty"`
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool 工具定义
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function 函数定义
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Strict      bool                   `json:"strict,omitempty"`
}

// ToolChoice 工具选择
type ToolChoice struct {
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

// ChatCompletionRequest 聊天完成请求
type ChatCompletionRequest struct {
	Model             string                 `json:"model"`
	Messages          []Message              `json:"messages"`
	Tools             []Tool                 `json:"tools,omitempty"`
	ToolChoice        interface{}            `json:"tool_choice,omitempty"`
	Temperature       *float64               `json:"temperature,omitempty"`
	TopP              *float64               `json:"top_p,omitempty"`
	N                 *int                   `json:"n,omitempty"`
	Stream            bool                   `json:"stream,omitempty"`
	StreamOptions     map[string]interface{} `json:"stream_options,omitempty"`
	Stop              interface{}            `json:"stop,omitempty"`
	MaxTokens         *int                   `json:"max_tokens,omitempty"`
	PresencePenalty   *float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64               `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]interface{} `json:"logit_bias,omitempty"`
	ResponseFormat    map[string]interface{} `json:"response_format,omitempty"`
	Seed              *int                   `json:"seed,omitempty"`
	ParallelToolCalls *bool                  `json:"parallel_tool_calls,omitempty"`
	User              string                 `json:"user,omitempty"`
	ExtraParams       map[string]interface{} `json:"-"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
}

// Choice 选择
type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason string   `json:"finish_reason"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionStreamResponse 流式响应
type ChatCompletionStreamResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
}
