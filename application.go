package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/solarknight/simple_golang_web/common"
	"github.com/solarknight/simple_golang_web/service"
)

type Controller func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

func main() {
	port := parseArgs()
	log.Println("Start on port", port)
	log.Println(http.ListenAndServe(":"+strconv.Itoa(port), initRouter()))
}

func parseArgs() (port int) {
	p := flag.Int("p", 8080, "port")
	flag.Parse()
	port = *p
	return
}

func initRouter() http.Handler {
	router := httprouter.New()
	router.GET("/", RestWrap(Index))
	router.GET("/hello", RestWrap(Hello))
	router.GET("/user/:id", RestWrap(GetUser))
	return router
}

func RestWrap(c Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		data, err := c(w, r, ps)
		if err != nil {
			log.Println("Error during process request", err)
			json.NewEncoder(w).Encode(&common.RestResponse{1, err.Error(), nil})
			return
		}

		json.NewEncoder(w).Encode(&common.RestResponse{0, "", data})
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, error) {
	return "Welcome!", nil
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return fmt.Sprintf("Hello, %s", name), nil
}

func GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (interface{}, error) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		return nil, err
	}
	return service.QueryByID(id)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, req *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, req *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
