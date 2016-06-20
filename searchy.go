// A sample RESTful API for search engines
// Used for tutorial on mattharris.org
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"golang.org/x/net/html"
)

// Handle RESTful Routing
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/google/{search_term}", GoogleIt)
	router.HandleFunc("/bing/{search_term}", BingIt)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle Google Searching
func GoogleIt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get("http://google.com/search?q=" + vars["search_term"])
	if err != nil {
		// Oh no! Google isn't working!!!
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		links := GetLinks(response)
		json.NewEncoder(w).Encode(links)
	}
}

// Handle Bing Searching
func BingIt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get("http://www.bing.com/search?q=" + vars["search_term"])
	if err != nil {
		// Oh no! Bing isn't working!!!
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		links := GetLinks(response)
		json.NewEncoder(w).Encode(links)
	}
}

// Parse out Links
func GetLinks(response *http.Response) (links []string) {
	z := html.NewTokenizer(response.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// EOD
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"

			if !isAnchor {
				continue
			}

			ok, url := GetHref(t)
			if !ok {
				continue
			}

			links = append(links, url)
		}
	}

	return
}

// Get the link
func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			matched, err := regexp.MatchString("^http", a.Val)
			if err != nil {
				ok = false
			} else {
				if matched {
					href = a.Val
					ok = true
				}
			}
		}
	}

	return
}
