package main

import (
    "fmt"
    //"io/ioutil"
    //"os"
    //"os.exec"
)

func (self *Config) StartDownload() *Config {
    logs := self.Logs

    fmt.Println(logs)
    for i := range logs {
        cands := GetCands(logs[i].Keyword, logs[i].Team_ids, logs[i].Episodes)
        fmt.Println(cands)

        // download(rows)
    }

    return self
}
