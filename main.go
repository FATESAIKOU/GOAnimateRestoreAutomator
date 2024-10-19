/*************************
@argv[1]: Storage path
@argv[2]: Error file path
@argv[3]: Email setting file path
@argv[4]: Animate request file path
 *************************/
package main

import (
	"fmt"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_downloader"
	"os"
)

func main() {
	storagePath := os.Args[1]
	errorFilePath := os.Args[2]
	emailFilePath := os.Args[3]
	animateRequestFilePath := os.Args[4]

	// Load config files
	downloadInfo := new(magnet_link_downloader.DownloadInfo).Load(
		storagePath, errorFilePath)
	mailInfo := new(notification_sender.MailInfo).LoadJson(emailFilePath)
	animateRequestInfo := new(magnet_link_crawler.AnimateRequestInfo).LoadJson(animateRequestFilePath)

	// Output checking info
	fmt.Println("[Storage Path]", downloadInfo.StoragePath)
	fmt.Println("[Error Log Path]", downloadInfo.ErrorFilePath)
	fmt.Println("[Publisher]", mailInfo.PublisherAccount)
	fmt.Println("[Mail List]", mailInfo.MailList)
	for animateKeyword, _ := range animateRequestInfo.AnimateStatusMap {
		fmt.Println("[Download Target] ", animateKeyword)
	}

	// Crawl website
	animateMagnetInfo := magnet_link_crawler.GetAnimateMagnetInfo(
		"https://share.dmhy.org/topics/list", animateRequestInfo)

	// Download
	//newDownloads := magnet_link_downloader.DownloadMagnet(animateMagnetInfo, *downloadInfo, *animateRequestInfo)
	magnet_link_downloader.DownloadMagnet(animateMagnetInfo, *downloadInfo, *animateRequestInfo)

	// Notificate
	//notification_sender.SendMail(newDownloads, *mailInfo)

	// Writeback
	animateRequestInfo.SaveJson(animateRequestFilePath)
}
