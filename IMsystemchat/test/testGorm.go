package mains

import (
	"IMsystemchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:000000@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connnect database")
	}

	//迁移schema 就是根据相关的结狗提生成相应的mysql数据表操作
	//db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	//db.AutoMigrate(&models.Contact{})
	//db.AutoMigrate(&models.GroupBasic{})

	//create
	//user := &models.UserBasic{}
	//user.Name = "小刘"
	//user.PassWord = "123"
	//db.Create(user)
	//
	////read
	//fmt.Println(db.First(user, 1))
	//
	////update
	//db.Model(user).Update("PassWord", "123456")

}
