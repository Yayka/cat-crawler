package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-rod/rod"
)

func main() {
	cats := crawl("https://airtable.com/appoPqE4I6DudZGIe/shr3aL7FonYv3m1b2/tblO7PpHl1a0mFxD4")

	for _, cat := range cats {
		printCat(cat)
		fmt.Println()
	}
	fmt.Printf("\nFound %d cats\n", len(cats))

	randomCat := cats[rand.Intn(len(cats))]
	randomPic := randomCat.pics[rand.Intn(len(randomCat.pics))]

	fmt.Printf("\n\nChecking %s\n", randomPic)
	
	annotations, err := Annotate(randomPic)
	if err != nil {
		fmt.Println(err)
		return;
	}

	var isCat = false
	for i := 0; i < len(annotations); i++ {
		label := annotations[i].Description
		score := annotations[i].Score
		fmt.Printf("\n\tLabel: %s (%f)", label, score)

		isCat = isCat || strings.Contains(strings.ToLower(label), "cat")
	}

	if isCat {
		fmt.Printf("\n\nA cat!\n")
	} else {
		fmt.Printf("\n\nNot a cat!\n")
	}
}

type cat struct {
	name         string
	pics         []string
	adoptionLink string
}

func crawl(url string) []cat {
	// Launch a new browser with default options, and connect to it
	browser := rod.New().MustConnect()
	defer browser.Close()

	// Create a new page
	page := browser.MustPage(url).MustWaitStable()

	// Zoom out to get more cats
	page.MustEval(`() => document.body.style.zoom = "10%"`)
	page.MustWaitStable()

	cats := []cat{}

	catDivs := page.MustElements("div.baymaxGalleryCard")
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

	return cats
}

func printCat(cat cat) {
	fmt.Printf("Name: %s\n", cat.name)
	fmt.Printf("Link: %s\n", cat.adoptionLink)

	fmt.Println("Pics:")
	for _, pic := range cat.pics {
		fmt.Println(pic)
	}
}
