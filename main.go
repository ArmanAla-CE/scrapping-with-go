package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const PAGE_COUNT int = 100

func main() {
	var accademy string
	fmt.Scanf("%s", &accademy)

	titleIndex := 0
	pageIndex := 0
	// Create a new collector
	c := colly.NewCollector()

	// Open a file for writing titles
	file, err := os.Create("titles.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Callback to be executed for each matched element
	c.OnHTML("h2 > a", func(e *colly.HTMLElement) {
		title := e.Text
		link := e.Attr("href")
		// fmt.Println(title, link)
		if title[0:len(accademy)] == accademy {
			titleIndex += 1

			// Write the title to the file
			_, err := file.WriteString(
				fmt.Sprintf(
					"%d- Course title: %s\nLink: %s\n--------------------------------------------------\n",
					titleIndex,
					title,
					link,
				),
			)
			if err != nil {
				log.Println("Error writing to file:", err)
			}
		}
	})

	// Callback to find and follow the next page link
	c.OnHTML("nav > div > a.next.page-numbers", func(e *colly.HTMLElement) {
		pageIndex += 1
		if pageIndex == PAGE_COUNT {
			log.Fatal("Maximum number of pages is reached")
		}
		fmt.Println("Reached to page: ", pageIndex)
		nextPageLink := e.Attr("href")
		if strings.HasPrefix(nextPageLink, "/") {
			// Concatenate the base URL if the next page link is relative
			nextPageLink = "https://downloadly.ir/download/elearning/video-tutorials" + nextPageLink
		}

		// Visit the next page
		if err := c.Visit(nextPageLink); err != nil {
			log.Fatal(err)
		}
	})

	// Start scraping from the initial page
	if err := c.Visit("https://downloadly.ir/download/elearning/video-tutorials/page/1/"); err != nil {
		log.Fatal(err)
	}
}
