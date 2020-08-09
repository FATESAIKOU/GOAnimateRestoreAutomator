package main

import (
    "log"
    "fmt"
    "os"
    "time"
    "bytes"
	"context"
    "os/exec"
    "net/smtp"
)

func Download(cands []Candidate, log *QueryAtom, storage_path string) string {
    var download_log bytes.Buffer
    for i := range cands {
        msg := fmt.Sprintf("Download: %s %v", cands[i].keyword, cands[i].episodes)
        download_log.WriteString(msg + "\n")

        for {
            fmt.Println(msg)
            err := DownloadCand(cands[i].magnet, storage_path)
            
            if err == nil {
                break
            }
            
            fmt.Println("Restart")
        }

        (*log).Episodes = append((*log).Episodes, cands[i].episodes...)
    }

    return download_log.String()
}

func DownloadCand(magnet_link string, storage_path string) error {
    ctxt, cancel := context.WithTimeout(context.Background(), 10800*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctxt, "webtorrent", magnet_link)
    cmd.Dir = storage_path
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        fmt.Println("Run Error")
        if ctxt.Err() == context.DeadlineExceeded {
            fmt.Println("Command time out")
            return err
        }
    }
    
    cmd.Wait()
    
    return nil
}

func sendMail(content string, user string, passwd string) {
    fmt.Printf("[%s][Start SendMail]\n", time.Now().Format("01/02 15:04:05"))
    auth := smtp.PlainAuth(
        "",
        user,
        passwd,
        "mail.ncku.edu.tw",
    )

    from := "mail service account"
    to   := "to email"
    to2  := "to email2"
    msg  := "From: " + from + "\n" +
            "To: " + to + "\n" +
            "Subject: Download Log\n\n" +
            content

    err := smtp.SendMail(
        "mail service url",
        auth,
        from,
        []string{to, to2},
        []byte(msg),
    )

    fmt.Printf("[%s][End SendMail]\n", time.Now().Format("01/02 15:04:05"))
    if err != nil {
        log.Fatal(err)
    }
}

func (self *Config) StartDownload(storage_path string, user string, passwd string) *Config {
    var logs *[]QueryAtom
    logs = &self.Logs

    var download_log bytes.Buffer
    for i := range *logs {
        cands := GetCands((*logs)[i].Keyword, (*logs)[i].Team_ids, (*logs)[i].Episodes)

        download_log.WriteString(Download(cands, &(*logs)[i], storage_path))
    }

    content := download_log.String()
    if len(content) > 0 {
        sendMail(content, user, passwd)
    }

    return self
}
