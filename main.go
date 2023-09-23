package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

func main() {
	cats := crawl("https://airtable.com/appoPqE4I6DudZGIe/shr3aL7FonYv3m1b2/tblO7PpHl1a0mFxD4")
	
	for _, cat := range cats {
		printCat(cat)
		fmt.Println()
	}
	fmt.Printf("Found %d cats\n", len(cats))
}

type cat struct {
	name         string
	pics         []string
	adoptionLink string
}

func crawl(url string) []cat {
	// Launch a new browser with default options, and connect to it
	browser := rod.New().MustConnect()

	// Close it after main process ends
	defer browser.Close()

	// Create a new page
	page := browser.MustPage(url).MustWaitStable()

	// Zoom out to get more cats
	page.MustEval(`() => document.body.style.zoom = "10%"`)
	page.MustWaitStable()

	cats := []cat{}

	catDivs := page.MustElements("div.baymaxGalleryCard");
	for _, catDiv := range catDivs {	
		nameA := catDiv.MustElement("a.galleryCardPrimaryCell")
		name := nameA.MustText()
		adoptionLink := nameA.MustProperty("href").String()
		
		pics := []string{}
		for _, image := range catDiv.MustElements(".coverAttachmentCell link") {
			pics = append(pics, image.MustProperty("href").String())
		}
		
		cats = append(cats, cat{name, pics, adoptionLink})
	}
	
	return cats;
}


func printCat(cat cat) {
	fmt.Printf("Name: %s\n", cat.name)
	fmt.Printf("Link: %s\n", cat.adoptionLink)
	
	fmt.Println("Pics:")
	for _, pic := range cat.pics {
		fmt.Println(pic)
	}
}
