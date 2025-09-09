package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// 全局验证器实例
var validate = validator.New()

// UserStorageService 用户存储服务接口
type UserStorageService interface {
	CreateUser(username, email, password string) (interface{}, error)
	ValidateUser(username, password string) (interface{}, error)
	UpdateUserLastLogin(userID int) error
	GetAllUsers() (interface{}, error)
}

// AuthHandlers 认证处理器结构体
type AuthHandlers struct {
	userStorage UserStorageService
	jwtService  *JWTService
}

// NewAuthHandlers 创建认证处理器实例
func NewAuthHandlers(userStorage UserStorageService, jwtService *JWTService) *AuthHandlers {
	return &AuthHandlers{
		userStorage: userStorage,
		jwtService:  jwtService,
	}
}

// RegisterHandler 用户注册处理器
func (ah *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "无效的JSON格式")
		return
	}

	// 验证请求参数
	if err := validate.Struct(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	log.Printf("用户注册请求: %s (%s)", req.Username, req.Email)

	// 创建用户
	user, err := ah.userStorage.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("用户注册失败: %v", err)
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	// 生成JWT token
	token, err := ah.generateTokenForUser(user)
	if err != nil {
		log.Printf("生成JWT失败: %v", err)
		respondWithError(w, http.StatusInternalServerError, "生成认证令牌失败")
		return
	}

	// 返回用户信息和token
	response := UserLoginResponse{
		User:  user,
		Token: token,
	}

	log.Printf("用户注册成功: %s", req.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "注册成功",
		"data":    response,
		"status":  "success",
	})
}

// LoginHandler 用户登录处理器
func (ah *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "无效的JSON格式")
		return
	}

	// 验证请求参数
	if err := validate.Struct(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	log.Printf("用户登录请求: %s", req.Username)

	// 验证用户凭据
	user, err := ah.userStorage.ValidateUser(req.Username, req.Password)
	if err != nil {
		log.Printf("用户登录失败: %v", err)
		respondWithError(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 更新最后登录时间
	if userID := ah.getUserID(user); userID > 0 {
		if err := ah.userStorage.UpdateUserLastLogin(userID); err != nil {
			log.Printf("更新用户最后登录时间失败: %v", err)
		}
	}

	// 生成JWT token
	token, err := ah.generateTokenForUser(user)
	if err != nil {
		log.Printf("生成JWT失败: %v", err)
		respondWithError(w, http.StatusInternalServerError, "生成认证令牌失败")
		return
	}

	// 返回用户信息和token
	response := UserLoginResponse{
		User:  user,
		Token: token,
	}

	log.Printf("用户登录成功: %s", req.Username)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "登录成功",
		"data":    response,
		"status":  "success",
	})
}

// ProfileHandler 获取用户资料处理器（需要认证）
func (ah *AuthHandlers) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从上下文获取用户信息
	user, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未找到用户信息")
		return
	}

	// 返回用户资料
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "获取用户资料成功",
		"data":    user,
		"status":  "success",
	})
}

// RefreshTokenHandler 刷新token处理器（需要认证）
func (ah *AuthHandlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从上下文获取用户信息
	user, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未找到用户信息")
		return
	}

	// 生成新的JWT token
	token, err := ah.generateTokenForUser(user)
	if err != nil {
		log.Printf("刷新JWT失败: %v", err)
		respondWithError(w, http.StatusInternalServerError, "生成认证令牌失败")
		return
	}

	log.Printf("用户刷新token成功")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "token刷新成功",
		"data": map[string]string{
			"token": token,
		},
		"status": "success",
	})
}

// GetUsersHandler 获取用户列表处理器（管理员功能，需要认证）
func (ah *AuthHandlers) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从上下文获取当前用户信息
	currentUser, ok := GetUserFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未找到用户信息")
		return
	}

	username := ah.getUserField(currentUser, "Username")
	log.Printf("管理员请求用户列表: %s", username)

	// 获取所有用户
	users, err := ah.userStorage.GetAllUsers()
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		respondWithError(w, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "获取用户列表成功",
		"data":    users,
		"status":  "success",
	})
}

// LogoutHandler 用户登出处理器（客户端处理，服务端记录日志）
func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从上下文获取用户信息（如果有的话）
	if user, ok := GetUserFromContext(r); ok {
		username := ah.getUserField(user, "Username")
		log.Printf("用户登出: %s", username)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "登出成功",
		"status":  "success",
		"note":    "请在客户端删除存储的token",
	})
}

// generateTokenForUser 为用户生成token的辅助方法
func (ah *AuthHandlers) generateTokenForUser(user interface{}) (string, error) {
	userID := ah.getUserID(user)
	username := ah.getUserField(user, "Username")
	email := ah.getUserField(user, "Email")

	return ah.jwtService.GenerateToken(userID, username, email)
}

// getUserID 获取用户ID的辅助方法
func (ah *AuthHandlers) getUserID(user interface{}) int {
	if userMap, ok := user.(map[string]interface{}); ok {
		if id, ok := userMap["id"].(int); ok {
			return id
		}
	}
	// 使用反射或类型断言来获取ID
	return 0
}

// getUserField 获取用户字段的辅助方法
func (ah *AuthHandlers) getUserField(user interface{}, field string) string {
	if userMap, ok := user.(map[string]interface{}); ok {
		if value, ok := userMap[field].(string); ok {
			return value
		}
	}
	return ""
}
