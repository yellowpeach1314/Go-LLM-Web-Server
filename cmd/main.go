package main

import (
	"log"
	"net/http"

	"go-base-web-server/internal/auth"
	"go-base-web-server/internal/config"
	"go-base-web-server/internal/handlers"
	"go-base-web-server/internal/llm"
	"go-base-web-server/internal/middleware"
	"go-base-web-server/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	qaStorage, err := storage.NewQAStorage(cfg.DBPath)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer qaStorage.Close()

	// 初始化用户存储
	userStorage := storage.NewUserStorage(qaStorage.GetDB())
	if err := userStorage.InitUserTables(); err != nil {
		log.Fatalf("初始化用户表失败: %v", err)
	}

	// 初始化LLM客户端
	llmClient := llm.NewClient(llm.Config{
		Provider: cfg.LLMProvider,
		APIKey:   cfg.LLMAPIKey,
		APIURL:   cfg.LLMAPIURL,
		Model:    cfg.LLMModel,
	})
	if err := llmClient.CheckConnection(); err != nil {
		log.Printf("LLM连接检查失败: %v", err)
	}

	// 初始化JWT服务
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// 创建应用实例
	app := handlers.NewApp(qaStorage, llmClient)

	// 创建认证处理器
	authHandlers := auth.NewAuthHandlers(userStorage, jwtService)

	// 创建路由器
	r := mux.NewRouter()

	// 添加全局中间件
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// 公开路由（不需要认证）
	r.HandleFunc("/", app.HomeHandler).Methods("GET")
	r.HandleFunc("/api/health", app.HealthHandler).Methods("GET")

	// 认证相关路由
	r.HandleFunc("/api/auth/register", authHandlers.RegisterHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", authHandlers.LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/logout", authHandlers.LogoutHandler).Methods("POST", "OPTIONS")

	// 可选认证路由（支持匿名和认证用户）
	optionalAuth := r.PathPrefix("/api").Subrouter()
	optionalAuth.Use(auth.OptionalAuthMiddleware(jwtService, userStorage))
	optionalAuth.HandleFunc("/ask", app.AskHandler).Methods("GET")
	optionalAuth.HandleFunc("/ask/stream", app.AskStreamHandler).Methods("GET")
	optionalAuth.HandleFunc("/records", app.GetRecordsHandler).Methods("GET")
	optionalAuth.HandleFunc("/records/{id:[0-9]+}", app.GetRecordHandler).Methods("GET")

	// 需要认证的路由
	authRequired := r.PathPrefix("/api/user").Subrouter()
	authRequired.Use(auth.AuthMiddleware(jwtService, userStorage))
	authRequired.HandleFunc("/profile", authHandlers.ProfileHandler).Methods("GET", "OPTIONS")
	authRequired.HandleFunc("/refresh-token", authHandlers.RefreshTokenHandler).Methods("POST", "OPTIONS")
	authRequired.HandleFunc("/records", app.GetUserRecordsHandler).Methods("GET", "OPTIONS")
	authRequired.HandleFunc("/users", authHandlers.GetUsersHandler).Methods("GET", "OPTIONS") // 管理员功能

	// 服务器配置
	port := ":" + cfg.Port

	log.Printf("🚀 LLM问答系统后端API启动成功")
	log.Printf("📍 API服务器地址: http://localhost%s", port)
	log.Printf("💾 数据库路径: %s", cfg.DBPath)
	log.Printf("🤖 LLM模式: %s", cfg.GetLLMMode())
	log.Printf("🔐 用户认证: 已启用")
	log.Printf("🌐 前端项目: ./frontend/ (独立运行)")
	log.Printf("📚 API文档: http://localhost%s", port)

	// 打印API路由信息
	printAPIRoutes()

	// 启动服务器
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// printAPIRoutes 打印API路由信息
func printAPIRoutes() {
	log.Println("📋 API路由列表:")
	log.Println("   公开路由:")
	log.Println("     GET  /                    - API信息")
	log.Println("     GET  /api/health          - 健康检查")
	log.Println("     POST /api/auth/register   - 用户注册")
	log.Println("     POST /api/auth/login      - 用户登录")
	log.Println("     POST /api/auth/logout     - 用户登出")
	log.Println("   可选认证路由:")
	log.Println("     GET  /api/ask             - 智能问答（支持匿名）")
	log.Println("     GET  /api/ask/stream      - 流式智能问答（支持匿名）- SSE")
	log.Println("     GET  /api/records         - 获取所有记录")
	log.Println("     GET  /api/records/{id}    - 获取特定记录")
	log.Println("   需要认证路由:")
	log.Println("     GET  /api/user/profile    - 获取用户资料")
	log.Println("     POST /api/user/refresh-token - 刷新token")
	log.Println("     GET  /api/user/records    - 获取用户记录")
	log.Println("     GET  /api/user/users      - 获取用户列表（管理员）")
}
