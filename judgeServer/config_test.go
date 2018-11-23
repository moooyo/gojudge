package main


import (
    "testing"
    "reflect"
)

func TestParseConfig(t *testing.T) {
    configCase := Config {
        ProcessConfig: ProcessServerConfig {
            Port: 8081,
        },
        ListenConfig: ListenServerConfig {
            Port: 8080,
        },
        DispatcherConfig: DispatcherConfig {
            QueueSize: 1024,
            DispatchChannelSize: 1024,
        },
    }

    configPath := "./config.json"
    config, err := ParseConfig(configPath)
    

    if err != nil {
        t.Error(err)
    }

    if !reflect.DeepEqual(configCase, config)  {
        t.Error("test case: ", configCase, "\n not equal\n", "result :", config)
    }
}
