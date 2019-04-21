package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const FILE_REDIRECT = "./redirect.yaml"

type redirectItem struct {
	key   string
	value string
}

type redirectList map[string]string

func main() {

	redirectKey := flag.String("a", "", "new redirect key")
	redirectURL := flag.String("u", "", "new redirect url")
	removeKey := flag.String("d", "", "Remove key from list")
	listRedirect := flag.Bool("l", false, "List redirect")
	port := flag.Int("p", 0, "Run server on port")
	printUsage := flag.Bool("h", false, "Print usage info")
	flag.Parse()

	if *printUsage == true {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *port != 0 {
		startServer(*port)
		return
	}

	list := getRedirectList()

	if *redirectKey != "" && *redirectURL != "" {
		list.append(redirectItem{key: *redirectKey, value: *redirectURL})
		return
	}

	if *removeKey != "" {
		list.remove(*removeKey)
		return
	}

	if *listRedirect == true {
		list.print()
		return
	}

}

func listen(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	key := strings.TrimPrefix(path, "/")
	list := getRedirectList()
	rediretURL := list[key]
	if rediretURL == "" {
		w.WriteHeader(404)
		w.Write([]byte("Not found redirect config"))
		return
	}

	http.Redirect(w, r, rediretURL, 301)
	return
}

func startServer(port int) {

	http.HandleFunc("/", listen)
	fmt.Printf("Starting server listen on %v", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func (list redirectList) append(item redirectItem) {
	list[item.key] = fmt.Sprintf("http://%v", item.value)
	d, _ := yaml.Marshal(&list)
	ioutil.WriteFile("redirect.yaml", d, 0644)
}

func (list redirectList) remove(key string) {
	delete(list, key)
	d, _ := yaml.Marshal(&list)
	ioutil.WriteFile(FILE_REDIRECT, d, 0644)
}

func (list redirectList) print() {
	d, err := yaml.Marshal(&list)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Print(string(d))
}

func getRedirectList() redirectList {
	yamlFile, err := ioutil.ReadFile(FILE_REDIRECT)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	m := make(map[string]string)

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return m
}
