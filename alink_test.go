package alink

import (
	"bytes"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"testing"
)

func TestIsValidUrl(t *testing.T) {

	str := "https://www.google.com/ab"
	if IsValidUrl(str) != true {
		t.Error(`IsValidUrl("https://www.google.com/ab")==false`)
	}

	str5 := "http://ab.f.f.nc.f.-._"
	if !IsValidUrl(str5) {
		t.Error(`IsValidUrl("http://ab.f.f.nc.f.-._")==false`)
	}
}

func TestNotIsValidUrl(t *testing.T) {

	var str1 = "/category.php?id=95929"
	str2 := "f.abc.co.cf"
	str3 := "w.b.com"
	str4 := "http:/ab.com.cn"
	str5 := "www.exp.com"
	expected := false
	actual := IsValidUrl(str1)
	if expected != actual {
		t.Error(`IsValidUrl("category.php?id=95929")==false`)
	}
	//t.Log(IsValidUrl(str1) )
	if expected != IsValidUrl(str2) {
		t.Error(`IsValidUrl("f.abc.co.cf")==true`)
	}
	if expected != IsValidUrl(str3) {
		t.Error(`IsValidUrl("w.b.com")==true`)
	}
	if IsValidUrl(str4) != expected {
		t.Error(`IsValidUrl("http:/ab.com.cn")==true`)
	}
	if IsValidUrl(str5) != expected {
		t.Error(`IsValidUrl("www.exp.com")==true`)
	}
}

func TestAlink(t *testing.T) {
	var reader = `<a href="http://jjjj.com">1</a>
  <a href='http://news.google.com'>2</a>
  <a style=\"\" href=http://imgur.com>3</a>
  http://alink.com
</p>`
	var links []string
	var tmp int
	//string to byte.reader
	c := []byte(reader)
	b := bytes.NewReader(c)
	mm, _ := Alink(b)
	//log.Print(mm)
	for i, k := range *mm {
		links = append(links, k)
		tmp = i
	}
	//log.Print(links,tmp)
	if tmp != 2 || len(links) != 3 {
		t.Error("Wrong number of links returned")
	}
	if links[0] != "http://jjjj.com" {
		t.Error("The first link is incorrect")
	}
	if links[1] != "http://news.google.com" {
		t.Error("The second link is incorrect")
	}
	if links[2] != "http://imgur.com" {
		t.Error("The third link is incorrect")
	}
}

func TestNewRespBody(t *testing.T) {
	s := "<div></div>"
	str := []byte(s)
	reader := bytes.NewReader(str)
	//str := strings.NewReader(s)
	abc, err := NewRespBody(reader)
	if err != nil {
		t.Error(err)
	}
	cc := html.NewTokenizer(abc)

	tokenType := cc.Next()

	token := cc.Token()
	log.Print(tokenType, token)
	if token.DataAtom.String() != "div" {
		t.Error("Wrong body ")
	}
}

func TestVideo(t *testing.T) {
	var reader = `<a href="http://jjjj.com">1</a>
   <video src="http://abc.com/ab.mp4">
  <a style=\"\" href=http://imgur.com>3</a>
  http://alink.com
</p>`

	//string to byte.reader
	c := []byte(reader)
	//b := bytes.NewReader(c)
	type args struct {
		httpBody *bytes.Reader
	}

	f := args{
		bytes.NewReader(c),
	}

	var tests = []struct {
		name    string
		args    args
		wantS   [] string
		wantErr bool
	}{
		{"video", f, []string{"http://abc.com/ab.mp4"}, false},
	}

	//log.Print(tests[0].wantS)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := VideoSrc(tt.args.httpBody)

			if (err != nil) != tt.wantErr {
				t.Errorf("Video() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("Video() gotS = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestTitle(t *testing.T) {

	type args struct {
		httpBody *bytes.Reader
	}

	page := "<html><header><title>the title</title></header><body></body></html>"
	p := []byte(page)
	a := args{bytes.NewReader(p)}

	page1 := "<html><header><title>test1</title></header><body></body></html>"
	p1 := []byte(page1)
	a1 := args{bytes.NewReader(p1)}

	tests := []struct {
		name    string
		args    args
		wantT   string
		wantErr bool
	}{
		{"title", a, "the title", false},
		{"title", a1, "test1", false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := Title(tt.args.httpBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("Title() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotT != tt.wantT {
				t.Errorf("Title() gotT = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func BenchmarkAlink(b *testing.B) {
	page := "<html><header><title>the title</title></header><body></body></html>"
	p := []byte(page)

	for i := 0; i < b.N; i++ {
		Alink(bytes.NewReader(p))
	}
}

//go test -cover -v -coverprofile=c.out
