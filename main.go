package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gin-swagger/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type (
	User struct {
		Id   int    `json:"id" example:"1"`      // ID
		Name string `json:"name" example:"jack"` // 用户名
	}
	FailResult struct {
		Code int    `json:"code" example:"-1"`  // 状态码
		Msg  string `json:"msg" example:"fail"` // 消息
	}
	SuccessResult struct {
		Code int         `json:"code"` // 状态码
		Msg  string      `json:"msg"`  // 消息
		Data interface{} `json:"data"` // 数据
	}
	LoginRequest struct {
		Name     string `json:"name" example:"admin"`      // 用户名
		Password string `json:"password" example:"123456"` // 密码
	}
)

func setFailResult(c *gin.Context, msg interface{}) {
	res := FailResult{
		Code: -1,
		Msg:  fmt.Sprint(msg),
	}
	c.JSON(http.StatusOK, res)
}

func setSuccessResult(c *gin.Context, data interface{}) {
	res := SuccessResult{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Tags 用户
// @Produce json
// @Param 请求参数 body LoginRequest true "请求参数"
// @Router /api/v1/user/login [post]
// @Success 200 {object} SuccessResult
// @Failure 400 {object} FailResult
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		setFailResult(c, "参数错误")
		return
	}
	if req.Name != Username || req.Password != Password {
		setFailResult(c, "用户名或密码错误")
		return
	}
	token := MD5([]byte(req.Name + req.Password + TokenSalt))
	data := map[string]interface{}{"token": token}
	setSuccessResult(c, data)
}

// @Summary 根据ID获取用户信息
// @Description 根据ID获取用户信息
// @Accept json
// @Tags 用户
// @Produce json
// @Param id path int true "id"
// @Security ApiKeyAuth
// @Router /api/v1/user/queryById/{id} [get]
// @Success 200 {object} SuccessResult{data=User}
// @Failure 400 {object} FailResult
func queryById(c *gin.Context) {
	sid := c.Param("id")
	if sid == "" {
		setFailResult(c, "id不能为空")
	} else {
		id, err := strconv.Atoi(sid)
		if err != nil {
			setFailResult(c, "参数不正确")
		}
		data := User{
			Id:   id,
			Name: "jack",
		}
		setSuccessResult(c, data)
	}
}

// @Summary 根据用户名获取用户信息
// @Description 根据用户名获取用户信息
// @Accept json
// @Tags 用户
// @Produce json
// @Param name query string true "用户名"
// @Security ApiKeyAuth
// @Router /api/v1/user/queryByName [get]
// @Success 200 {object} SuccessResult{data=User}
// @Failure 400 {object} FailResult
func queryByName(c *gin.Context) {
	name := c.Query("name")
	if name != "" {
		data := User{
			Id:   1,
			Name: "jack",
		}
		setSuccessResult(c, data)
		return
	}
	setFailResult(c, "用户名不能为空")
}

// @Summary 添加用户信息
// @Description 添加用户信息
// @Accept json
// @Tags 用户
// @Produce json
// @Param 请求参数 body User true "请求参数"
// @Security ApiKeyAuth
// @Router /api/v1/user/addUser [post]
// @Success 200 {object} SuccessResult{data=User}
// @Failure 400 {object} FailResult
func addUser(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		setFailResult(c, "参数错误")
		return
	}
	data := User{
		Id:   1,
		Name: user.Name,
	}
	setSuccessResult(c, data)
}

// @Summary Excel文件上传
// @Description Excel文件上传
// @Accept mpfd
// @Tags Excel文件上传
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "Excel文件"
// @Router /api/v1/upload [post]
// @Success 200 {object} SuccessResult
// @Failure 400 {object} FailResult
func uploadFile(c *gin.Context) {
	setSuccessResult(c, nil)
}

const (
	// 可自定义盐值
	TokenSalt = "default_salt"
	// 默认用户名密码
	Username = "admin"
	Password = "123456"
)

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authentication") // 访问令牌
		if strings.ToLower(MD5([]byte(Username+Password+TokenSalt))) == strings.ToLower(token) {
			// 验证通过，会继续访问下一个中间件
			c.Next()
		} else {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
		}
	}
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authentication
func main() {
	docs.SwaggerInfo.Title = "Restful API接口文档"
	docs.SwaggerInfo.Description = "这是后端接口文档"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", "127.0.0.1", "8888")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()

	r.POST("/api/v1/user/login", login)

	// 接口文档
	url := ginSwagger.URL("http://127.0.0.1:8888/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Use(Authorize())

	user := r.Group("/api/v1/user")
	{
		user.GET("/queryById/:id", queryById)
		user.GET("/queryByName", queryByName)
		user.POST("/addUser", addUser)
	}

	r.POST("/api/v1/upload", uploadFile)

	if err := r.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}
