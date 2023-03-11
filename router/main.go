package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type routerParam map[string]string

type routerFunc func(routerParam, http.ResponseWriter, *http.Request)

type routerItem struct {
	method string
	matcher *regexp.Regexp
	fnc routerFunc
}

type router struct {
	items []routerItem
}

func (rt *router) GET(prefix string, fnc routerFunc) {
	rt.items = append(rt.items, routerItem{
		method: http.MethodGet,
		matcher: regexp.MustCompile(prefix),
		fnc: fnc,
	})
}

func (rt *router) POST(prefix string, fnc routerFunc) {
	rt.items = append(rt.items, routerItem{
		method: http.MethodPost,
		matcher: regexp.MustCompile(prefix),
		fnc: fnc,
	})
}

// http.Handlerのインターフェースを満たすために必要
func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, v := range rt.items {
		if v.method == r.Method && v.matcher.MatchString(r.RequestURI) {
			match := v.matcher.FindStringSubmatch(r.RequestURI)
			param := make(routerParam) // mapを初期化
			for i, name := range v.matcher.SubexpNames() {
				param[name] = match[i]
			}
			v.fnc(param, w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	rt := router{}
	rt.GET(`^/$`, func(p routerParam, w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%+v\n", p)
		fmt.Fprintln(w, "Hello")
	})
	rt.GET(`^/(?P<name>\w+)$`, func(p routerParam, w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n", p)
		fmt.Fprintf(w, "Hello: %v\n", p["name"])
	})
	rt.POST(`^/api$`, func(p routerParam, w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n", p)
		w.Header().Set("content-type", "text/json")
		fmt.Fprintln(w, `{"status": "OK"}`)
	})
	http.ListenAndServe(":8080", &rt)
}