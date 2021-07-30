package main

import "github.com/Toolnado/autotest/internal/scrape"

const (
	url = "http://test.youplace.net"
)

func main() {

	bot := scrape.NewScrapeBot()

	bot.Run(url)
}
