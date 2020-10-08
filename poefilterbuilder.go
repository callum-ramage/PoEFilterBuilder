package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"./filters"
	"./replacement"
	"./resources"

	"github.com/callum-ramage/jsonconfig"
	"github.com/pkg/browser"
)

func servePage(file string, templateData interface{}, w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.ParseFiles(resources.ServerRoot() + file)
	if err != nil {
		log.Fatal(err)
	}

	err = homeTemplate.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func serveFile(base string, w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	// fmt.Println(upath)
	if len(upath) == 1 || upath[len(upath)-1] == '/' {
		servePage(base+upath+"Home.html", nil, w, r)
		return
	}
	// fmt.Println(resources.ServerRoot() + filepath.FromSlash(base))
	dir := http.Dir(resources.ServerRoot() + filepath.FromSlash(base))
	_, err := dir.Open(path.Clean(upath))
	if err != nil {
		forohfortemplate, err := template.ParseFiles(resources.ServerRoot() + "404.html")
		if err != nil {
			return
		}

		forohfortemplate.Execute(w, nil)
		return
	}
	http.FileServer(dir).ServeHTTP(w, r)
}

func respond(w http.ResponseWriter, r *http.Request, content []byte) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		// fmt.Fprintf(w, "%+v", string(value))
		encodeGzip(w, content)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%+v", string(content))
	}
}

func encodeGzip(w http.ResponseWriter, content []byte) {
	gz := gzip.NewWriter(w)
	defer gz.Close()

	_, err := gz.Write(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func localTestHandler(w http.ResponseWriter, r *http.Request) {
	servePage("Test.html", nil, w, r)
}

func localCloseHandler(w http.ResponseWriter, r *http.Request) {
	servePage("Close.html", nil, w, r)
	go func() {
		time.Sleep(1000 * time.Millisecond)
		log.Fatal("Shutdown")
	}()
}

func localHomeHandler(w http.ResponseWriter, r *http.Request) {
	servePage("Home.html", nil, w, r)
}

func localReplaceHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.MarshalIndent(replacement.Replacements, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	type ConfigData struct {
		// FileNames []string
		Files string
	}
	// servePage("Replace.html", ConfigData{FileNames: replacement.GetReplacementKeys()}, w, r)
	servePage("Replace.html", ConfigData{Files: string(json)}, w, r)
}

func localSubmitReplaceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileName := r.FormValue("filename")
	fileContent := r.FormValue("replaceArea")
	// fmt.Println(r.FormValue("submit"))
	// fmt.Println(fileName)
	// fmt.Println(fileContent)
	switch strings.ToLower(r.FormValue("submit")) {
	case "save":
		replacement.SaveReplacement(fileName, fileContent)
	case "delete":
		replacement.DeleteReplacement(fileName)
	}

	localReplaceHandler(w, r)
	// servePage("Replace.html", ConfigData{JSON: string(json)}, w, r)
}

func localFilterHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.MarshalIndent(filters.Filters, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	type ConfigData struct {
		// FileNames []string
		Files string
	}
	// servePage("Replace.html", ConfigData{FileNames: replacement.GetReplacementKeys()}, w, r)
	servePage("Filter.html", ConfigData{Files: string(json)}, w, r)
}

func localSubmitFilterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileName := r.FormValue("filename")
	fileContent := r.FormValue("filterArea")
	// fmt.Println(r.FormValue("submit"))
	// fmt.Println(fileName)
	// fmt.Println(fileContent)
	switch strings.ToLower(r.FormValue("submit")) {
	case "save":
		filters.SaveFilter(fileName, fileContent)
	case "delete":
		filters.DeleteFilter(fileName)
	case "compile":
		filters.SaveFilter(fileName, fileContent)
		filters.CompileFilter(fileName)
	}

	localFilterHandler(w, r)
	// servePage("Replace.html", ConfigData{JSON: string(json)}, w, r)
}

func localConfigHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	type ConfigData struct {
		JSON string
	}

	servePage("Config.html", ConfigData{JSON: string(json)}, w, r)
	// respond(w, r, json)
}

func localSaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	saveConfig := false

	if r.FormValue("replacements") != config.ReplaceFiles {
		saveConfig = true
		config.ReplaceFiles = r.FormValue("replacements")
		replacement.LoadReplacements(config.ReplaceFiles)
	}

	if r.FormValue("filters") != config.FilterFiles {
		saveConfig = true
		config.FilterFiles = r.FormValue("filters")
	}

	if r.FormValue("output") != config.OutputLocation {
		saveConfig = true
		config.OutputLocation = r.FormValue("output")
	}

	// if r.FormValue("") != config. {
	// 	saveConfig = true
	// 	config. = r.FormValue("")
	// }

	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if saveConfig {
		err = ioutil.WriteFile("./config.conf", json, 0644)
		if err != nil {
			log.Fatal(err)
		}
		loadConfig()
	}

	type ConfigData struct {
		JSON string
	}

	servePage("Config.html", ConfigData{JSON: string(json)}, w, r)
	// respond(w, r, json)
}

func localHandlerFunc(w http.ResponseWriter, r *http.Request) {
	serveFile("", w, r)
}

type localHandler map[string]http.Handler

func (h localHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.Split(r.Host, ":")
	domainParts = strings.Split(domainParts[0], ".")

	r.URL, _ = url.Parse(r.URL.String())

	mux := h[strings.ToLower(domainParts[0])]
	if mux != nil {
		mux.ServeHTTP(w, r)
	} else {
		// http.FileServer(http.Dir(resources.ServerRoot()+"404.html")).ServeHTTP(w, r)
	}
}

type configuration struct {
	ReplaceFiles   string
	FilterFiles    string
	OutputLocation string
}

var config configuration

func loadConfig() {
	replacement.LoadReplacements(config.ReplaceFiles)
	filters.LoadFilters(config.FilterFiles, config.OutputLocation)
}

func main() {
	var err error

	config = configuration{
		ReplaceFiles:   "./replacements/",
		OutputLocation: "./output/",
		FilterFiles:    "./filters/",
	}
	err = jsonconfig.Load("./config.conf", &config)

	if err != nil {
		log.Fatal(err)
	}

	loadConfig()

	fmt.Println("ReplaceFiles:", config.ReplaceFiles)
	fmt.Println("FilterFiles:", config.FilterFiles)

	localMux := http.NewServeMux()
	localMux.HandleFunc("/test", localTestHandler)
	localMux.HandleFunc("/close", localCloseHandler)
	localMux.HandleFunc("/home", localHomeHandler)
	localMux.HandleFunc("/config", localConfigHandler)
	localMux.HandleFunc("/saveconfig", localSaveConfigHandler)
	localMux.HandleFunc("/replace", localReplaceHandler)
	localMux.HandleFunc("/submitreplace", localSubmitReplaceHandler)
	localMux.HandleFunc("/filter", localFilterHandler)
	localMux.HandleFunc("/submitfilter", localSubmitFilterHandler)
	localMux.HandleFunc("/", localHandlerFunc)

	globalHandler := make(localHandler)
	globalHandler["localhost"] = localMux

	browser.OpenURL("http://localhost:3000/home")

	var addr = 3000
	for {
		fmt.Printf("Trying %d\n", addr)
		tempServerVar := &http.Server{Addr: ":" + strconv.Itoa(addr), Handler: globalHandler}
		err = tempServerVar.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			addr++
		}
	}
}
