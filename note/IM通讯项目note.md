Gin +webSocket项目实战

# 一、需求分析

## 1，项目目的：

当然是提升自己的coding 技术

## 2.我能获得什么

熟悉开发环境，熟练相关技术栈gin+grom+swagger+logrus auth等中间件。高性能

## 3.核心功能

发送和接受消息，文字表情图片 音频，访客，点对点，群聊，广播，快捷回复，撤回，心跳检测

## 4.技术栈：

前端，后端(websocket,channel/goroutine,gin,template,gorm.sql,nosql)单元测试，日志....



## 5，系统架构

四层：前端、接入层、逻辑层、持久层

![](D:\myfile\GO\project\IMsystem\note\pic\架构.jpg)



## 6.消息发送流程

A>登录>鉴权（游客）>消息类型(群聊广播) >B

# 二、环境搭建

安装号golang

在idea新建一个project

go mod tidy

# 三、代码编写

### 1.引入gorm

http://pkg.go.dev/可以搜索到gorm

https://gorm.io/zh_CN/docs/中文API

首先快速开启代码块

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

1》.用户模块设计

创建models包user_basic.go再写一个struct

```go
package models

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LoginOutTime  uint64 `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

// 这个结构体的函数,返回一个表名
func (table *UserBasic) TableName() string {
	return "user_basic"
}
```

2》.建一个test包，进行数据库的连接实验

```go
package mains

import (
	"IMsystemchat/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:123456@tcp(192.168.20.10:3306)/ginchat?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connnect database")
	}

	//迁移schema
	db.AutoMigrate(&models.UserBasic{})

	//create
	user := &models.UserBasic{}
	user.Name = "xiaoming"
	db.Create(user)

	//read
	fmt.Println(db.First(user, 1))

	//update
	db.Model(user).Update("PassWord", "123456")

}

```

并进行查看数据情况

### 2.引入gin框架

dev搜

```go
go get -u github.com/gin-gonic/gin
```

main.go

```go
package main

import (
	"IMsystemchat/router"
)

func main() {
	r := router.Router()
	r.Run(":8081")
}

```

router包下面的app.go

```go
package router

import (
	"IMsystemchat/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/index", service.GetIndex)

	return r

}

```

service下的index.go

```go
package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome",
	})
}

```

### 3.将数据和请求关联起来

1》在main方法里面初始化配置文件以及数据库


*utils.InitConfig()
*utils.InitMYSQL()

2》在config包里面的app.yml

```yml
mysql:
  dns: root:123456@tcp(192.168.20.10:3306)/ginchat?charset=utf8mb4&parseTime=True

```

3》新建utils包system_ini.go

```go
package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))
}

func InitMYSQL() {
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
    //一下code是用来测试数据库是否连接成功，查找除数据表中的第一条data
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

```

4》在model的user_basic.go下建立getUserList()方法

```go
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

```

5》到service包中userService.go后建立getUserList()方法

```go
package service

import (
	"IMsystemchat/models"
	"github.com/gin-gonic/gin"
	"net/http"
)
//获取表中所有字段并将其json化
func GetUserList(c *gin.Context) {
	//拿到数据
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
```

6》到router里面加上app.go

```go
package router

import (
	"IMsystemchat/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/index", service.GetIndex)
    //加上这一行即可
	r.GET("/user/getUserList", service.GetUserList)
	return r
}

```

4.整合swagger'

download

```go
 go get -u github.com/swaggo/swag/cmd/swag
//应该引入这个包
go get -u github.com/swaggo/gin-swagger
//对swagger进行初始化
swag init 一下
查看项目是否多一个docs目录
然后拉取
改造router app.go
package router

import (
	"IMsystemchat/docs"
	"IMsystemchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	return r

}
然后测试index页面 http://127.0.0.1:8081/swagger/index.html

然后再到service中的index.go中加入相应的注解注释（方法上面加）
package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome !!",
	})
}
以及userservice.go
package service

import (
	"IMsystemchat/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserList
// @Tags 首页
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	//拿到数据
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

swag init 一下


```

test的pic

![](D:\myfile\GO\project\IMsystem\note\pic\swagger.jpg)

![](D:\myfile\GO\project\IMsystem\note\pic\swagger2.jpg)

### 4.日志打印

在util下的system_init.go中修改initMYSQL方法加入自己的logger

```go

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
}

```

# 四.完成用户模块基本的功能

### 1.CRUD

#### -1.显示所有用户

userService.go

```go
// GetUserList
// @Summary 所有用户
// @Tags 首页
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	//拿到数据
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
```

model层下的user_basic

```go
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
```

Router/app.go

```go
r.GET("/user/getUserList", service.GetUserList)
```

#### -2.新增用户

userService.go

```go
// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致！",
		})
		return
	}
	//将密码给user对象
	user.PassWord = password
	models.CreateUser(user) //推入数据库中
	c.JSON(200, gin.H{
		"message": "新增用户成功",
	})
}
```

model层下的user_basic

```go
// 创建user
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
```

Router/app.go

```go
r.GET("/user/createUser", service.CreateUser)
```

#### -3.删除用户

userService.go

```go
// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user) //推入数据库中
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}
```

model层下的user_basic

```go
// 删除用户
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}
```

Router/app.go

```go
r.GET("/user/deleteUser", service.DeleteUser)
```

#### -4.修改用户

userService.go

```go
// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")

	models.UpdateUser(user) //推入数据库中
	c.JSON(200, gin.H{
		"message": "修改用户成功",
	})
}

```

model层下的user_basic

```go
//修改用户

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord})
}
```

Router/app.go

```go
r.POST("/user/updateUser", service.UpdateUser)
```

Router/app下的所有Code

```go
package router

import (
	"IMsystemchat/docs"
	"IMsystemchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	return r

}

```

#### -5.对电话和email进行校验操作

加入修改电话号码和邮箱并校验

先引入

```go
go get github.com/asaskevich/govalidator
结构体字段后加入校验规则
Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
最后service 
如：
_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user) //推入数据库中
		c.JSON(200, gin.H{
			"message": "修改用户成功",
		})
	}
```

#### -6.重复注册的校验

```go
models/user_basic.go
// 通过名字查找对象
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =?", name).First(&user)
	return user
}

// 通过电话查找对象
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone =?", phone).First(&user)
}

// 通过email查找对象
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email =?", email).First(&user)
}


再去service/UserService.go修改Create的功能

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已注册！",
		})
		return
}
```

加密操作

```go
1.封装md5的工具类
md5.go
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 随机数加密操作
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + "         " + password)
	return md == password
}
写好工具类之后，在userService的creater用起来
// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
//在这里开始写入加密操作
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已注册！",
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致！",
		})
		return
	}
	//将密码给user对象
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt

	models.CreateUser(user) //推入数据库中
	c.JSON(200, gin.H{
		"message": "新增用户成功",
	})
}
```

#### -7.登录（解密）操作

```go
user_basic.go

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =? and pass_word=?", name, password).First(&user)
	return user
}


userService.go

// FindUserByNameAndPwd
// @Summary 登录用户
// @Tags 首页
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
    //判断用户是否存在
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "该用户不存在",
		})
		return
	}
	fmt.Println(user)
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码不正确",
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)

	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
```

#### -8.token的加入对返回的结构做了调整

models/user_Basic.go

```go
// 通过名字查找对象
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =?", name).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id =?", user.ID).Update("Identity", temp)
	return user
}
//返回的结果
c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "密码不正确",
			"data":    data,
		})
```

### 2.加入redis

导入redis

```go
1.载入相关的redis包
go get  github.com/garyburd/redigo/redis
go get github.com/go-redis/redis

2.在main()调用
utils.InitRedis()

3.在config/appyml下对redis进行配置
redis:
  addr: "127.0.0.1:6379"
  password: ""
  DB: 0
  poolSize: 30
  minIdleConn: 30
4.system_init.go完善InitRedis()方法
func InitRedis() {

	//数据库连接
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Red.Ping().Result()
	if err != nil {
		fmt.Println("init redis .....", err)
	} else {
		fmt.Println("Redis inited ......", pong)
	}

}
```

### 3.通过websocket连通

```go
1.下载包 go get github.com/gorilla/websocket
        go get github.com/go-redis/redis/v8
2.system_init.go
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

3.userService.go加入以下代码
// 防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 发送消息
func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHander(ws, c)
}
func MsgHander(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息：", msg)
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}


4.最后在router/app.go加入get请求
//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
```

测试：http://www.jsons.cn/websocket/

![](D:\myfile\GO\project\IMsystem\note\pic\websockettest.jpg)

查看是否测试成功

![](D:\myfile\GO\project\IMsystem\note\pic\websockettest2.jpg)

### 4.设计关系表、群信息表、消息表

message.go

```go
package models

import (
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FormId   uint   //发送者
	TargetId uint   //消息的接收者
	Type     string // 消息类型 群聊，私聊，广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他的数字统计

}

func (table *Message) TableName() string {
	return "message"
}

```

group_basic.go

```go
package models

import (
	"gorm.io/gorm"
)

// 群信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}

```

contact.go

```go
package models

import (
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //对应的类型 0 1
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

```

### 5.消息传递

发送消息

​	需要：发送者ID,接收者ID，消息类型，发送的内容，发送类型

​    校验Token，关系

接收消息

message.go

```go
package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 消息
type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //消息的接收者
	Type     int    // 发送类型 群聊，私聊，广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他的数字统计

}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface // go get gopkg.in/fatih/set.v0
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并校验token等合法性
	//token :=query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgType :=query.Get("type")
	//targetId :=query.Get("targetId")
	//context :=query.Get("context")
	isvalida := true //check token() 待。。。。。
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取连接Conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid 跟node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送的逻辑
	go sendProc(node)
	//6.完成接收的逻辑
	go recevProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recevProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<<", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		//case 2:
		//	sendGroupMsg() //群发
		//case 3:
		//	sendAllMsg() //广播
		//case 4:

	}
}
func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

```

### 6.集成html和前端页面

service/index

```go
package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	fmt.Println("进来了")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome !!",
	//})
}

// 跳转到注册页面进行注册操作
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	fmt.Println("进来了")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

```



router/app.go引入静态资源

```go
//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")
//首页相关的
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
```

service/userService.go

```go
// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//拿到数据
	user := models.UserBasic{}
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println("repassword: ", repassword)
	fmt.Println(user.Name, ">>>>>", password, repassword)
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	//看数据有没有在前端填写好，没有就返回错误信息
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "用户名或密码不能为空",
			"data":    user,
		})
		return
	}
	//当查询到了结果之后
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名已注册！",
			"data":    user,
		})
		return
	}
	//比对两次输入的密码是否正确
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	//将密码给user对象
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	fmt.Println(user.PassWord)
	models.CreateUser(user) //推入数据库中
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失败
		"message": "新增用户成功",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 登录用户
// @Tags 首页
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}
	fmt.Println(user)
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)

	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失败
		"message": "登录成功",
		"data":    data,
	})
}

```

head.html

```html
{{define "/chat/head.shtml"}}
<script>
    function userId(id){
        if(typeof  id =="undefined"){
            var r = sessionStorage.getItem("userid");
            if(!r){
                return 0;
            }else{
                return parseInt(r)
            }
        }else{
            sessionStorage.setItem("userid",id);
        }
    }
    function userInfo(o){
        if(typeof  o =="undefined"){
            var r = sessionStorage.getItem("userinfo");
            if(!!r){
                return JSON.parse(r);
            }else{
                return null
            }
        }else{
            sessionStorage.setItem("userinfo",JSON.stringify(o));
        }
    }
    var url = location.href;
    var isOpen = url.indexOf("/login")>-1 || url.indexOf("/register")>-1
    if (!userId() && !isOpen){
      // location.href = "login.shtml";
    }

</script>

    <!--聊天所需-->
<meta name="viewport" content="width=device-width, initial-scale=1,maximum-scale=1,user-scalable=no">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black">
<title>IM解决方案</title>
<meta name="Description" content="马士兵教育IM通信系统">
<meta name="Keywords" content="无人售货机，小程序，推送，群聊,单聊app">
<link rel="stylesheet" href="/asset/plugins/mui/css/mui.css" />
<link rel="stylesheet" href="/asset/css/chat.css" />
<link rel="stylesheet" href="/asset/css/audio.css" />
<!--登录所需 -->
<link rel="stylesheet" href="/asset/css/login.css" />
<link rel="stylesheet" href="/asset/iconfont/iconfont.css" />
<link rel="icon" href="asset/images/favicon.ico" type="image/x-icon"/>  
<script src="/asset/plugins/mui/js/mui.js" ></script>
<script src="/asset/js/vue.min.js" ></script>
<script src="/asset/js/vue-resource.min.js" ></script>
<script src="/asset/js/util.js" ></script>
<script>
   function post(uri,data,fn){
                var xhr = new XMLHttpRequest();
                xhr.open("POST","//"+location.host+"/"+uri, true);
                // 添加http头，发送信息至服务器时内容编码类型
                xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
                xhr.onreadystatechange = function() {
                    if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
                        fn.call(this, JSON.parse(xhr.responseText));
                    }
                };
                var _data=[];
                if(!! userId()){
                   // data["userid"] = userId();
                }
                for(var i in data){
                    _data.push( i +"=" + encodeURI(data[i]));
                }
                xhr.send(_data.join("&"));
            }
            function uploadfile(uri,dom,fn){
                var xhr = new XMLHttpRequest();
                xhr.open("POST","//"+location.host+"/"+uri, true);
                // 添加http头，发送信息至服务器时内容编码类型
                xhr.onreadystatechange = function() {
                    if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
                        fn.call(this, JSON.parse(xhr.responseText));
                    }
                };
                var _data=[];
                var formdata = new FormData();
                if(!! userId()){
                    formdata.append("userid",userId());
                }
                formdata.append("file",dom.files[0])
                xhr.send(formdata);
            }
   function uploadblob(uri,blob,filetype,fn){
       var xhr = new XMLHttpRequest();
       xhr.open("POST","//"+location.host+"/"+uri, true);
       // 添加http头，发送信息至服务器时内容编码类型
       xhr.onreadystatechange = function() {
           if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
               fn.call(this, JSON.parse(xhr.responseText));
           }
       };
       var _data=[];
       var formdata = new FormData();
       formdata.append("filetype",filetype);
       if(!! userId()){
           formdata.append("userid",userId());
       }
       formdata.append("file",blob)
       xhr.send(formdata);
   }
       function uploadaudio(uri,blob,fn){
                uploadblob(uri,blob,".mp3",fn)
        }
       function uploadvideo(uri,blob,fn){
           uploadblob(uri,blob,".mp4",fn)
       }
</script>

<style>
    .flex-container{
        display:flex;
        flex-direction:row;
        width:100%;
        padding-top: 10px;
        position: fixed;
        bottom: 0px;
        background-color: #FFFFFF;
    }
    .item-1{
        height:50px;
        height:50px;
        padding: 5px 5px 5px 5px;
    }
    .item-2{
                margin-right:auto;
        height:50px;
        width: 100%;
    }
    .txt{
        margin-right:auto;
    }
    .item-3{
        height:50px;
        height:50px;
        padding: 5px 5px 5px 5px;
    }
    .item-4{
        height:50px;
        height:50px;
        padding: 5px 5px 5px 5px;
    }

     li.chat{
         justify-content: flex-start;
         align-items: flex-start;
         display: flex;

     }
     .chat.other{
         flex-direction: row;
     }
    .chat.mine{
        flex-direction: row-reverse;
    }
    img.avatar{
        width: 54px;
        height: 54px;
    }
    .other .avatar{
        margin-left:10px;
    }
    .mine .avatar{
        margin-right:10px;
    }
    .other span{
        display: none;
        border: 10px solid;
        border-color: transparent #FFFFFF transparent transparent ;
        margin-top: 10px;
    }
    .mine span{
        display: none;
        border: 10px solid;
        border-color: transparent  transparent transparent #32CD32;
        margin-top: 10px;
    }
    .other>.content{
        background-color: #FFFFFF;

    }
    .mine>.content{
        background-color: #e3eafa;

    }
    div.content{
        min-width: 60px;
        clear: both;
        display: inline-block;
        padding: 16px 16px 16px 10px;
        margin: 0 0 20px 0;
        font: 16px/20px 'Noto Sans', sans-serif;
        border-radius: 10px;

        min-height: 54px;
    }
    .content>img.pic{
        width: 100%;
        margin:3px 3px 3px 3px;
    }
    .content>img.audio{
        width: 32px;
        color: white;
    }
    #panels{
        background-color: #FFFFFF;
        display: flex;
        position: fixed;
        bottom: 50px;
    }
    .doutures{
        flex-direction: row;
        flex-wrap: wrap;
        display: flex;
    }
    .doutures img{
        margin: 10px 10px 10px 10px;
    }
    .doutupkg{
        flex-direction: row;
        flex-wrap: wrap;
        display: flex;
    }
    .plugins{
        flex-direction: row;
        flex-wrap: wrap;
        display: flex;
    }
    .plugin{
        padding: 10px 10px 10px 20px;
        margin-left: 10px;
        margin-right: 10px;
    }
    .plugin img{
        width: 40px;
    }
    .plugin p{
        text-align: center;
        font-size: 16px;
    }
    .doutupkg img{
        width: 32px;
        height: 32px;
        margin: 5px 5px 5px 5px;
    }
    .upload{
        width: 64px;
        height: 64px;
        position: absolute;
        top: 1px;
        opacity:0;
    }
    .tagicon{
        width: 32px;
        height:32px;
    }
    
    .small{
        width: 32px;
        height:32px;
    }
    .middle{
        width: 64px;
        height:64px;
    }
    .large{
        width: 96px;
        height:96px;
    }
    .res image{
        width: 32px;
        height:32px;
    }
    .mui-content {
                padding-top: 44px;
                position: absolute;
                left: 0;
                top: 0;
                background: #fff;
                width: 100%;
                height: 100%;
        }
</style>
{{end}}
```



index.html

```html
<!DOCTYPE html>
<html>

<head>
    <!--js include-->
    {{template "/chat/head.shtml"}}
<!--    <title>够浪</title>-->
<!--    <link rel="stylesheet" href="/asset/plugins/mui/css/mui.css"/>-->
<!--    <link rel="stylesheet" href="/asset/css/login.css"/>-->
<!--    <script src="/asset/plugins/mui/js/mui.js"></script>-->
<!--    <script src="/asset/js/vue.min.js"></script>-->
<!--    <script src="/asset/js/util.js"></script>-->
</head>
<body>
<p>进入登录</p>
<header class="mui-bar mui-bar-nav">
    <h1 class="mui-title">登录</h1>
</header>
{{.}}
<div class="mui-content login-page" id="pageapp">
    <form id='login-form' class="mui-input-group login-from">
        <div class="mui-input-row">
            <input v-model="user.name" placeholder="请你输入用户名" type="text" class="mui-input-clear mui-input" >
        </div>
        <div class="mui-input-row">
            <input v-model="user.password" placeholder="请你输入密码"  type="password" class="mui-input-clear mui-input" >
        </div>
    </form>
    <div class="mui-content-padded">
        <button @click="login"  type="button"  class="mui-btn mui-btn-block mui-btn-primary btn-login">登录</button>
        <div class="link-area"><a id='reg' href="toRegister">注册账号</a> <span class="spliter">|</span> <a  id='forgetPassword'>忘记密码</a>
        </div>
    </div>
    <div class="mui-content-padded oauth-area">
    </div>
</div>
</body>
</html>
<script>
    var app = new Vue({
        el:"#pageapp",
        data:function(){
          return {
              user:{
                name:"",
                password:"",
              }
          }
        },
        methods:{
            login:function(){
                //检测手机号是否正确
                console.log("login")
                //检测密码是否为空

                //网络请求
                //封装了promis
                util.post("user/findUserByNameAndPwd",this.user).then(res=>{
                    console.log(res)
                    if(res.code!=0){
                        mui.toast(res.message)
                    }else{         
                        var url = "/toChat?userId="+res.data.ID+"&token="+res.data.Identity
                        userInfo(res.data)
                        userId(res.data.ID)
                        mui.toast("登录成功,即将跳转")
                        location.href = url
                    }
                })
            },
        }
    })
</script>
```

register.html

```html
<!DOCTYPE html>
<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>IM解决方案</title>
    <link rel="stylesheet" href="/asset/plugins/mui/css/mui.css" />
    <link rel="stylesheet" href="/asset/css/login.css" />
    <link rel="icon" href="asset/images/favicon.ico" type="image/x-icon" />
    <script src="/asset/plugins/mui/js/mui.js"></script>
    <script src="/asset/js/vue.min.js"></script>
    <script src="/asset/js/util.js"></script>
</head>

<body>

    <header class="mui-bar mui-bar-nav">
        <h1 class="mui-title">注册</h1>
    </header>
    <div class="mui-content register-page" id="pageapp">
        <form id='login-form' class="mui-input-group register-form">
            <div class="mui-input-row">
                <input v-model="user.name" placeholder="请输入用户名" type="text" class="mui-input-clear mui-input">
            </div>
            <div class="mui-input-row">
                <input v-model="user.password" placeholder="请输入密码" type="password" class="mui-input-clear mui-input">
            </div>
            <div class="mui-input-row">
                <input v-model="user.Identity" placeholder="再输入密码" type="password" class="mui-input-clear mui-input">
            </div>
        </form>
        <div class="mui-content-padded">
            <button @click="login" type="button" class="mui-btn mui-btn-block mui-btn-primary btn-register">注册</button>
            <div class="link-area"><a id='reg' href="/index">登录账号</a> <span class="spliter">|</span> <a
                    id='forgetPassword'>忘记密码</a>
            </div>
        </div>
        <div class="mui-content-padded oauth-area">
        </div>
    </div>
</body>

</html>
<script>
    var app = new Vue({
        el: "#pageapp",
        data: function () {
            return {
                user: {
                    name: "",
                    password: "",
                    Identity: "",
                }
            }
        },
        methods: {
            login: function () {
                //检测密码是否为空
                console.log(this.user)
                //网络请求
                //封装了promis
                util.post("/user/createUser", this.user).then(res => {
                    console.log(res)
                    if (res.code != 0) {
                        mui.toast(res.message)
                    } else {
                        location.replace("//127.0.0.1:8081/index")
                       // location.href = "/"
                        mui.toast("注册成功,即将跳转")
                    }
                })
            },
        }
    })
</script>
```







看完43集接下来看44集了
