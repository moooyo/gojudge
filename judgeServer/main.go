package main

import (
    "flag"
    "fmt"
)


var configPath string

func init() {
    flag.StringVar(&configPath, "configfile", "./config.json", "config file of judgeServer")
}

func main() {
    flag.Parse()

    config, err := ParseConfig(configPath)

    if err != nil {
        panic(err)
    }
    
    fmt.Println(config)
}
