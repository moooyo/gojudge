package main

import (
    "flag"
    "strconv"
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
    var listenServer ListenServer
    addr := "127.0.0.1:" + strconv.Itoa(config.ListenConfig.Port)
    RunServer(&listenServer, addr)
}
