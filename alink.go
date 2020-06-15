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

// NewRespBody is the func create a new io.reader body
// It returns a point of bytes.Reader
func NewRespBody(respBody io.Reader) (*bytes.Reader, error) {

	b, err := ioutil.ReadAll(respBody)
	reader := bytes.NewReader(b)
	return reader, err
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
func titleText(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		//log.Print(n)
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

// videoSrc get video src
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

// VideoSrc get the video tags src
// It returns []string
func VideoSrc(httpBody *bytes.Reader) (s [] string, err error) {
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

// Title to get pages title return a string
func Title(httpBody *bytes.Reader) (t string, err error) {
	title := ""
	node, err := html.Parse(httpBody)
	if err != nil {
		return title,err
	}

	title, _ = titleText(node)

	return title, nil
}

// Alink get all links
// It returns point []string and a bool value to check the page has a tags
func Alink(httpBody *bytes.Reader) (l *[]string, b bool) {
	var links []string
	node, err := html.Parse(httpBody)
	if err != nil {
		return &links,false
	}
	ff, _ := alLink(node, &links)
	return ff, true
}

// alLink Get href url
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