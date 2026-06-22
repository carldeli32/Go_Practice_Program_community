package config

import (
	"fmt"

	"community/backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	LoadEnv()
	fmt.Println("🚀 配置加载完成")
}

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN()), &gorm.Config{})
	if err != nil {
		panic("❌ 数据库连接失败: " + err.Error())
	}
	fmt.Println("✅ 数据库连接成功")

	// 自动建表
	DB.AutoMigrate(
		&models.User{}, &models.Post{}, &models.Comment{}, &models.Message{},
		&models.Follow{}, &models.Announcement{}, &models.Thread{},
		&models.Category{}, &models.Role{}, &models.ModeratorCategory{},
		&models.Mute{}, &models.CategoryBan{},
	)

	// 种子：默认分类
	seedCategories()

	// 种子：默认角色
	seedRoles()

	// 种子：超级管理员 root
	seedRootUser()
}

func seedCategories() {
	categories := []models.Category{
		{Name: "综合讨论", Description: "各类话题，畅所欲言"},
		{Name: "技术交流", Description: "编程、数码、科技相关讨论"},
		{Name: "军事纵横", Description: "军事历史、装备、战略讨论"},
		{Name: "历史长廊", Description: "古今中外，谈史论道"},
		{Name: "文学艺术", Description: "诗词歌赋，文艺创作"},
		{Name: "生活杂谈", Description: "日常趣事，生活分享"},
	}
	for _, cat := range categories {
		var existing models.Category
		if err := DB.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			DB.Create(&cat)
			fmt.Printf("📂 分类「%s」已创建\n", cat.Name)
		}
	}
}

func seedRoles() {
	roles := []models.Role{
		{Name: "super_admin", Description: "超级管理员"},
		{Name: "admin", Description: "管理员"},
		{Name: "moderator", Description: "版主"},
		{Name: "user", Description: "普通用户"},
	}
	for _, role := range roles {
		var existing models.Role
		if err := DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			DB.Create(&role)
			fmt.Printf("🔑 角色「%s」已创建\n", role.Name)
		}
	}
}

func seedRootUser() {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "root").Count(&count)
	if count > 0 {
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("iamking"), bcrypt.DefaultCost)
	user := models.User{
		Username: "root",
		Password: string(hashed),
		Motto:    "系统管理员",
	}
	DB.Create(&user)

	// 绑定 super_admin + admin 角色
	var superAdmin, admin models.Role
	DB.Where("name = ?", "super_admin").First(&superAdmin)
	DB.Where("name = ?", "admin").First(&admin)
	DB.Model(&user).Association("Roles").Append(&superAdmin, &admin)

	fmt.Println("👑 超级管理员 root 已创建（角色: super_admin, admin）")
}
