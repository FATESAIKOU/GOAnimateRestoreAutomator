package main

import (
    "fmt"
    "os"
    "time"
    "strconv"
)

func main() {
    for {
        fmt.Printf("[%s][Start Download]\n", time.Now().Format("01/02 15:04:05"))

        // Load cfg
        config := LoadCfg(os.Args[1])

        // download
        config.StartDownload(os.Args[2], os.Args[3], os.Args[4])

        // restore result
        config.Save(os.Args[1])

        // sleep
        fmt.Printf("[%s][End of Downloading]\n", time.Now().Format("01/02 15:04:05"))
        num, _ := strconv.Atoi(os.Args[5])
        time.Sleep(time.Duration(num) * time.Second)
    }
}
