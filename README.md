# alink
- [中文](#中文)

Golang package to read href,video,title ...  tags from an HTML page。



## 中文
一个Golang package 用来读取HTML页面中的 <title> ，<video>，<a> 等元素
输入一个   http.Get 返回的 response 使用 html.Parse 解析后返回一个字符串数组指针



### 例子
```go

package main
import (
	"https://github.com/gitgitcode/alink"
	"golang.org/x/net/html"
	"fmt" 
)

func main(){
    resp,_ := http.Get("http://www.testtest.com")
    newResp ,err := alink.NewRespBody(resp.Body)
    	if err !=nil{
    		log.Print(err.Error())
    	}
    links,_ := alink.Alink(newResp)
    fmt.Println(links)
    
}
```
