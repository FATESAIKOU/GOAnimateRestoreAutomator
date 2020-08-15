package notification_sender

import (
	"fmt"
	"github.com/FATESAIKOU/GOAnimateRestoreAutomator/magnet_link_crawler"
	"log"
	"net/smtp"
	"strings"
)

func SendMail(newDownloadeds map[string]*magnet_link_crawler.AnimateStatus, mailInfo MailInfo) {
	log.Println("[Start send notification]")

	auth := smtp.PlainAuth(
		"", mailInfo.PublisherAccount, mailInfo.PublisherPassword, mailInfo.SmtpDomain)

	content := genContent(newDownloadeds)
	if len(content) == 0 {
		return
	}

	from := mailInfo.PublisherAccount
	to   := mailInfo.MailList
	msg	 := fmt.Sprintf("From: %s\nTo: %s\nSubject: Download Log\n\n%s",
		from, strings.Join(to, ","), content)

	err := smtp.SendMail(
		mailInfo.SmtpServiceUrl,
		auth,
		from,
		to,
		[]byte(msg))

	if err != nil {
		log.Fatal("Fail to send mail", err)
	}

	log.Println("[End sned notification]")
}

func genContent(newDownloaded map[string]*magnet_link_crawler.AnimateStatus) string {
	content := ""
	for animeteKey, animateStatus := range newDownloaded {
		for _, episode := range animateStatus.CompletedEpisodes {
			content += fmt.Sprintf("[Complete][%s]: %f\n",animeteKey, episode)
		}
	}

	return content
}