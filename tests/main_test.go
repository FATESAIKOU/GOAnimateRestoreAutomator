package tests

import (
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_downloader"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/notification_sender"
	"testing"
)

func TestDownloadInfoLoad(t *testing.T) {
	downloadInfo := new(magnet_link_downloader.DownloadInfo).Load("/", "error.log")
	if (downloadInfo.StoragePath != "/") || (downloadInfo.ErrorFilePath != "error.log") {
		t.Error("TestDownloadInfoLoad Failed!")
	}
}

func TestMailInfoLoad(t *testing.T) {
	mailInfo := new(notification_sender.MailInfo).LoadJson("../cfg/example_email_setting.json")
	if (mailInfo.PublisherAccount != "publisher@example.com") ||  (mailInfo.PublisherPassword != "publisher_password") ||
		(mailInfo.ImapService != "google.com:587") || (len(mailInfo.MailList) != 3) {
		t.Error("TestMailInfoLoad Failed!")
	}
}

func TestAnimateRequestInfoLoad(t *testing.T) {
	animateRequestInfo := new(magnet_link_crawler.AnimateRequestInfo).LoadJson("../cfg/example_animate_request_v2.json")
	if animateRequestInfo.AnimateStatus["animate keyword"].PreferParser[0] != "[\\d+]" {
		t.Error("TestMailInfoLoad Failed!")
	}
}
