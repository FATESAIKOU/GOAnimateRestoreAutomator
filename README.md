# GOAnimateRestoreAutomator
An automized btorrent downloader for downloading animate from https://dmhy.org 

## Config File

### AnimateRequestConfig
This config file will contain the information about what animate you want to subscribe, and what interpretation team you prefer(as id), the prefered pattern to extract episode from the title.

example:
```json
{
  "租借女友":
  {
    "CompletedEpisodes": [
      1,
      2,
      3,
      4,
      6,
      5
    ],
    "PreferTeamIds": [],
    "PreferParser": ""
  }
}
```
### MailConfig
GOAnimateRestoreAutomator also support to send notification after every time animate been downloaded newly.
Currently, the notificator only support using smtp mail service to do this task. 

example:
```json
{
  "PublisherAccount": "publisher@example.com",
  "PublisherPassword": "publisher_password",
  "SmtpDomain":  "mail.example.com",
  "SmtpServiceUrl": "mail.exmaple.com:587",
  "MailList": [
    "subscriber1@example.com",
    "subscriber2@example.com",
    "subscriber3@example.com"
  ]
}
```

## To execute

### Prerequisites
* Follow the instructions at [golang official document](https://golang.org/doc/) to install golang.
* Created your own `AnimateRequestConfig` and `MailConfig`.

### Execution
```shell script
$ go run github.com/FATESAIKOU/GOAnimateRestoreAutomator \ 
[the path to restore the animate] \
[the path to put the error log  (but currently not used)] \
[your AnimateRequestConfig file] \
[your MailConfig file]
```