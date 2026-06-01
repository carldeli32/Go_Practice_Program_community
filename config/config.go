package config

import (
	"fmt"

	"community/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const dsn = "root:Asd12714052458!@tcp(127.0.0.1:3306)/community?charset=utf8mb4&parseTime=True"

func Init() {
	fmt.Println("🚀 配置加载完成")
}

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ 数据库连接失败: " + err.Error())
	}
	fmt.Println("✅ 数据库连接成功")

	// 自动建表
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Message{}, &models.Follow{}, &models.Announcement{}, &models.Thread{})

	// 种子：创建超级管理员 root，密码 iamking
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "root").Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("iamking"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Username: "root",
			Password: string(hashed),
			IsAdmin:  true,
			Motto:    "系统管理员",
		})
		fmt.Println("👑 超级管理员 root 已创建")
	}
}
