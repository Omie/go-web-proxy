package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/p/", proxy)
	http.HandleFunc("/t/", proxyWget)
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

func proxyWget(res http.ResponseWriter, req *http.Request) {
	target := req.URL.Query().Get("target")
	if len(target) != 0 {
		//make temp dir
		tempPath, err := ioutil.TempDir("", "proxy")
		if err != nil {
			io.WriteString(res, err.Error())
			return
		}
		fmt.Println(tempPath)
		os.Chdir(tempPath)

		_, err = exec.Command("wget", "-p", target).Output()
		if err != nil {
			io.WriteString(res, "1003"+err.Error())
			return
		}
		u, _ := url.ParseRequestURI(target)
		h := u.Host
		fp := filepath.Join(tempPath, h)
		os.Chdir(fp)
		//fp = filepath.Join(fp, "index.html")
		//http.ServeFile(res, req, fp)
		go http.ListenAndServe("0.0.0.0:8080", http.FileServer(http.Dir(fp)))
		io.WriteString(res, "done")
	}
}
