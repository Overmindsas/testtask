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

func Proxy(addres string) (*httputil.ReverseProxy, error) {

	url, err := url.Parse(addres)
	if err != nil {
		fmt.Println(err)
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func Cache() {

}

func main() {
	var p ProxySettings
	dir, err := filepath.Abs("conf.json")
	if err != nil {
		fmt.Println(err)
	}

	file, _ := ioutil.ReadFile(dir)

	err = json.Unmarshal([]byte(file), &p)
	if err != nil {
		fmt.Println(err)
	}

	data := ProxySettings{
		Addres: p.Addres,
		Port:   p.Port,
	}
	newproxy, err := Proxy(data.Addres)

	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newproxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":"+data.Port, nil))
}
