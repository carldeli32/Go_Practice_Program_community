package middlewares

import (
	"community/backend/data"
	"community/backend/models"
)

// CheckRelation 关系推理：用户是否对资源拥有指定操作权限
//
// 三条路径满足其一即可：
//  1. 直接关系：userID == ownerID（作者本人）
//  2. 全局关系：HasPerm(roles, permAny)（如 post.manage_any）
//  3. 间接关系：HasPerm(roles, permCategory) + moderator OF category（分区版主）
func CheckRelation(userID, ownerID, categoryID uint, roles []string, permAny, permCategory string) bool {
	// 路径 1：作者本人
	if userID == ownerID {
		return true
	}
	// 路径 2：全局管理（admin / super_admin）
	if models.HasPerm(roles, permAny) {
		return true
	}
	// 路径 3：分区管理（moderator of this category）
	if models.HasPerm(roles, permCategory) {
		ok, _ := data.ExistsModeratorCategory(userID, categoryID)
		return ok
	}
	return false
}
