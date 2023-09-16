package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	_, err := crawl("https://airtable.com/appoPqE4I6DudZGIe/shr3aL7FonYv3m1b2/tblO7PpHl1a0mFxD4")
	if err != nil {
		log.Fatalf("could not crawl the url: %v", err)
	}
}

type record struct {
	name         string
	pics         []string
	adoptionLink string
}

func crawl(url string) ([]record, error) {
	// get the html of webpage
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error: could not Get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: Status code %d", resp.StatusCode)

	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}
	// go through each element in DOM
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nil, nil
	// look up the element
}
