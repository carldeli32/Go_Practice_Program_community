package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 通过 gorm.DB 来操作数据库
var DB *gorm.DB

// 数据库连接信息
// 格式：用户名:密码@tcp(地址:端口)/数据库名?charset=utf8mb4
const dsn = "root:Asd12714052458!@tcp(127.0.0.1:3306)/community?charset=utf8mb4&parseTime=True"

func Init() {
	fmt.Println("配置加载完成")
}

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("数据库连接成功")
}
