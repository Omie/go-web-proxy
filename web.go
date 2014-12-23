package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/p/", proxy)
	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello, world from %s", runtime.Version())
}

func proxy(res http.ResponseWriter, req *http.Request) {
	target := req.URL.Query().Get("target")
	if len(target) != 0 {
		client := &http.Client{}
		r, err := http.NewRequest("GET", target, nil)
		if err != nil {
			io.WriteString(res, err.Error())
			return
		}

		resp, err := client.Do(r)
		if err != nil {
			io.WriteString(res, err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			io.WriteString(res, err.Error())
			return
		}
		io.WriteString(res, string(body))
	}
}
