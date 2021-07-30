package scrape

import (
	"log"
	"os"
	"strconv"
	"time"

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

func (s *ScrapeBot) Run(url string) {
	count := 1
	s.collector.OnHTML("a", func(h *colly.HTMLElement) {
		questionUrl := h.Attr("href")
		s.Sqrape(url+questionUrl, &count)
	})

	if err := s.collector.Visit(url); err != nil {
		log.Printf("Error connection: %s", err.Error())
	}

}

func (s *ScrapeBot) Sqrape(url string, count *int) {
	var (
		success bool
		err     error
	)

	data := map[string]string{}
	newCount := *count + 1
	timer := time.NewTimer(time.Second * 1)

	s.collector.OnHTML("form", func(h *colly.HTMLElement) {

		h.DOM.Find("p").Each(func(i int, s *goquery.Selection) {

			item := &model.Group{}
			items := s.Find("input")
			lenItems := len(items.Nodes)

			if lenItems > 0 {
				data, err = setInputsSettings(items, item, data)
				if err != nil {
					log.Printf("Error of inputs settings: %s", err.Error())
				}
			}

			sItem := &model.Group{}
			selects := s.Find("select")
			lenSelects := len(selects.Nodes)

			if lenSelects > 0 {
				data, err = setSelectSettings(selects, sItem, data)
				if err != nil {
					log.Printf("Error of inputs settings: %s", err.Error())
				}
			}
		})

	})

	if err := s.collector.Visit(url); err != nil {
		success = true
		if success {
			log.Printf("Test successfully passed\n")
			os.Exit(0)
		}
	}

	if err := s.collector.Post(url, data); err != nil {
		log.Printf("Error post: %s\n", err.Error())
	}

	log.Printf("Question %d passed: %v\n", *count, data)

	<-timer.C

	s.Sqrape("http://test.youplace.net/question/"+strconv.Itoa(newCount), &newCount)

}

func searchLongValue(items []string) string {
	long := items[0]
	for i := 0; i < len(items); i++ {

		if len([]rune(long)) <= len([]rune(items[i])) {
			long = items[i]
		}
	}

	return long
}

func setInputsSettings(items *goquery.Selection, item *model.Group, data map[string]string) (map[string]string, error) {
	inputType := ""
	items.Each(func(i int, s *goquery.Selection) {
		typeInput, ok := s.Attr("type")

		if !ok {
			log.Println("type not found")
		}

		switch typeInput {
		case "radio":
			nameInput, ok := s.Attr("name")
			if !ok {
				log.Println("value not found")
			}
			inputType = typeInput
			item.Name = nameInput

			val, ok := s.Attr("value")

			if !ok {
				log.Println("value not found")
			}

			item.Items = append(item.Items, val)
		case "text":
			nameInput, ok := s.Attr("name")
			if !ok {
				log.Println("value not found")
			}
			inputType = typeInput
			item.Name = nameInput
			s.SetAttr("value", "test")
			data[nameInput] = "test"
		default:
		}
	})

	if inputType == "radio" {

		hightItem := searchLongValue(item.Items)

		items.Each(func(i int, s *goquery.Selection) {
			val, ok := s.Attr("value")

			if !ok {
				log.Printf("Value not found\n")
			}

			name, ok := s.Attr("name")

			if !ok {
				log.Printf("Name not found\n")
			}

			if val == hightItem {
				s.SetAttr("checked", "true")
				data[name] = val
			}
		})

	}

	return data, nil
}

func setSelectSettings(items *goquery.Selection, item *model.Group, data map[string]string) (map[string]string, error) {
	items.Each(func(i int, s *goquery.Selection) {
		selectName, ok := s.Attr("name")
		if !ok {
			log.Printf("Name not found\n")
		}
		item.Name = selectName

		s.Each(func(i int, s *goquery.Selection) {
			options := s.Find("option")

			options.Each(func(i int, s *goquery.Selection) {
				val, ok := s.Attr("value")

				if !ok {
					log.Printf("Value not found\n")
				}

				item.Items = append(item.Items, val)
			})
		})
	})

	hightItem := searchLongValue(item.Items)

	items.Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")

		if !ok {
			log.Printf("Name not found\n")
		}
		s.Each(func(i int, s *goquery.Selection) {
			options := s.Find("option")

			options.Each(func(i int, s *goquery.Selection) {
				val, ok := s.Attr("value")

				if !ok {
					log.Printf("Value not found\n")
				}

				if val == hightItem {
					s.SetAttr("selected", "true")
					data[name] = val
				}
			})
		})

	})

	return data, nil
}
