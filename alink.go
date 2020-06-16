// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package alink  a html tag collect package
//
//
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

// GetBytesReaderWithIoReader create a new bytes reader
func GetBytesReaderWithIoReader(respBody io.Reader)(reader *bytes.Reader ,err error){

	c, err := ioutil.ReadAll(respBody)
	if err == nil{
		reader = bytes.NewReader(c)
	}
	return reader, err
}

// GetByteWithIoReader is the func use ioutil.ReadAll() change to byte
// It returns []byte
func GetByteWithIoReader(respBody io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(respBody)
	return b, err
}


// GetByteReader use  bytes.NewReader create a new reapBody to read
func GetByteReader(respBody []byte) *bytes.Reader {
	reader := bytes.NewReader(respBody)
	return reader
}

// 检查是否是url
// check string is url
// It return bool
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

// isTitleElement
func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

// Check video
func isVideoElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "video"
}

// isAHrefElement
func isAHrefElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

// isImgElement
func isImgElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "img"
}

// Get page title
func getTitleText(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := getTitleText(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

// videoSrc get video src
func getVideoSrc(node *html.Node) (string, bool) {
	if isVideoElement(node) {
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				return attr.Val, true
			}
		}
		return "", true
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		mark, ok := getVideoSrc(c)
		if ok {
			return mark, ok
		}
	}
	return "", false
}

// VideoSrc get the video tags src
// It returns []string
func GetVideoSrcWithBytesReader(httpBody *bytes.Reader) (s []string, err error) {
	var src []string
	node, err := html.Parse(httpBody)
	if err != nil {
		return src, err
	}
	link, flag := getVideoSrc(node)
	if flag {
		src = append(src, link)
	}
	return src, nil
}

// TitleBytes to get pages title return a string
func TitleBytes(httpBody *bytes.Reader) (t string, err error) {
	title := ""
	node, err := html.Parse(httpBody)
	if err != nil {
		return title, err
	}

	title, _ = getTitleText(node)

	return title, nil
}

// GetTitleWithByte
func GetTitleWithByte(httpBody []byte) (t string, err error) {
	title := ""
	body:= GetByteReader(httpBody)

	node, err := html.Parse(body)
	if err != nil {
		return title, err
	}

	title, _ = getTitleText(node)

	return title, nil
}

// GetImgSrcWithBytesReader get all img urls
func GetImgSrcWithBytesReader(httpBody *bytes.Reader )(i *[]string, err error){
	ul:= []string{}
	page,err := html.Parse(httpBody)
	if err != nil{
		return &ul,err
	}
	ll , _ := getImgUrl(page,&ul)
	return ll,nil

}

// GetImgSrcWithByte
func GetImgSrcWithByte(httpBody []byte )(i *[]string, err error){
	var ul []string
	mm := GetByteReader(httpBody)

	page,err := html.Parse(mm)
	if err != nil{
		return &ul,err
	}
	ll , _ := getImgUrl(page,&ul)
	return ll,nil

}

// getImgUrl
func getImgUrl(node *html.Node, ad *[]string) (l *[]string, b bool) {
	flag := false
	if isImgElement(node){
		for _, v := range node.Attr{
			if v.Key == "src" {
				if check(ad, v.Val) == false {
					*ad = append(*ad, v.Val)
				}
			}
		}
		return ad ,true
	}

	for p:= node.FirstChild;p!=nil;p= p.NextSibling{
		ul,f := getImgUrl(p,ad)
		if f {
			flag = f
			ad = ul
		}

	}
	return ad, flag
}

// GetHrefWithBytesReader get all links
// It returns point []string
func GetHrefWithBytesReader(httpBody *bytes.Reader) (l *[]string, err error) {
	var links []string
	node, err := html.Parse(httpBody)
	if err != nil {
		return &links, err
	}
	ff, _ := getHref(node, &links)
	return ff, nil
}

// GetHrefWithByte
func GetHrefWithByte(httpBody []byte) (l *[]string, err error) {
	var links []string
	mm := GetByteReader(httpBody)
	node, err := html.Parse(mm)
	if err != nil {
		return &links, err
	}
	ff, _ := getHref(node, &links)
	return ff, nil
}


// getHref get url
func getHref(node *html.Node, h *[]string) (f *[]string, n bool) {
	b := false
	if isAHrefElement(node) {
		for _, a := range node.Attr {
			if a.Key == "href" {
				s := trimHash(a.Val)
				if check(h, s) == false {
					//*h = make([]string,0,len(*h)+len(s))
					//log.Println(cap(f))
					*h = append(*h, s)
				}

				return h, true
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		all, flag := getHref(c, h)
		h = all
		b = flag
	}
	return h, b
}

// TrimHash
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

// Check url exits
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
