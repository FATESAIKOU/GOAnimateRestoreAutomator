package main

import (
    "fmt"
    "os"
    "time"
    "strconv"
)

func main() {
    for {
        fmt.Println("[Start Download]")

        // Load cfg
        config := LoadCfg(os.Args[1])

        // download
        config.StartDownload(os.Args[2])

        // restore result
        config.Save(os.Args[1])

        // sleep
        fmt.Println("[End of Downloading]")
        num, _ := strconv.Atoi(os.Args[3])
        time.Sleep(time.Duration(num) * time.Second)

    }
}
