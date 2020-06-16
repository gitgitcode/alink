# alink

[![Build Status](https://travis-ci.org/gitgitcode/alink.svg?branch=master)](https://travis-ci.org/gitgitcode/alink)
[![Coverage Status](https://coveralls.io/repos/github/gitgitcode/alink/badge.svg?branch=master)](https://coveralls.io/github/gitgitcode/alink?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gitgitcode/alink)](https://goreportcard.com/report/github.com/gitgitcode/alink)
- [ZH](#简介)

Golang package to read href,video,title,img ...  tags from an HTML page。


## 简介

一个简单的Golang package 主要用来读取HTML页面中的 ``` <title> ，<video>的src，<a>的href，<img>的src``` 等元素的内容.
在库里提供了两种方式处理 ```http.Get``` 返回的```response.Body```内容，一是通过  ```alink.GetBytesReaderWithIoReader```方法处理可以读取 ```http.Get``` 返回的```response.Body```内容。
但是如果要***多次***读取使用io.Reader 要通过 ```body, err := ioutil.ReadAll(b.Body)```读取后再次新建 ``` readerHref := bytes.NewReader(body)``` 的方式来进行。
第二中就是使用 ```alink.GetByteWithIoReader``` 方法读取```http.Get``` 返回的```response.Body``` 使用``WithByte``后缀的方进行多次读取.
内部方法使用html.Parse 解析后内容。


### 例子 Example

- 一个读取google/baidu主页的例子。获取页面的img和全部a连接并打印出来

- Use http client Get google/baidu Index Page and collect tags img ,href     

```go
package main

import (
	    "github.com/gitgitcode/alink"
	    "bytes"
       	"fmt"
       	"io/ioutil"
       	"log"
       	"math/rand"
       	"net/http"
       	"time"
    )
    var userAgentList = []string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
    	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
    	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
    	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
    	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
    	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
    	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
    	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
    	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
    	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
    	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
    	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
    	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
    	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}
    
    func GetRandomUserAgent() string{
    	r := rand.New(rand.NewSource(time.Now().UnixNano()))
    	return userAgentList[r.Intn(len(userAgentList))]
    }
    
    var accept = "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
    
    func ReqAdd(req *http.Request) {
    	req.Header.Set("Cookie","sug=3; a=1; ORIGIN=0; bdime=21110")
    	req.Header.Add("User-Agent",GetRandomUserAgent() )
    	req.Header.Add("Accept",accept)
    	req.Header.Add("Upgrade-Insecure-Requests","1")
    }
    
    func main() {
    
    	Response ,_:= GetHttpResponseP()
    	body, err := ioutil.ReadAll(Response.Body)
    	if err !=nil{
    		panic(err)
    	}
    	GetWithByte(body)
    	GetWithBytesReaderCreateTwiceNewReader(body)
    
    }
    
    func GetHttpResponseP() (*http.Response,error){
    	str:="https://google.co.jp"
    	str1:="https://www.baidu.com"
    
    	//fmt.Print(alink.IsValidUrl(str1))
    	client:= http.Client{Timeout: 2 * time.Second}
    	req,err := http.NewRequest("GET",str,nil)
    	req1,err1 := http.NewRequest("GET",str1,nil)
    
    	if err != nil{
    		log.Printf("google is err:%s",err.Error())
    	}
    
    	if err1 != nil{
    		log.Printf("baidu is err:%s",err1.Error())
    	}
    
    	ReqAdd(req)
    	ReqAdd(req1)
    	b,err := client.Do(req)
    	defer client.CloseIdleConnections()
    	if err != nil{
    		log.Printf("request google err %s",err.Error())
    		b1,err1 := client.Do(req1)
    		if err1 !=nil{
    			log.Printf("request baidu err %s",err.Error())
    			panic(err1)
    		}
    		b = b1
    	}
    	return b ,nil
    }
    func GetWithByte(body []byte)  {
    
    	 title, err:= alink.GetTitleWithByte(body)
    	 if err == nil{
    	 	fmt.Println(title)
    	 }else{
    	 	fmt.Println("GetWithByte GetTitleWithByte err")
    	 }
    	src,err := alink.GetImgSrcWithByte(body)
    	if err == nil{
    		for _,s :=range *src{
    			fmt.Println(s)
    		}
    	}else{
    		fmt.Println("GetWithByte GetImgSrcWithByte err")
    	}
    
    }
    
    func GetWithBytesReaderCreateTwiceNewReader(body []byte){
    	fmt.Println("<=================>")
    	//for read twice create new reader
    	readerHref := bytes.NewReader(body)
    	//创建两个新 reader
    	readerImg := bytes.NewReader(body)
    
    	t,f := alink.GetHrefWithBytesReader (readerImg)
    
    	if f !=nil {
    		log.Print(f)
    	}
    	fmt.Printf("Href:%s \n",t)
    
    	a,bl := alink.GetImgSrcWithBytesReader(readerHref)
    
    	if bl ==nil{
    		for i,v := range *a{
    			fmt.Printf("index:%d=href:%s\n",i,v)
    		}
    	}
    
    	//title:百度一下,你就知道
    	//index:0=href:/
    	// index:1=href:javascript:;
    	// index:2=href:https://passport.baidu.com/v2
    	//or
    	//title:Google
    	//index:0=href:/
    	// index:1=href:javascript:;
    	// index:2=href:https://wwww.google.com/
    }
```
