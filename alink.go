package alink

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

//create new io.reader body
func NewRespBody(respBody io.Reader) (*bytes.Reader, error) {

	b, err := ioutil.ReadAll(respBody)
	reader := bytes.NewReader(b)
	return reader, err
}

//检查是否是url
//check string is url
func IsValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	//log.Print( u.Scheme,"--", u.Host)
	return true
}

//check title
func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

//check video
func isVideoElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "video"
}

//check a
func isAHrefElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

//check img
func isImgElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

//get page title
func titleText(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := titleText(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

func videoSrc(node *html.Node) (string, bool) {
	if isVideoElement(node) {
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				return attr.Val, true
			}
		}
		return "", true
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		mark, ok := videoSrc(c)
		if ok {
			return mark, ok
		}
	}
	return "", false
}

//get video src
func Video(httpBody *bytes.Reader) (s [] string, err error) {
	var src []string
	node, err := html.Parse(httpBody)
	if err != nil {
		return src, err
	}
	link, flag := videoSrc(node)
	if flag {
		src = append(src, link)
	}
	return src, nil
}

//get pages title
func Title(httpBody *bytes.Reader) (t string, err error) {
	title := ""
	//httpBodyf = ioutil.NopCloser(httpBody)
	node, err := html.Parse(httpBody)
	if err != nil {
		return title,err
	}
	title, _ = titleText(node)
	return title, nil
}

//get all links
func Alink(httpBody *bytes.Reader) (l *[]string, b bool) {
	var links []string
	node, err := html.Parse(httpBody)
	if err != nil {
		return &links,false
	}
	ff, _ := alLink(node, &links)
	return ff, true
}

//get href url
func alLink(node *html.Node, h *[]string) ( f *[]string,n bool) {
	b := false

	if isAHrefElement(node) {
		for _, a := range node.Attr {
			if a.Key == "href" {
				s := trimHash(a.Val)
				if check(h,s)==false{
					//*h = make([]string,0,len(*h)+len(s))
					//log.Println(cap(f))
					*h= append(*h, s)
				}

				return h, true
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		all, flag := alLink(c, h)
		h = all
		b = flag
	}
	return h, b
}


func trimHash(l string) string {
	if strings.Contains(l, "#") {
		var index int
		for n, str := range l {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		return l[:index]
	}
	return l
}


//check url exits
func check(sl *[]string, s string) bool {
	var check bool
	for _, str := range *sl {
		if str == s {
			check = true
			break
		}
	}
	return check
}