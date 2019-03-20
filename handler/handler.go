package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/Nrehearsal/wifi_auth/db"
	"github.com/Nrehearsal/wifi_auth/jwt"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
	return
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
	return
}

func LoginCheck(c *gin.Context) {
	gwAddress := c.DefaultQuery("gw_address", "")
	if gwAddress == "" {
		log.Println("no address")
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	gwPort := c.DefaultQuery("gw_port", "")
	if gwPort == "" {
		log.Println("no port")
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	clientMac := c.DefaultQuery("mac", "")
	if clientMac == "" {
		log.Println("no mac")
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	gwSSLOn := c.DefaultQuery("gw_ssl_on", "no")
	originUrl := c.DefaultQuery("url", "")
	clientIP := c.ClientIP()

	log.Println(gwAddress, gwPort, clientMac)

	username := c.PostForm("username")
	password := c.PostForm("password")

	log.Println(username, password)

	user, err := db.GetUserByName(username)
	if err != nil {
		log.Println("no such a user")
		c.Redirect(http.StatusFound, "/login?"+c.Request.URL.RawQuery)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("password is incorrect")
		c.Redirect(http.StatusFound, "/login?"+c.Request.URL.RawQuery)
		return
	}

	//Generate token
	token, err := jwt.GenerateToken(user.Id, user.Usernmae, clientIP, clientMac, user.AccessDuration)
	if err != nil {
		log.Println("username or password is incorrect")
		c.Redirect(http.StatusFound, "/login?"+c.Request.URL.RawQuery)
		return
	}

	//all things in token
	var url string
	if gwSSLOn == "yes" {
		url = fmt.Sprintf("https://%s:%s/auth?token=%s&stage=login&mac=%s&ip=%s&url=%s", gwAddress, gwPort, token, clientMac, clientIP, originUrl)
	} else {
		url = fmt.Sprintf("http://%s:%s/auth?token=%s&stage=login&mac=%s&ip=%s&url=%s", gwAddress, gwPort, token, clientMac, clientIP, originUrl)
	}

	c.Redirect(http.StatusFound, url)
	return
}

func Portal(c *gin.Context) {
	log.Println(c.Request.URL.RawQuery)

	//c.Redirect(http.StatusFound, "https://www.cuit.edu.cn")
	originUrl := c.DefaultQuery("url", "")
	if originUrl == "" {
		c.HTML(http.StatusOK, "portal.html", nil)
		return
	}

	decodeString, err := base64.StdEncoding.DecodeString(originUrl)
	if err != nil {
		c.HTML(http.StatusOK, "portal.html", nil)
		return
	}

	c.Redirect(http.StatusFound, string(decodeString))
	return
}

func Auth(c *gin.Context) {
	clientIP := c.DefaultQuery("ip", "")
	clientMac := c.DefaultQuery("ip", "")
	stage := c.DefaultQuery("stage", "")
	token := c.DefaultQuery("token", "")

	if stage != "login" && stage != "logout" {
		log.Println("Unknown stage")
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Println("token is incorrect")
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	if claims.IP != clientIP || claims.Mac != clientMac {
		log.Println("auth failed")
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	c.String(http.StatusOK, string(claims.AccessDuration))
	return
}

func Msg(c *gin.Context) {
	msg := c.DefaultQuery("msg", "")
	if msg == "" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	c.String(http.StatusOK, msg)
	return
}

func AddUserAccount(c *gin.Context) {

}

func GetValidMacList(c *gin.Context) {

}
