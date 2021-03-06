package handler

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"wifi_auth/db"
	"wifi_auth/jwt"
	"wifi_auth/template"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var AdminKey string = "c16cbe118a80436b5b6fe3eb15ffc37d"

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

	clientIP := c.DefaultQuery("ip", "")
	if clientIP == "" {
		log.Println("no ip")
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	gwSSLOn := c.DefaultQuery("gw_ssl_on", "no")
	originUrl := c.DefaultQuery("url", "")

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
	token, err := jwt.GenerateToken(user.Id, user.Username, clientIP, clientMac, user.Level)
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

	log.Println(username, "login success")

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
	clientMac := c.DefaultQuery("mac", "")
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
		log.Println("clientIP: ", clientIP)
		log.Println("tokenIP: ", claims.IP)
		log.Println("clientMac: ", clientMac)
		log.Println("tokenMac: ", claims.Mac)
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	ol := db.OnlineUser{
		Username: claims.Username,
		Level:    claims.Level,
		IP:       clientIP,
		Mac:      clientMac,
	}

	/*
		set expiration time
	*/
	if claims.Level == 1 {
		ol.ExpiredAt = time.Now().Add(time.Duration(49*24) * time.Hour)
	} else {
		ol.ExpiredAt = time.Now().Add(time.Duration(2) * time.Hour)
	}
	ol.ExpiredTimeStamp = ol.ExpiredAt.Unix()

	err = db.AddUser2List(&ol)
	if err != nil {
		log.Println("auth failed")
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	log.Println("[log]Message: Auth Success")
	ret := fmt.Sprintf(`Auth: %d`, claims.Level)
	log.Println("[log]Message: ", ret)
	c.String(http.StatusOK, ret)
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

func AddUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key != AdminKey {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}
	userInfo := template.User{}
	err := c.BindJSON(&userInfo)
	if err != nil {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}
	log.Println("go here")
	log.Println(userInfo)

	cipherPwd, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	newUser := db.User{}
	newUser.Username = userInfo.Username
	newUser.Password = string(cipherPwd)
	newUser.Level = userInfo.Level

	err = db.CreateUser(&newUser)
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}
	c.String(http.StatusCreated, "user added successfully")
	return
}

func GetOnlineUserList(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key != AdminKey {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	users, err := db.GetOnlineUserList()
	if err != nil {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	c.JSON(http.StatusOK, users)
	return
}

func KickOutUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key != AdminKey {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	mac := c.DefaultQuery("mac", "")
	if mac == "" {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	username := c.DefaultQuery("username", "")
	if username == "" {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	err := db.KickOutUser(username, mac)
	if err != nil {
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	c.String(http.StatusOK, "kick out user successfully")
	return
}
