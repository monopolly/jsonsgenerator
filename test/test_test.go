package news

//testing

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Marshal(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	// js=NewsJson sql=news ts=news swift

	fmt.Println(Between("{{tag1 value1}} {{tag2 value2}}", "{{tag1 ", "}}"))
	fmt.Println(Between("{{tag1 value1}} {{tag2 value2}}", "{{tag2 ", "}}"))

	var r NewsQuery
	r.Fields = []string{
		IndexKeys.SQLName(),
		IndexTitle.SQLName(),
		IndexCreated.SQLName(),
		IndexCount.SQLName(),
		IndexActive.SQLName(),
		IndexKYC.SQLName(),
	}
	r.EQ = map[string]any{
		"title":  "James",
		"active": true,
		"kyc":    false,
	}
	r.GT = map[string]any{
		"created": 194859789278,
	}
	r.LT = map[string]any{
		"count": 200,
	}
	r.Desc = true
	r.Like = map[string]string{
		"title": "max%",
	}
	r.Limit = 100
	r.Offset = 50
	r.Sort = "title"

	sql, fields, v := r.Render()

	fmt.Println(sql)
	fmt.Println(fields)
	fmt.Println(v)

}

func TestGotiny(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	p := testStructure()
	b := p.Gotiny()
	fmt.Println("gotiny size", len(b))

	p1 := ParseNewsGotiny(b)

	if p1.AccountID != 111 {
		panic("gotiny id")
	}

}

func TestMessagePack(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	p := testStructure()
	b := p.MessagePack()

	fmt.Println("msgp size", len(b))

	p1, err := ParseNewsMessagePack(b)
	if err != nil {
		panic(err)
	}

	if p1.AccountID != 111 {
		panic("msgp id")
	}

}

func TestJSONFast(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	p := testStructure()
	b := p.Pack()

	fmt.Println("jsonfast size", len(b))
	os.WriteFile("test_test_fastjson.json", b, os.ModePerm)

	n, err := ParseNews(b)
	if err != nil {
		panic(err)
	}

	if n.AccountID != 111 {
		panic("fastjson id")
	}

}

func TestToJSON(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	p := testStructure()
	b := p.ToJson()

	fmt.Println("toJson size", len(b))
	os.WriteFile("test_test_to.json", b, os.ModePerm)

}

func testStructure() News {

	var p News
	p.AccountID = 111
	p.Active = true
	p.BID = 24552
	p.Features = map[string]bool{"some": true}
	p.Floats = 0.5243224
	p.Html = []byte("nice oevinoen elfrj sljalj alkj alfkaldknf ,adlfk nadflknadm. adlknadflknadlfknaldknad,m adflknadlfknadf")
	p.Tags = []string{"131441", "fksjfkh", "sljfljslj", "42j4lkj"}
	return p
}

func BenchmarkJSON(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.ReportAllocs()
	bbbbbbbb.ResetTimer()
	for n := 0; n < bbbbbbbb.N; n++ {
		p.ToJson()
	}
}

func BenchmarkJSONParallel(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.ToJson()
		}
	})
}

func BenchmarkJSONFast(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.ReportAllocs()
	bbbbbbbb.ResetTimer()
	for n := 0; n < bbbbbbbb.N; n++ {
		p.Pack()
	}
}

func BenchmarkJSONFastParallel(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Pack()
		}
	})
}

func BenchmarkGotiny(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.ReportAllocs()
	bbbbbbbb.ResetTimer()
	for n := 0; n < bbbbbbbb.N; n++ {
		p.Gotiny()
	}
}

func BenchmarkGotinyParallel(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Gotiny()
		}
	})
}

func BenchmarkMsgPack(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.ReportAllocs()
	bbbbbbbb.ResetTimer()
	for n := 0; n < bbbbbbbb.N; n++ {
		p.MessagePack()
	}
}

func BenchmarkMsgPackParallel(bbbbbbbb *testing.B) {
	p := testStructure()
	bbbbbbbb.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.MessagePack()
		}
	})
}

// any {{tag value}}
func Between(raw string, startTag, endTag string) (value string) {
	indexStart := strings.Index(raw, startTag)
	if indexStart == -1 {
		return
	}

	from := indexStart + len(startTag)
	indexEnd := strings.Index(raw[from:], endTag)
	if indexEnd == -1 {
		return
	}
	to := from + indexEnd //+ len(endTag)

	return raw[from:to]
}
