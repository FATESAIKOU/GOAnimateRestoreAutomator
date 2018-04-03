package main

import (
    //"fmt"
    "os"
)

func main() {
    // Load cfg
    config := LoadCfg(os.Args[1])

    // download
    config.StartDownload()

    // restore result
    //config.Save(os.Args[1])
}
