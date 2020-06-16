package alink

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
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

func TestGetHrefWithBytesReader(t *testing.T) {
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
	mm, _ := GetHrefWithBytesReader(b)
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

func TestGetImgSrcWithBytesReader(t *testing.T) {
	type args struct {
		httpBody *bytes.Reader
	}


	var html = `<a href="http://jjjj.com">1</a> <video src="http://abc.com/ab.mp4">
  <a style=\"\" href=http://imgur.com>3</a> <img src="abc.com/img.jpg">http://alink.com</p>`

	c := []byte(html)
	i :=  args{
		bytes.NewReader(c),
	}

	var html1 = `<a href="http://jjjj.com">1</a> <video src="http://abc.com/ab.mp4">
  <a style=\"\" href=http://imgur.com>3</a> <img lin="abc.com/img.jpg">http://alink.com</p>`

	c1 := []byte(html1)
	i1 :=  args{
		bytes.NewReader(c1),
	}

	var tests =[]struct {
		name string
		args args
		wantS *[]string
		wantErr bool
	}{
		{"img",i,&[]string{"abc.com/img.jpg"},false },
		{"imgNoSrc",i1,&[]string{},false },
	}
	for _, tt :=range tests{
		t.Run(tt.name,func(t *testing.T){
			gotS ,err:= GetImgSrcWithBytesReader(tt.args.httpBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImgSrcWithBytesReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("GetImgSrcWithBytesReader() gotS = %v, want %v", gotS, tt.wantS)
			}
		})
	}

}

func TestGetBytesReaderWithIoReader(t *testing.T) {
	s := "<div></div>"
	str := []byte(s)
	reader := bytes.NewReader(str)

	abc, err := GetBytesReaderWithIoReader(reader)
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

func TestGetVideoSrcWithBytesReader(t *testing.T) {
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
		wantS   []string
		wantErr bool
	}{
		{"video", f, []string{"http://abc.com/ab.mp4"}, false},
	}

	//log.Print(tests[0].wantS)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := GetVideoSrcWithBytesReader(tt.args.httpBody)

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

func TestTitleBytes(t *testing.T) {

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
			gotT, err := TitleBytes(tt.args.httpBody)
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
		GetHrefWithBytesReader(bytes.NewReader(p))
	}
}



func TestGetImgSrcWithByte(t *testing.T) {
	type args struct {
		httpBody []byte
	}
	var reader = `<img src="http://jjjj.com"> 
   <video src="http://abc.com/ab.mp4">
  <a style=\"\" href=http://imgur.com>3</a>
  http://alink.com
</p>`

	 //string to byte.reader
	  httpBody :=args{[]byte(reader) }

	  tests := []struct {
		name    string
		args    args
		wantI   *[]string
		wantErr bool
	}{
		{"img",httpBody, &[]string{"http://jjjj.com"},false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, err := GetImgSrcWithByte(tt.args.httpBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImgSrcWithByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("GetImgSrcWithByte() gotI = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

//go test -cover -v -coverprofile=c.out

func TestGetTitleWithByte(t *testing.T) {
	type args struct {
		httpBody []byte
	}

	page1 := "<html><header><title>test1</title></header><body></body></html>"
	p1 := []byte(page1)
	a1 := args{p1}

	tests := []struct {
		name    string
		args    args
		wantT   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"title",a1,"test1",false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := GetTitleWithByte(tt.args.httpBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTitleWithByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotT != tt.wantT {
				t.Errorf("GetTitleWithByte() gotT = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestGetByteWithIoReader(t *testing.T) {
	type args struct {
		respBody io.Reader
	}

	af :=[]byte("abc")

	body := args{
		bytes.NewReader(af),
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"readToByte",body,[]byte{97,98,99},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByteWithIoReader(tt.args.respBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByteWithIoReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByteWithIoReader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestGetHrefWithByte(t *testing.T) {
//	type args struct {
//		httpBody []byte
//	}
//	html :=`<p>test<a href="test.com">one</a></p>`
//	html1 :=`<p>test <a href="#">one</a></p>`
//
//	h := []byte(html)
//	h1 := []byte(html1)
//	arr := args{
//		h,
//	}
//	arr1 := args{h1}
//	tests := []struct {
//		name  string
//		args  args
//		wantL *[]string
//		wantB bool
//	}{
//
//		{"one",arr,&[]string{"test.com"},true},
//		{"two",arr1,&[]string{""},true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotL, gotB := GetHrefWithByte(tt.args.httpBody)
//			if !reflect.DeepEqual(gotL, tt.wantL) {
//				t.Errorf("GetHrefWithByte() gotL = %v, want %v", gotL, tt.wantL)
//			}
//			if gotB != tt.wantB {
//				t.Errorf("GetHrefWithByte() gotB = %v, want %v", gotB, tt.wantB)
//			}
//		})
//	}
//}

func TestGetHrefWithByte1(t *testing.T) {
	type args struct {
		httpBody []byte
	}
	html :=`<p>test<a href="test.com">one</a></p>`
	html1 :=`<p>test <a href="#">one</a></p>`

	h := []byte(html)
	h1 := []byte(html1)
	arr := args{
		h,
	}
	arr1 := args{h1}
	tests := []struct {
		name    string
		args    args
		wantL   *[]string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"one",arr,&[]string{"test.com"},false},
		{"two",arr1,&[]string{""},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL, err := GetHrefWithByte(tt.args.httpBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHrefWithByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("GetHrefWithByte() gotL = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}