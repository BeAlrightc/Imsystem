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

### 5.完成用户模块基本的功能

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

#### -5.完成用户模块基本的功能

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

学完36集
