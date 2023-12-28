package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

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

func InitRedis() {

	//数据库连接
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	//pong, err := Red.Ping().Result()
	//if err != nil {
	//	fmt.Println("init redis .....", err)
	//} else {
	//	fmt.Println("Redis inited ......", pong)
	//}
}

const (
	PublishKey = "websocket"
)

// Publish发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish ......", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe ......", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe ......", msg.Payload)
	return msg.Payload, err
}
