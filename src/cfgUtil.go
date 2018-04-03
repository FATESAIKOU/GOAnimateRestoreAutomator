package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "time"
    "os"
)

type Config struct {
    Etime int64 // time.Now().Unix()
    Logs []QueryAtom
}

type QueryAtom struct {
    Team_ids []int
    Keyword string
    Episodes []float32
}


func LoadCfg(config_path string) Config {
    // Read Config contain
    raw_cfg_json, e := ioutil.ReadFile(config_path)
    if e != nil {
        fmt.Printf("File error: %v\b", e)
        os.Exit(1)
    }

    // Construct config object
    var config Config
    json.Unmarshal(raw_cfg_json, &config)

    // Return
    return config
}

func (self *Config) Save(config_path string) *Config {
    // Update Etime
    self.Etime = time.Now().Unix()

    // Marshal and Restore
    raw_cfg_json, _ := json.Marshal(self)
    ioutil.WriteFile(config_path, raw_cfg_json, 0644)

    return self
}

func (self *Config) UpdateEpisode(keyword string, value float32) *Config {
    for i := range self.Logs {
        if self.Logs[i].Keyword == keyword {
            self.Logs[i].Episodes = append(self.Logs[i].Episodes, value)
        }
    }

    return self
}
