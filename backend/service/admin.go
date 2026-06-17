package service

import (
	"community/backend/data"
	"community/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// ─── 用户管理 ───

// BanUser 封禁用户（管理员保护），返回用户名
func BanUser(targetID uint) (string, error) {
	user, err := data.FindUserByIDWithRoles(targetID)
	if err != nil {
		return "", ErrNotFound("用户")
	}
	if user.IsAdminLike() {
		return "", ErrCannotBanAdmin
	}
	return user.Username, data.UpdateUser(targetID, map[string]interface{}{"is_banned": true})
}

// UnbanUser 解封用户，返回用户名
func UnbanUser(targetID uint) (string, error) {
	user, err := data.FindUserByID(targetID)
	if err != nil {
		return "", ErrNotFound("用户")
	}
	return user.Username, data.UpdateUser(targetID, map[string]interface{}{"is_banned": false})
}

// AdminCreateUser 管理员创建用户（带默认 user 角色）
func AdminCreateUser(username, password string) (*models.User, error) {
	// 查重
	if _, err := data.FindUserByUsername(username); err == nil {
		return nil, ErrUserExists
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{Username: username, Password: string(hashed)}
	if err := data.CreateUser(user); err != nil {
		return nil, ErrDBOpFail
	}

	// 赋予默认 user 角色
	_ = data.AppendUserRole(user.ID, "user")
	return user, nil
}

// AdminDeleteUser 删除用户（管理员保护 + 级联清理，单事务原子操作）
func AdminDeleteUser(targetID uint) error {
	user, err := data.FindUserByIDWithRoles(targetID)
	if err != nil {
		return ErrNotFound("用户")
	}
	if user.IsAdminLike() {
		return ErrCannotDeleteAdmin
	}

	return data.DeleteUserCascade(targetID)
}

// AssignRoles 给用户设置角色
func AssignRoles(targetID uint, roleNames []string) error {
	if _, err := data.FindUserByID(targetID); err != nil {
		return ErrNotFound("用户")
	}
	return data.SetUserRoles(targetID, roleNames)
}

// ─── 公告 ───

// SetAnnouncement 设置公告（覆盖同名分类旧公告）
func SetAnnouncement(content string, categoryID uint) error {
	_ = data.DeleteAnnouncementsByCategory(categoryID)
	return data.CreateAnnouncement(&models.Announcement{Content: content, CategoryID: categoryID})
}

// DeleteAnnouncement 删除公告
func DeleteAnnouncement(categoryID uint) error {
	return data.DeleteAnnouncementsByCategory(categoryID)
}
