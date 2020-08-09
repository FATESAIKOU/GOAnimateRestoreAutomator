package main

import (
	"fmt"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_downloader"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/request_parser"
)

func main() {
    fmt.Println("Hello World!")
    fmt.Println(magnet_link_crawler.CrawlMagnetLinkV0())
    fmt.Println(magnet_link_downloader.DownloadMagnetLinkV0())
	fmt.Println(request_parser.ParseRequestV0())
}
