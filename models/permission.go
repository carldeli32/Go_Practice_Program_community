package models

// RolePermissions 角色 → 权限列表
var RolePermissions = map[string][]string{
	"super_admin": {
		// 内容
		"post.create", "post.edit_own", "post.delete_own",
		"comment.create", "comment.edit_own", "comment.delete_own",
		// 版内管理
		"post.manage_category", "comment.manage_category",
		"post.pin", "mute.category", "ban.category", "announce.category",
		// 全站管理
		"post.manage_any", "comment.manage_any",
		"mute.any", "ban.any", "announce.any",
		"user.list",
		// 超管专属
		"user.create", "user.delete", "role.assign", "category.manage",
	},
	"admin": {
		"post.create", "post.edit_own", "post.delete_own",
		"comment.create", "comment.edit_own", "comment.delete_own",
		"post.manage_category", "comment.manage_category",
		"post.pin", "mute.category", "ban.category", "announce.category",
		"post.manage_any", "comment.manage_any",
		"mute.any", "ban.any", "announce.any",
		"user.list",
	},
	"moderator": {
		"post.create", "post.edit_own", "post.delete_own",
		"comment.create", "comment.edit_own", "comment.delete_own",
		"post.manage_category", "comment.manage_category",
		"post.pin", "mute.category", "ban.category", "announce.category",
	},
	"user": {
		"post.create", "post.edit_own", "post.delete_own",
		"comment.create", "comment.edit_own", "comment.delete_own",
	},
}

// HasPerm 检查角色列表是否包含指定权限
func HasPerm(roles []string, perm string) bool {
	for _, role := range roles {
		for _, p := range RolePermissions[role] {
			if p == perm {
				return true
			}
		}
	}
	return false
}
