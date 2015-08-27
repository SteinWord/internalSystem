package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	_rootdir    = "setup/first"
	_mkdirdir   = "setup/first"
)

type Server struct {
	Port string
}

type NewProject struct {
	Type string
	Name string
}

type Page struct {
	Result string
}

func randURL(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(36)]
	}
	return string(b)
}

func (n NewProject) Run() string {
	n.Name += ".9en.co"
	if n.Type == "beta" {
		n.Name = randURL(5) + "-" + n.Name
	}
	err := os.MkdirAll(_mkdirdir+"/"+n.Name+"/public_html", 0775)
	if err != nil {
		panic(err)
	}
	return n.Name
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	r.ParseForm()
	t := ""
	if r.Method == "POST" && r.Form["appName"][0] != "" {
		n := &NewProject{Name: r.Form["appName"][0], Type: r.Form["appType"][0]}
		t = n.Run()
	}
	html, _ := ioutil.ReadFile(_rootdir + "/tmpl/prot-design.html")
	data := Page{Result: t}
	tmpl, _ := template.New("page").Parse(string(html))
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", Handler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(_rootdir+"/css"))))
	http.ListenAndServe(":9000", nil)
}
