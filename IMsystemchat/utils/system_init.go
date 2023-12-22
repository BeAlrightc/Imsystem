package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// 初始化配置文件操作
func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	viper.ReadInConfig()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app0 inited .....")
	//fmt.Println("config mysql:", viper.Get("mysql"))
}

func InitMYSQL() {
	newLogger := logger.New(
		//自定义日志模板，打印SQL语句
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢ssql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	//数据库连接
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("config mysql:", viper.Get("mysql"))
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
	//看完了swagger的整合操作
}
