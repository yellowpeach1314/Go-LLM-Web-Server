# 🚀 前端开发者API使用指南 v2.0

本文档专为前端开发者提供后端API的完整使用指南。后端项目已升级到v2.0版本，新增了用户认证系统和模块化架构。

## 📋 重要变更说明

### 🆕 新增功能
- **用户认证系统**：支持用户注册、登录、JWT认证
- **权限控制**：区分匿名用户和认证用户访问
- **用户管理**：个人资料、用户记录、管理员功能
- **Token刷新**：自动续期机制

### 🔄 接口变更
- 所有API响应格式统一化
- 新增认证相关接口
- 部分接口支持可选认证
- 错误处理标准化

## 🌐 服务器信息

- **基础URL**: `http://localhost:8080`
- **API版本**: v2.0
- **认证方式**: JWT Bearer Token
- **内容类型**: `application/json`

## 📚 API接口详解

### 1. 公开接口（无需认证）

#### 1.1 获取API信息
```http
GET /
```

**响应示例：**
```json
{
  "name": "LLM 问答系统 API",
  "version": "2.0.0",
  "description": "现代化的前后端分离问答系统后端API - 支持用户认证",
  "endpoints": { ... },
  "authentication": {
    "type": "Bearer Token (JWT)",
    "header": "Authorization: Bearer <token>"
  }
}
```

#### 1.2 健康检查
```http
GET /api/health
```

**响应示例：**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01 10:00:00",
  "services": {
    "database": "ok",
    "llm": "ok"
  },
  "llm_provider": {
    "provider": "OpenAI Provider",
    "status": "active"
  }
}
```

#### 1.3 用户注册
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "string (必填，3-50字符)",
  "email": "string (必填，有效邮箱格式)",
  "password": "string (必填，6-100字符)"
}
```

**响应示例：**
```json
{
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "api_key": "generated-api-key",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "status": "success"
}
```

#### 1.4 用户登录
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "string (必填)",
  "password": "string (必填)"
}
```

**响应示例：**
```json
{
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "api_key": "generated-api-key",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "status": "success"
}
```

#### 1.5 用户登出
```http
POST /api/auth/logout
```

**响应示例：**
```json
{
  "message": "登出成功",
  "status": "success",
  "note": "请在客户端删除存储的token"
}
```

### 2. 可选认证接口（支持匿名访问）

这些接口支持匿名访问，但如果提供了有效的JWT Token，会记录用户信息。

#### 2.1 智能问答
```http
GET /api/ask?prompt=你的问题
Authorization: Bearer <token> (可选)
```

**查询参数：**
- `prompt` (string, 必填): 要提问的内容

**响应示例：**
```json
{
  "id": 1,
  "question": "你好",
  "answer": "你好！我是AI助手，很高兴为您服务。有什么我可以帮助您的吗？",
  "user_id": 1,
  "status": "success"
}
```

#### 2.2 流式智能问答 (SSE)
```http
GET /api/ask/stream?prompt=你的问题
Authorization: Bearer <token> (可选)
```

**查询参数：**
- `prompt` (string, 必填): 要提问的内容

**响应格式：** Server-Sent Events (SSE)
**Content-Type：** `text/event-stream`

**SSE事件类型：**

1. **开始事件**
```json
data: {
  "type": "start",
  "record_id": 1,
  "question": "你的问题",
  "user_id": 1
}
```

2. **增量数据事件**
```json
data: {
  "type": "delta",
  "content": "AI回答的片段"
}
```

3. **结束事件**
```json
data: {
  "type": "end",
  "record_id": 1,
  "answer": "完整的AI回答"
}
```

4. **错误事件**
```json
data: {
  "type": "error",
  "error": "错误描述"
}
```

**前端使用示例：**
```typescript
function askStreamQuestion(prompt: string): EventSource {
  const token = localStorage.getItem('authToken');
  const url = `/api/ask/stream?prompt=${encodeURIComponent(prompt)}`;
  
  const eventSource = new EventSource(url, {
    headers: token ? { 'Authorization': `Bearer ${token}` } : {}
  });

  let fullAnswer = '';
  
  eventSource.onmessage = function(event: MessageEvent) {
    const data: SSEEvent = JSON.parse(event.data);
    
    switch(data.type) {
      case 'start':
        console.log('开始接收回答，记录ID:', data.record_id);
        fullAnswer = '';
        break;
        
      case 'delta':
        fullAnswer += data.content;
        // 实时更新UI显示
        updateAnswerDisplay(fullAnswer);
        break;
        
      case 'end':
        console.log('回答完成:', data.answer);
        eventSource.close();
        break;
        
      case 'error':
        console.error('流式响应错误:', data.error);
        eventSource.close();
        break;
    }
  };

  eventSource.onerror = function(event: Event) {
    console.error('SSE连接错误:', event);
    eventSource.close();
  };
  
  return eventSource;
}

function updateAnswerDisplay(answer: string): void {
  // 实现UI更新逻辑
  const answerElement = document.getElementById('answer');
  if (answerElement) {
    answerElement.textContent = answer;
  }
}
```

#### 2.3 获取所有问答记录
```http
GET /api/records
Authorization: Bearer <token> (可选)
```

**响应示例：**
```json
{
  "message": "获取记录成功",
  "data": [
    {
      "id": 1,
      "question": "你好",
      "answer": "你好！我是AI助手...",
      "user_id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:01Z"
    }
  ],
  "status": "success"
}
```

#### 2.4 获取特定问答记录
```http
GET /api/records/{id}
Authorization: Bearer <token> (可选)
```

**路径参数：**
- `id` (integer, 必填): 记录ID

**响应示例：**
```json
{
  "message": "获取记录成功",
  "data": {
    "id": 1,
    "question": "你好",
    "answer": "你好！我是AI助手...",
    "user_id": 1,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:01Z"
  },
  "status": "success"
}
```

### 3. 需要认证的接口

这些接口必须提供有效的JWT Token才能访问。

#### 3.1 获取用户资料
```http
GET /api/user/profile
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "message": "获取用户资料成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "api_key": "generated-api-key",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  "status": "success"
}
```

#### 3.2 刷新Token
```http
POST /api/user/refresh-token
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "message": "token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "status": "success"
}
```

#### 3.3 获取用户问答记录
```http
GET /api/user/records
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "message": "获取用户记录成功",
  "data": [
    {
      "id": 1,
      "question": "你好",
      "answer": "你好！我是AI助手...",
      "user_id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:01Z"
    }
  ],
  "status": "success"
}
```

#### 3.4 获取用户列表（管理员功能）
```http
GET /api/user/users
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "message": "获取用户列表成功",
  "data": [
    {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "api_key": "generated-api-key",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "status": "success"
}
```

## 🔐 认证机制详解

### JWT Token格式
```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "exp": 1640995200,
  "iat": 1640908800,
  "iss": "go-base-web-server"
}
```

### 认证流程
1. **用户注册/登录** → 获取JWT Token
2. **存储Token** → 保存到localStorage或sessionStorage
3. **请求头添加** → `Authorization: Bearer <token>`
4. **自动刷新** → Token过期前调用刷新接口

### Token使用示例
```typescript
// 存储token
localStorage.setItem('authToken', response.data.token);

// 使用token发送请求
const token: string | null = localStorage.getItem('authToken');
fetch('/api/user/profile', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

## ❌ 错误处理

### 标准错误响应格式
```json
{
  "error": "错误描述信息",
  "status": "error",
  "code": "HTTP状态码"
}
```

### 常见错误码
- `400` - 请求参数错误
- `401` - 未认证或Token无效
- `403` - 权限不足
- `404` - 资源不存在
- `409` - 资源冲突（如用户名已存在）
- `500` - 服务器内部错误

### 错误处理示例
```typescript
interface LoginRequest {
  username: string;
  password: string;
}

async function loginUser(username: string, password: string): Promise<APIResponse<AuthResponse>> {
  try {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password } as LoginRequest)
    });

    if (!response.ok) {
      const error: APIResponse = await response.json();
      throw error;
    }

    return await response.json();
  } catch (error: any) {
    console.error('登录失败:', error.error);
    throw error;
  }
}
```

## � TypeScript 类型定义

```typescript
// API响应类型定义
interface APIResponse<T = any> {
  message?: string;
  data?: T;
  status: 'success' | 'error';
  error?: string;
}

interface User {
  id: number;
  username: string;
  email: string;
  api_key: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface QARecord {
  id: number;
  question: string;
  answer: string;
  user_id: number | null;
  created_at: string;
  updated_at: string;
}

interface AuthResponse {
  user: User;
  token: string;
}

interface HealthResponse {
  status: string;
  timestamp: string;
  services: {
    database: string;
    llm: string;
  };
  llm_provider: {
    provider: string;
    status: string;
  };
}

// SSE事件类型定义
interface SSEStartEvent {
  type: 'start';
  record_id: number;
  question: string;
  user_id: number | null;
}

interface SSEDeltaEvent {
  type: 'delta';
  content: string;
}

interface SSEEndEvent {
  type: 'end';
  record_id: number;
  answer: string;
}

interface SSEErrorEvent {
  type: 'error';
  error: string;
}

type SSEEvent = SSEStartEvent | SSEDeltaEvent | SSEEndEvent | SSEErrorEvent;

// 回调函数类型
type StreamMessageCallback = (message: {
  type: 'start' | 'delta';
  data?: SSEStartEvent;
  content?: string;
  fullAnswer?: string;
}) => void;

type StreamErrorCallback = (error: string) => void;
type StreamCompleteCallback = (finalAnswer: string) => void;
```

## �💻 前端集成示例

### 1. 用户认证管理类
```typescript
class AuthManager {
  private baseURL: string;
  private token: string | null;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.token = localStorage.getItem('authToken');
  }

  async register(username: string, email: string, password: string): Promise<APIResponse<AuthResponse>> {
    const response = await fetch(`${this.baseURL}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, email, password })
    });
    
    if (response.ok) {
      const data: APIResponse<AuthResponse> = await response.json();
      if (data.data?.token) {
        this.setToken(data.data.token);
      }
      return data;
    }
    throw await response.json();
  }

  async login(username: string, password: string): Promise<APIResponse<AuthResponse>> {
    const response = await fetch(`${this.baseURL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    
    if (response.ok) {
      const data: APIResponse<AuthResponse> = await response.json();
      if (data.data?.token) {
        this.setToken(data.data.token);
      }
      return data;
    }
    throw await response.json();
  }

  async logout(): Promise<void> {
    await fetch(`${this.baseURL}/api/auth/logout`, {
      method: 'POST',
      headers: this.getAuthHeaders()
    });
    this.clearToken();
  }

  setToken(token: string): void {
    this.token = token;
    localStorage.setItem('authToken', token);
  }

  clearToken(): void {
    this.token = null;
    localStorage.removeItem('authToken');
  }

  getAuthHeaders(): Record<string, string> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }
    return headers;
  }

  isAuthenticated(): boolean {
    return !!this.token;
  }

  getToken(): string | null {
    return this.token;
  }
}
```

### 2. API客户端类
```typescript
class APIClient {
  private baseURL: string;
  private auth: AuthManager;

  constructor(authManager: AuthManager, baseURL: string = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.auth = authManager;
  }

  async askQuestion(prompt: string): Promise<{
    id: number;
    question: string;
    answer: string;
    user_id: number | null;
    status: string;
  }> {
    const url = `${this.baseURL}/api/ask?prompt=${encodeURIComponent(prompt)}`;
    const response = await fetch(url, {
      headers: this.auth.getAuthHeaders()
    });
    
    if (response.ok) {
      return await response.json();
    }
    throw await response.json();
  }

  async getRecords(): Promise<APIResponse<QARecord[]>> {
    const response = await fetch(`${this.baseURL}/api/records`, {
      headers: this.auth.getAuthHeaders()
    });
    
    if (response.ok) {
      return await response.json();
    }
    throw await response.json();
  }

  async getUserProfile(): Promise<APIResponse<User>> {
    const response = await fetch(`${this.baseURL}/api/user/profile`, {
      headers: this.auth.getAuthHeaders()
    });
    
    if (response.ok) {
      return await response.json();
    }
    throw await response.json();
  }

  async getUserRecords(): Promise<APIResponse<QARecord[]>> {
    const response = await fetch(`${this.baseURL}/api/user/records`, {
      headers: this.auth.getAuthHeaders()
    });
    
    if (response.ok) {
      return await response.json();
    }
    throw await response.json();
  }

  // 流式问答方法
  askQuestionStream(
    prompt: string, 
    onMessage?: StreamMessageCallback, 
    onError?: StreamErrorCallback, 
    onComplete?: StreamCompleteCallback
  ): EventSource {
    const url = `${this.baseURL}/api/ask/stream?prompt=${encodeURIComponent(prompt)}`;
    
    const eventSource = new EventSource(url);
    
    // 注意：EventSource不支持自定义headers，如需认证可考虑：
    // 1. 将token作为查询参数传递（不推荐，安全性较低）
    // 2. 使用WebSocket替代SSE
    // 3. 在服务端支持cookie认证
    
    let fullAnswer = '';
    
    eventSource.onmessage = function(event: MessageEvent) {
      try {
        const data: SSEEvent = JSON.parse(event.data);
        
        switch(data.type) {
          case 'start':
            onMessage?.({ type: 'start', data });
            fullAnswer = '';
            break;
            
          case 'delta':
            fullAnswer += data.content;
            onMessage?.({ 
              type: 'delta', 
              content: data.content, 
              fullAnswer 
            });
            break;
            
          case 'end':
            onComplete?.(data.answer);
            eventSource.close();
            break;
            
          case 'error':
            onError?.(data.error);
            eventSource.close();
            break;
        }
      } catch (err) {
        onError?.('解析响应数据失败');
        eventSource.close();
      }
    };

    eventSource.onerror = function(event: Event) {
      onError?.('SSE连接错误');
      eventSource.close();
    };
    
    return eventSource;
  }
}
```

### 3. 使用示例
```typescript
// 初始化
const authManager = new AuthManager();
const apiClient = new APIClient(authManager);

// 用户注册
async function registerUser(): Promise<void> {
  try {
    const result = await authManager.register('testuser', 'test@example.com', 'password123');
    console.log('注册成功:', result);
  } catch (error: any) {
    console.error('注册失败:', error.error);
  }
}

// 用户登录
async function loginUser(): Promise<void> {
  try {
    const result = await authManager.login('testuser', 'password123');
    console.log('登录成功:', result);
  } catch (error: any) {
    console.error('登录失败:', error.error);
  }
}

// 提问
async function askQuestion(): Promise<void> {
  try {
    const result = await apiClient.askQuestion('你好，请介绍一下自己');
    console.log('AI回答:', result.answer);
  } catch (error: any) {
    console.error('提问失败:', error.error);
  }
}

// 获取用户资料
async function getUserProfile(): Promise<void> {
  if (authManager.isAuthenticated()) {
    try {
      const profile = await apiClient.getUserProfile();
      console.log('用户资料:', profile.data);
    } catch (error: any) {
      console.error('获取资料失败:', error.error);
    }
  }
}

// 流式提问示例
function startStreamChat(): EventSource {
  const eventSource = apiClient.askQuestionStream(
    '请写一首关于春天的诗',
    // onMessage回调
    (message) => {
      if (message.type === 'start') {
        console.log('开始接收回答...');
        const answerElement = document.getElementById('answer');
        if (answerElement) {
          answerElement.innerHTML = '';
        }
      } else if (message.type === 'delta') {
        // 实时更新显示
        const answerElement = document.getElementById('answer');
        if (answerElement && message.fullAnswer) {
          answerElement.innerHTML = message.fullAnswer;
        }
      }
    },
    // onError回调
    (error: string) => {
      console.error('流式问答错误:', error);
      alert('问答失败: ' + error);
    },
    // onComplete回调
    (finalAnswer: string) => {
      console.log('问答完成:', finalAnswer);
      const answerElement = document.getElementById('answer');
      if (answerElement) {
        answerElement.innerHTML = finalAnswer;
      }
    }
  );

  return eventSource;
}

// 可以手动关闭连接
// eventSource.close();
```

## 🎯 最佳实践建议

### 1. Token管理
- 使用localStorage持久化存储Token
- 实现Token自动刷新机制
- 在Token过期时自动跳转到登录页

### 2. 错误处理
- 统一处理HTTP错误状态码
- 为用户提供友好的错误提示
- 实现重试机制

### 3. 用户体验
- 支持匿名用户访问基础功能
- 为认证用户提供个性化体验
- 实现加载状态和进度提示

### 4. 安全考虑
- 不在URL中传递敏感信息
- 使用HTTPS传输（生产环境）
- 定期刷新Token

### 5. 性能优化
- 缓存不经常变化的数据
- 使用防抖处理用户输入
- 实现分页加载

## 🔧 开发工具推荐

### API测试工具
- **Postman**: 接口测试和文档
- **Insomnia**: 轻量级API客户端
- **curl**: 命令行测试

### 前端框架集成
- **React**: 使用axios或fetch
- **Vue**: 使用axios或Vue Resource
- **Angular**: 使用HttpClient

### 状态管理
- **Redux/Zustand**: React状态管理
- **Vuex/Pinia**: Vue状态管理
- **NgRx**: Angular状态管理

## 📞 技术支持

如有问题，请联系后端开发团队或查看：
- 项目README.md
- API健康检查接口
- 服务器日志

---

**文档版本**: v2.0  
**最后更新**: 2024年1月  
**兼容性**: 支持所有现代浏览器