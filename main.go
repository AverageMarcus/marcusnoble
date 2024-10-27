package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

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

	// Set up "slugs" that can be used as redirects to social media accounts
	// E.g. https://marcusnoble.com/mastodon -> https://k8s.social/@Marcus
	for _, l := range data["social"].([]interface{}) {
		link := l.(map[interface{}]interface{})
		if link["slug"] != "" {
			http.Handle(link["slug"].(string), http.RedirectHandler(link["url"].(string), http.StatusTemporaryRedirect))
		}
	}

	// Filter out any events that have passed already
	futureEvents := []map[interface{}]interface{}{}
	dateLayout := "2006-01-02"
	for _, e := range data["events"].([]interface{}) {
		event := e.(map[interface{}]interface{})
		t, err := time.Parse(dateLayout, event["date"].(string))
		if err == nil && time.Now().Before(t) {
			futureEvents = append(futureEvents, event)
		}
	}
	data["events"] = futureEvents

	funcMap := template.FuncMap(map[string]interface{}{
		"join": func(objs []interface{}, key, joiner string) template.HTML {
			vals := []string{}
			for _, obj := range objs {
				val := obj.(map[interface{}]interface{})[key]
				vals = append(vals, val.(string))
			}
			return template.HTML(strings.Join(vals, joiner))
		},
		"html": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := "src" + r.URL.Path
		if strings.HasSuffix(page, "/") {
			page = page + "index.html"
		}

		if strings.HasSuffix(page, ".html") || strings.HasSuffix(page, ".md") {
			tpl, err := template.New(path.Base(page)).Funcs(funcMap).ParseFS(res, page)
			if err != nil {
				log.Printf("page %s (%s) not found...", r.RequestURI, page)
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			if err := tpl.Execute(w, data); err != nil {
				log.Println(err)
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
