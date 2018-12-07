package main

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"github.com/lllwan/redirect/model"
	"fmt"
)

var conf = model.Config

func Login(g *gin.Context){
	username := g.DefaultPostForm("username", "")
	password := g.DefaultPostForm("password", "")
	if username == "" || password == ""{
		g.JSON(403, gin.H{
			"errcode": 403,
			"errmsg": "用户名或密码不可为空！",
		})
		g.Abort()
		return
	}
	var user model.Users
	if err := model.DB.Where("username=? and password=?",
		username, model.Convert(username, password)).First(&user).Error;err != nil{
		g.JSON(403, gin.H{
			"errcode": 403,
			"errmsg": "登入失败！",
		})
		g.Abort()
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user"] = username
	claims["ip"] = g.Request.RemoteAddr
	token.Claims = claims
	tokenStr, err := token.SignedString([]byte(conf.SECRET))
	if err != nil{
		g.JSON(500, gin.H{
			"errcode":500,
			"errmsg":"内部错误：Token信息生成失败！",
		})
		log.Println(username, "认证成功，但是生成token失败！", err)
		g.Abort()
		return
	}
	g.JSON(200, gin.H{
		"errcode":0,
		"errmsg":"OK",
		"token":tokenStr,
	})
}
func Validate() gin.HandlerFunc{
	return func(g *gin.Context) {
		token := g.GetHeader("authorization")
		if token == ""{
			g.JSON(403, gin.H{
				"errcode":403,
				"errmsg":"请求非法！",
			})
			g.Abort()
			return
		}
		auth,err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.SECRET), nil
		})
		if !auth.Valid{
			g.JSON(403, gin.H{
				"errcode":403,
				"errmsg":"Not authorized！",
				"result":err.Error(),
			})
			g.Abort()
		}
		g.Set("claims", auth.Claims)
	}
}
func UrlForward(g *gin.Context){
	Host := g.Request.Host
	var acl model.Acls
	if err := model.DB.Where("domain=?",Host).First(&acl).Error;err != nil{
		log.Println("数据库查询失败，Host：", g.Request.Host, err)
		g.JSON(403, gin.H{
			"errcode":403,
			"errmsg":"数据查询失败！",
		})
		g.Abort()
		return
	}
	switch acl.Method {
	default:
		g.HTML(200, "template.html", gin.H{
			"title": acl.Title,
			"count": acl.Count,
			"countid": acl.Countid,
			"keywords": acl.Keywords,
			"description": acl.Description,
			"favicon": acl.Favicon,
			"url": acl.Url,
			"path":g.Request.URL.Path,
		})
	case "301":
		g.Redirect(301, acl.Url)
	case "301all":
		g.Redirect(301, acl.Url + g.Request.URL.Path)
	}
}
func ChangePassword(g *gin.Context){
	claims := g.MustGet("claims")
	oldpassword := g.DefaultPostForm("oldpass", "")
	password := g.DefaultPostForm("password", "")
	if len([]rune(password)) >=8 {
		var user model.Users
		claiminfo := claims.(jwt.MapClaims)
		username := claiminfo["user"].(string)
		if err := model.DB.Where("username=? and password=?", username, model.Convert(username, oldpassword)).First(&user).Error;err != nil{
			g.JSON(403, gin.H{
				"errcode": 403,
				"errmsg": "旧密码验证失败！",
			})
			g.Abort()
			return
		}
		New := model.Convert(claiminfo["user"].(string), password)
		err := model.DB.Model(&user).Where("username = ?", claiminfo["user"].(string)).Update("password", New).Error
		if err != nil{
			g.JSON(500, gin.H{
				"errcode": 500,
				"errmsg": fmt.Sprintf("密码修改失败，因为:%s", err.Error),
			})
			g.Abort()
			return
		}
		g.JSON(200, gin.H{
			"errcode": 200,
			"errmsg": "OK",
		})
	} else {
		g.JSON(403, gin.H{
			"errcode": 403,
			"errmsg": "密码位数必须大于8！",
		})
	}
}

func RemoveAcl(g *gin.Context){
	domain := g.PostForm("domain")
	if domain == ""{
		g.JSON(403, gin.H{
			"errcode": 403,
			"errmsg": "缺少必需参数：domain",
		})
		g.Abort()
		return
	}
	err := model.DB.Delete(model.Users{}, "domain = ？", domain)
	if err != nil{
		g.JSON(500, gin.H{
			"errcode": 500,
			"errmsg": fmt.Sprintf("删除失败:%s", err.Error),
		})
		g.Abort()
		return
	}
	g.JSON(200, gin.H{
		"errcode": 200,
		"errmsg": "OK",
	})
}

func CreateAcl(g *gin.Context){
	domain := g.DefaultPostForm("domain", "")
	url := g.DefaultPostForm("url", "")
	keywords := g.DefaultPostForm("keywords", "")
	description := g.DefaultPostForm("description", "")
	favicon := g.DefaultPostForm("favicon", "")
	method := g.DefaultPostForm("method", "hide")
	comment := g.DefaultPostForm("comment", "")
	count := g.DefaultPostForm("count", "")
	countid := g.DefaultPostForm("countid", "")
	title := g.DefaultPostForm("title", "")
	claims := g.MustGet("claims")
	if  domain == "" || url == ""{
		g.JSON(403, gin.H{
			"errcode": 403,
			"errmsg": "域名或目标URL为空！",
		})
		g.Abort()
		return
	}
	Acl := model.Acls{
		Domain:domain,
		Url:url,
		Keywords:keywords,
		Description:description,
		Favicon:favicon,
		Method:method,
		Comment:comment,
		Count:count,
		Countid:countid,
		Title:title,
		Username:claims.(jwt.MapClaims)["user"].(string),
	}
	if err := model.DB.Create(Acl).Error;err != nil{
		g.JSON(500, gin.H{
			"errcode": 500,
			"errmsg": "内部错误, 数据添加失败！",
		})
		g.Abort()
		return
	}
	g.JSON(200, gin.H{
		"errcode": 200,
		"errmsg": "OK",
	})
}


func main(){
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.POST("/login", Login)
	api := r.Group("/api", Validate())
	api.POST("/ChangePassword", ChangePassword)
	api.POST("/CreateAcl", CreateAcl)
	api.POST("/RemoveAcl", CreateAcl)
	r.GET("", UrlForward)
	r.Run(conf.HTTP_BIND)
}
