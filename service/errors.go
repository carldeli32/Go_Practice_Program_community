package service

import "fmt"

// AppError 业务层统一错误，携带 HTTP 状态码
type AppError struct {
	Code    int    // HTTP 状态码
	Message string // 用户可见的错误信息
}

func (e *AppError) Error() string {
	return e.Message
}

// ─── 构造动态错误 ───

func ErrMuted(until string) *AppError {
	return &AppError{Code: 403, Message: "你已被禁言，解禁时间: " + until}
}

func ErrMutedCategory(until string) *AppError {
	return &AppError{Code: 403, Message: "你已被该版块禁言，解禁时间: " + until}
}

func ErrNotFound(what string) *AppError {
	return &AppError{Code: 404, Message: what + "不存在"}
}

// ─── 静态哨兵错误 ───

var (
	ErrUserExists        = &AppError{Code: 409, Message: "用户名已被注册"}
	ErrBadCredentials    = &AppError{Code: 401, Message: "用户名或密码错误"}
	ErrBanned            = &AppError{Code: 403, Message: "账号已被封禁，请联系管理员"}
	ErrForbidden         = &AppError{Code: 403, Message: "无权操作"}
	ErrCannotSelfFollow  = &AppError{Code: 400, Message: "不能关注自己"}
	ErrCannotSelfMessage = &AppError{Code: 400, Message: "不能给自己发私信"}
	ErrCannotSelfThread  = &AppError{Code: 400, Message: "不能和自己对话"}
	ErrAlreadyFollowing  = &AppError{Code: 409, Message: "已关注该用户"}
	ErrNotFollowing      = &AppError{Code: 404, Message: "未关注该用户"}
	ErrCannotBanAdmin    = &AppError{Code: 400, Message: "不能封禁管理员"}
	ErrCannotDeleteAdmin = &AppError{Code: 400, Message: "不能删除管理员"}
	ErrCannotDeleteRoot  = &AppError{Code: 400, Message: "默认分类「综合讨论」不可删除"}
	ErrNoUpdateContent   = &AppError{Code: 400, Message: "没有要更新的内容"}
	ErrCategoryExists    = &AppError{Code: 409, Message: "分类名已存在"}
	ErrPasswordHashFail  = &AppError{Code: 500, Message: "密码加密失败"}
	ErrTokenGenFail      = &AppError{Code: 500, Message: "Token 生成失败"}
	ErrDBOpFail          = &AppError{Code: 500, Message: "操作失败"}

	// 注册失败
	ErrRegisterFail = func(msg string) *AppError {
		return &AppError{Code: 500, Message: msg}
	}
)

// ToHTTP 辅助：将 AppError 拆成 HTTP 状态码 + 消息
func ToHTTP(err error) (int, string) {
	if ae, ok := err.(*AppError); ok {
		return ae.Code, ae.Message
	}
	return 500, fmt.Sprintf("内部错误: %v", err)
}
