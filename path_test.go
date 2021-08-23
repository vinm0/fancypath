package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ReqReader struct{}

func (rr ReqReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// func TestParse(t *testing.T) {
// 	r := httptest.NewRequest(http.MethodGet, "/", ReqReader{})
// 	r1 := httptest.NewRequest(http.MethodGet, "/blog", ReqReader{})
// 	r2 := httptest.NewRequest(http.MethodGet, "/blog/new", ReqReader{})
// 	r3 := httptest.NewRequest(http.MethodGet, "/blog/new/article", ReqReader{})
// 	r4 := httptest.NewRequest(http.MethodGet, "/blog/new/article/42", ReqReader{})

// 	var tests = []struct {
// 		s    *[]string
// 		r    *http.Request
// 		path string
// 	}{
// 		{&[]string{}, r, ""},
// 		{&[]string{}, r, "/"},
// 		{&[]string{}, r, "/*"},
// 		{&[]string{}, r, "/*/action"},
// 		{&[]string{}, r, "/category/action/item/id"},
// 		{&[]string{}, r1, ""},
// 		{&[]string{}, r1, "/"},
// 		{&[]string{}, r1, "/*"},
// 		{&[]string{}, r1, "/*/action"},
// 		{&[]string{}, r1, "/category/action/item/id"},
// 		{&[]string{}, r2, ""},
// 		{&[]string{}, r2, "/"},
// 		{&[]string{}, r2, "/*"},
// 		{&[]string{}, r2, "/*/action"},
// 		{&[]string{}, r2, "/category/action/item/id"},
// 		{&[]string{}, r3, ""},
// 		{&[]string{}, r3, "/"},
// 		{&[]string{}, r3, "/*"},
// 		{&[]string{}, r3, "/*/action"},
// 		{&[]string{}, r3, "/category/action/item/id"},
// 		{&[]string{}, r4, ""},
// 		{&[]string{}, r4, "/"},
// 		{&[]string{}, r4, "/*"},
// 		{&[]string{}, r4, "/*/action"},
// 		{&[]string{}, r4, "/category/action/item/id"},
// 	}
// }

func TestVar(t *testing.T) {
	p := Newpath(httptest.NewRequest(http.MethodGet, "/", ReqReader{}), "//*/category/id")
	p1 := Newpath(httptest.NewRequest(http.MethodGet, "/blog", ReqReader{}), "//*/category/id")
	p2 := Newpath(httptest.NewRequest(http.MethodGet, "/blog/new", ReqReader{}), "//*/category/id")
	p3 := Newpath(httptest.NewRequest(http.MethodGet, "/blog/new/article", ReqReader{}), "//*/category/id")
	p4 := Newpath(httptest.NewRequest(http.MethodGet, "/blog/new/article/42", ReqReader{}), "//*/category/id")
	p5 := Newpath(httptest.NewRequest(http.MethodGet, "/blog/new/article/42/extra", ReqReader{}), "//*/category/id")

	var tests = []struct {
		p    *path
		key  string
		want string
	}{
		{p, "", ""},
		{p, "*", ""},
		{p, "action", ""},
		{p, "category", ""},
		{p, "id", ""},

		{p1, "", ""},
		{p1, "*", ""},
		{p1, "action", ""},
		{p1, "category", ""},
		{p1, "id", ""},

		{p2, "", ""},
		{p2, "*", ""},
		{p2, "action", ""},
		{p2, "category", ""},
		{p2, "id", ""},

		{p3, "", ""},
		{p3, "*", ""},
		{p3, "action", ""},
		{p3, "category", "article"},
		{p3, "id", ""},

		{p4, "", ""},
		{p4, "*", ""},
		{p4, "action", ""},
		{p4, "category", "article"},
		{p4, "id", "42"},

		{p5, "", ""},
		{p5, "*", ""},
		{p5, "action", ""},
		{p5, "category", "article"},
		{p5, "id", "42"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v,%s", tt.p, tt.key)

		t.Run(testname, func(t *testing.T) {
			ans := tt.p.Var(tt.key)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {

	var tests = []struct {
		s    *[]string
		path string
		want int
	}{
		{&[]string{}, "", 0},
		{&[]string{}, "/", 0},
		{&[]string{}, "/blog", 1},
		{&[]string{}, "/blog/", 1},
		{&[]string{}, "/blog/new", 2},
		{&[]string{}, "/blog/new/article", 3},
		{&[]string{}, "/blog/new/article/42", 4},

		{&[]string{}, "/*", 1},
		{&[]string{}, "/*/", 1},
		{&[]string{}, "/*//", 2},
		{&[]string{}, "/*//category", 3},
		{&[]string{}, "/*//category/id", 4},
	}

	for _, tt := range tests {
		testname := tt.path

		t.Run(testname, func(t *testing.T) {
			parse(tt.s, tt.path)
			ans := len(*tt.s)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
