package main

import (
    "fmt"
    "os"
    "os/exec"
)

func Download(cands []Candidate, log *QueryAtom, storage_path string) {
    for i := range cands {
        fmt.Printf("Download: %s %v\n", cands[i].keyword, cands[i].episodes)
        cmd := exec.Command("webtorrent", cands[i].magnet)
        cmd.Dir = storage_path
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()

        cmd.Wait()

        (*log).Episodes = append((*log).Episodes, cands[i].episodes...)
    }
}

func (self *Config) StartDownload(storage_path string) *Config {
    var logs *[]QueryAtom
    logs = &self.Logs

    for i := range *logs {
        cands := GetCands((*logs)[i].Keyword, (*logs)[i].Team_ids, (*logs)[i].Episodes)

        Download(cands, &(*logs)[i], storage_path)
    }

    return self
}
