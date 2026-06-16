package data

import (
	"community/config"
	"community/models"

	"gorm.io/gorm"
)

// ─── 基础查询 ───

// FindUserByID 按主键查用户
func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByIDWithRoles 按主键查用户（带角色）
func FindUserByIDWithRoles(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByUsername 按用户名查用户
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByUsernameWithRoles 按用户名查用户（带角色）
func FindUserByUsernameWithRoles(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ─── 写入 ───

// CreateUser 创建用户
func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// DeleteUser 删除用户（仅删 users 表，关联数据需单独清理）
func DeleteUser(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

// ─── 列表 ───

// ListUsers 获取用户列表（支持搜索，带角色）
func ListUsers(search string) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{}).Preload("Roles")
	if search != "" {
		query = query.Where("username LIKE ?", "%"+search+"%")
	}
	err := query.Order("id ASC").Find(&users).Error
	return users, err
}

// ─── 更新 ───

// UpdateUser 按 map 更新用户字段
func UpdateUser(id uint, updates map[string]interface{}) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// ─── 统计 ───

// CountPostsByUser 统计用户发帖数
func CountPostsByUser(userID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountCommentsByUser 统计用户评论数
func CountCommentsByUser(userID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountFollowers 粉丝数（关注此用户的人数）
func CountFollowers(userID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Follow{}).Where("followee_id = ?", userID).Count(&count).Error
	return count, err
}

// CountFollowing 关注数（此用户关注的人数）
func CountFollowing(userID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

// ─── 关注关系 ───

// IsFollowing 检查 follower 是否关注了 followee
func IsFollowing(followerID, followeeID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}

// ─── 角色 ───

// SetUserRoles 清空旧角色并设置新角色
func SetUserRoles(userID uint, roleNames []string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 清空
		if err := tx.Model(&models.User{ID: userID}).Association("Roles").Clear(); err != nil {
			return err
		}
		// 查找并追加
		var roles []models.Role
		for _, name := range roleNames {
			var role models.Role
			if err := tx.Where("name = ?", name).First(&role).Error; err == nil {
				roles = append(roles, role)
			}
		}
		return tx.Model(&models.User{ID: userID}).Association("Roles").Append(roles)
	})
}

// AppendUserRole 追加单个角色
func AppendUserRole(userID uint, roleName string) error {
	var role models.Role
	if err := config.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}
	return config.DB.Model(&models.User{ID: userID}).Association("Roles").Append(&role)
}

// ─── 级联清理 ───

// DeleteUserContent 删除某用户所有内容（帖、评论、私信、关注），在事务内执行
func DeleteUserContent(userID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return deleteUserContentTx(tx, userID)
	})
}

// DeleteUserCascade 删除用户 + 级联清理所有内容，单个事务保证原子性
func DeleteUserCascade(userID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := deleteUserContentTx(tx, userID); err != nil {
			return err
		}
		return tx.Delete(&models.User{}, userID).Error
	})
}

func deleteUserContentTx(tx *gorm.DB, userID uint) error {
	if err := tx.Where("user_id = ?", userID).Delete(&models.Post{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", userID).Delete(&models.Comment{}).Error; err != nil {
		return err
	}
	if err := tx.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Delete(&models.Message{}).Error; err != nil {
		return err
	}
	if err := tx.Where("follower_id = ? OR followee_id = ?", userID, userID).Delete(&models.Follow{}).Error; err != nil {
		return err
	}
	return nil
}
