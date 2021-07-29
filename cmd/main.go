package main

import "github.com/Toolnado/autotest/internal/scrape"

func main() {

	url := "http://test.youplace.net"

	bot := scrape.NewScrapeBot()

	bot.Scrape(url)
}
