package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	mac := c.DefaultQuery("mac", "")
	if mac == "" {
		log.Println("no mac")
		c.Redirect(http.StatusFound, "/msg?msg=Please contact the network administrator")
		return
	}

	gwSSLOn := c.DefaultQuery("gw_ssl_on", "no")
	originUrl := c.DefaultQuery("url", "")
	ip := c.ClientIP()

	log.Println(gwAddress, gwPort, mac)
	username := c.PostForm("username")
	password := c.PostForm("password")

	log.Println(username, password)

	if username == "root" && password == "123456" {
		token := "12345678"
		var url string
		if gwSSLOn == "yes" {
			url = fmt.Sprintf("https://%s:%s/auth?token=%s&stage=login&mac=%s&ip=%s&url=%s", gwAddress, gwPort, token, mac, ip, originUrl)
		} else {
			url = fmt.Sprintf("http://%s:%s/auth?token=%s&stage=login&mac=%s&ip=%s&url=%s", gwAddress, gwPort, token, mac, ip, originUrl)
		}

		c.Redirect(http.StatusFound, url)
		return
	} else {
		log.Println("password incorrect")
		c.Redirect(http.StatusFound, "/login?"+c.Request.URL.RawQuery)
		return
	}
	return
}

func Portal(c *gin.Context) {
	log.Println(c.Request.URL.RawQuery)
	//c.Redirect(http.StatusFound, "https://feiyu.com")
	c.HTML(http.StatusOK, "portal.html", nil)
	return
}

func Auth(c *gin.Context) {
	stage := c.DefaultQuery("stage", "")
	if stage == "" {
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	token := c.DefaultQuery("token", "")
	if token == "" {
		log.Println("token未找到")
		c.String(http.StatusOK, "Auth: 0")
		return
	}

	if stage == "login" {
		if token == "12345678" {
			c.String(http.StatusOK, "Auth: 1")
		} else {
			log.Println("token is incorrect")
			c.String(http.StatusOK, "Auth: 0")
		}
		return
	}

	c.String(http.StatusOK, "Auth: 1")
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
