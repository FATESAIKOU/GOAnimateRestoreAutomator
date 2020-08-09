package tests

import (
	"magnet_link_crawler"
	"magnet_link_downloader"
	"request_parser"
	"testing"
)

func TestCrawlMagnetLinkV0(t *testing.T) {
	if magnet_link_crawler.CrawlMagnetLinkV0() != "This is the response from CrawlMagnetLinkV0 for test." {
		t.Error("CrawlMagnetLinkV0 Failed!")
	}
}

func TestDownloadMagnetLinkV0(t *testing.T) {
	if magnet_link_downloader.DownloadMagnetLinkV0() != "This is the response from DownloadMagnetLinkV0 for test." {
		t.Error("DownloadMagnetLinkV0 Failed!")
	}
}

func TestParseRequestV0(t *testing.T) {
	if request_parser.ParseRequestV0() != "This is the response from ParseRequestV0 for test." {
		t.Error("ParseRequestV0 Failed!")
	}
}
