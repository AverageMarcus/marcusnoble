package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	//go:embed src data.yaml
	res  embed.FS
	port string
)

func init() {
	var ok bool
	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}
}

func main() {
	dataBytes, err := res.ReadFile("data.yaml")
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	if err := yaml.Unmarshal(dataBytes, &data); err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := "src" + r.URL.Path
		if strings.HasSuffix(page, "/") {
			page = page + "index.html"
		}

		if strings.HasSuffix(page, ".html") {
			tpl, err := template.ParseFS(res, page)
			if err != nil {
				log.Printf("page %s (%s) not found...", r.RequestURI, page)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

			tpl.Funcs(template.FuncMap(map[string]interface{}{
				"html": func(str string) template.HTML {
					return template.HTML(str)
				},
			}))

			if err := tpl.Execute(w, data); err != nil {
				return
			}
		} else {

			// Serve up the best file format for image
			if strings.Contains(page, "headshot-transparent.png") {
				if strings.Contains(r.Header.Get("Accept"), "image/avif") {
					page = strings.Replace(page, ".png", ".avif", 1)
				} else if strings.Contains(r.Header.Get("Accept"), "image/webp") {
					page = strings.Replace(page, ".png", ".webp", 1)
				}
			}

			body, err := res.ReadFile(page)
			if err != nil {
				log.Printf("file %s (%s) not found...", r.RequestURI, page)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(body)
			return
		}
	})
	// http.FileServer(http.FS(res))

	fmt.Println("Server started at port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
