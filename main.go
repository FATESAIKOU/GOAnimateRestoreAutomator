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
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/notification_sender"
	"os"
)

func main() {
	storagePath := os.Args[1]
	errorFilePath := os.Args[2]
	emailFilePath := os.Args[3]
	animateRequestFilePath := os.Args[4]

	downloadInfo := new(magnet_link_downloader.DownloadInfo).Load(
		storagePath, errorFilePath)
	mailInfo := new(notification_sender.MailInfo).LoadJson(emailFilePath)
	animateRequestInfo := new(magnet_link_crawler.AnimateRequestInfo).LoadJson(animateRequestFilePath)

	fmt.Println(downloadInfo.ErrorFilePath)
	fmt.Println(mailInfo)
	fmt.Println(animateRequestInfo.AnimateStatus["animate keyword"].PreferParser[0])
}
