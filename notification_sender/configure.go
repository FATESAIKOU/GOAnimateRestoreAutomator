package notification_sender

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// The settings for mailing to the subscriber
type MailInfo struct {
	PublisherAccount string
	PublisherPassword string
	SmtpDomain string
	SmtpServiceUrl string
	MailList []string
}

func (mailInfoSelf *MailInfo) LoadJson(jsonFilePath string) *MailInfo {
	rawJson, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatal("Fail to read email config file: ", err)
	}

	err = json.Unmarshal(rawJson, mailInfoSelf)

	if err != nil {
		log.Fatal("Fail to parse email config file: ", err)
	}

	return mailInfoSelf
}
