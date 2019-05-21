# wifi_auth

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

wifidog认证协议golang实现，使用详情参考wifidog认证协议标准

  - 在线示例，http://wifi.return0.top
  - 需要配合wifidog客户端使用，详情见[go_captive_portal](https://github.com/Nrehearsal/go_captive_portal)
  - 与标准实现不同，wifi_auth对用户原始url进行了url编码，还增加了用户状态持久化等功能
  - 网关id，客户都ip，mac等参数不能为空，否则会show error message。

### 安装使用说明
```sh
    $ go build
    $ mv wifi_auth resource/
    $ cd resource/
    $ tar -zxvf static.tar.gz
    $ ./wifi_auth
```
### API接口简单说明
```go
    //登录页面
    router.GET("/login", handler.Login) 
    
    //认证成功显示的门户页面
    router.GET("/portal", handler.Portal)
    
    //心跳API
    router.GET("/ping", handler.Ping)
    
    //登录页提交用户名密码的action，检查用户名，生成token等
    router.POST("/logincheck", handler.LoginCheck)
    
    //token校验API，校验参数?token=xxxxxx是否合法。
    router.GET("/auth", handler.Auth)
    
    //自定义错误页，例如/msg?msg=not found
    router.GET("/msg", handler.Msg)
    
    //添加用户的API，将新用户信息写入SQLite
    router.POST("/adduser", handler.AddUser)
    
    //获取在线用户的列表
    router.GET("/onlinelist", handler.GetOnlineUserList)
    
    //强制踢出用户
    router.POST("/kickout", handler.KickOutUser)
```
