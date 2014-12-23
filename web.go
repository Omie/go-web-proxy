package main

import (
	"crypto/rand"
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
	/* simple function that will fetch single resource and tunnel the response
	   useful for simple txt or image files
	*/
	target := req.URL.Query().Get("target")
	if len(target) != 0 {
		client := &http.Client{}
		r, err := http.NewRequest("GET", target, nil)
		if err != nil {
			io.WriteString(res, "1001 "+err.Error())
			return
		}

		resp, err := client.Do(r)
		if err != nil {
			io.WriteString(res, "1002 "+err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			io.WriteString(res, "1003 "+err.Error())
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
			io.WriteString(res, "1004 "+err.Error())
			return
		}
		os.Chdir(tempPath) //do this to make sure wget will fetch contents in this dir

		_, err = exec.Command("wget", "-p", target).Output()
		if err != nil {
			io.WriteString(res, "1005 "+err.Error())
			return
		}
		u, _ := url.ParseRequestURI(target)
		h := u.Host
		fp := filepath.Join(tempPath, h)
		b := "/" + rand_str(15) + "/"
		http.Handle(b, http.StripPrefix(b, http.FileServer(http.Dir(fp))))
		http.Redirect(res, req, b, 302)
	}
}

func rand_str(str_size int) string {
	//https://devpy.wordpress.com/2013/10/24/create-random-string-in-golang/
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
