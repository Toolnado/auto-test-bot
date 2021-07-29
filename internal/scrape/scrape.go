package scrape

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/Toolnado/autotest/internal/model"
	"github.com/gocolly/colly"
)

type ScrapeBot struct {
	collector *colly.Collector
}

func NewScrapeBot() *ScrapeBot {
	c := colly.NewCollector()
	return &ScrapeBot{
		collector: c,
	}
}

func (s *ScrapeBot) Scrape(url string) {
	s.collector.OnHTML("a", func(h *colly.HTMLElement) {
		questionUrl := h.Attr("href")

		clone := s.collector.Clone()

		clone.OnHTML("form", func(h *colly.HTMLElement) {
			h.DOM.Find("p").Each(func(i int, s *goquery.Selection) {

				item := &model.Group{}

				items := s.Find("input")
				lenItems := len(items.Nodes)
				if lenItems > 0 {
					items.Each(func(i int, s *goquery.Selection) {
						typeInput, ok := s.Attr("type")

						if !ok {
							log.Println("type not found")
							return
						}

						switch typeInput {
						case "radio":
							item.Name = typeInput
							val, ok := s.Attr("value")

							if !ok {
								log.Println("value not found")
								return
							}

							item.Items = append(item.Items, val)
						case "text":
							item.Name = typeInput
							s.SetAttr("value", "test")
							log.Printf("%v set to 'test'\n", item.Name)
						default:
						}
					})

					if item.Name == "radio" {
						fmt.Printf("Group: %v\n", item)

						hightItem := searchLongInput(item.Items)

						items.Each(func(i int, s *goquery.Selection) {
							val, ok := s.Attr("value")

							if !ok {
								log.Printf("Value not found\n")
							}

							if val == hightItem {
								s.SetAttr("checked", "true")
								log.Printf("%s checked set to true", val)
							}
						})

					}

				}

			})

		})
		if err := clone.Visit(url + questionUrl); err != nil {
			log.Printf("Error connection: %s", err.Error())
		}

	})

	if err := s.collector.Visit(url); err != nil {
		log.Printf("Error connection: %s", err.Error())
	}

}

func searchLongInput(items []string) string {
	var long string
	for i := 0; i < len(items); i++ {

		long = items[0]

		if len(long) < len(items[i]) {
			long = items[i]
		}
	}

	return long
}
