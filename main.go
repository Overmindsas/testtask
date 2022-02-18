package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type ProxySettings struct {
	Addres string `json:"addres"`
	Port   string `json:"port"`
}

func Proxy(addres string) *url.URL {

	url, err := url.Parse(addres)
	if err != nil {
		fmt.Println(err)
	}
	return url
}

func main() {

	var pr ProxySettings

	cache := make(map[string]string)

	dir, err := filepath.Abs("conf.json")
	if err != nil {
		fmt.Println(err)
	}

	file, _ := ioutil.ReadFile(dir)

	err = json.Unmarshal([]byte(file), &pr)
	if err != nil {
		fmt.Println(err)
	}

	data := ProxySettings{
		Addres: pr.Addres,
		Port:   pr.Port,
	}

	url := Proxy(data.Addres)
	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			cache[""] = data.Addres
			log.Println(r.URL)
			r.Host = url.Host
			w.Header().Set("X-Ben", "Rad")
			p.ServeHTTP(w, r)
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":"+data.Port, nil)
	if err != nil {
		panic(err)
	}
}
