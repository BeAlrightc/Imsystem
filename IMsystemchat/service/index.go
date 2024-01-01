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
