package main

import (
	"fmt"
	"magnet_link_crawler"
	"magnet_link_downloader"
	"request_parser"
)

func main() {
    fmt.Println("Hello World!")
    fmt.Println(magnet_link_crawler.CrawlMagnetLinkV0())
    fmt.Println(magnet_link_downloader.DownloadMagnetLinkV0())
	fmt.Println(request_parser.ParseRequestV0())
}
