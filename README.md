[![Badge](https://img.shields.io/badge/link-Tigo-blue.svg)](https://karldoenitz.github.io/Tigo/)
[![LICENSE](https://img.shields.io/badge/license-tission-blue.svg)](https://github.com/karldoenitz/tission/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/karldoenitz/tission.svg?branch=master)](https://travis-ci.org/karldoenitz/tission)
# tission

一个基于redis实现的session处理插件，可内嵌入`Tigo`框架使用。

## 安装

```shell
go get github.com/karldoenitz/tission
```

## 使用

将此包引入到`Tigo`项目中，简单配置之后即可使用，Demo如下所示：

```go
package main

import (
	"fmt"
	"github.com/karldoenitz/Tigo/TigoWeb"
	"github.com/karldoenitz/tission/session/redis"
)

type Test struct {
	A string
	B string
}

type FrontHandler struct {
	TigoWeb.BaseHandler
}

func (frontHandler *FrontHandler) Get() {
	param := frontHandler.GetParameter("a").ToString()
	t := Test{}
	if param == "a" {
		frontHandler.GetSession("abc", &t)    // 从session中取值
		frontHandler.ResponseAsJson(t)
		return
	}
	t.A = param
	t.B = "from session"
	e := frontHandler.SetSession("abc", t)  // 设置session值
	if e != nil {
		fmt.Println(e.Error())
	}
	frontHandler.ResponseAsJson(t)
}

type Testa struct {
	TigoWeb.BaseHandler
}

func (t * Testa) Get() {
	param := t.GetParameter("a").ToString()
	if param == "a" {
		tt := Test{}
		if e := t.GetSession("abc", &tt); e != nil {
			println(e.Error())
		}
		t.ResponseAsJson(tt)
		return
	}
}

var urlMapping = []TigoWeb.Router{
	{"/front", &FrontHandler{}, nil},
	{"/tt", &Testa{}, nil},
}

func main() {
	application := TigoWeb.Application{IPAddress: "0.0.0.0", Port: 8888, UrlRouters: urlMapping}
	t := redis.SessionInterface{
		IP:      "127.0.0.1",  // redis地址
		Port:    "6379",       // redis端口
		MaxIdle: 10,
		Timeout: 100,
		Expire:  3600,         // 3600s后session过期
	}
	application.StartSession(&t, "tid")
	application.Run()
}
```

