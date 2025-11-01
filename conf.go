package main

import (
	"encoding/json"
	"errors"
	"os"
)


type Config struct {
    Ipaddr string `json:"ipaddr"`
    Port string `json:"port"`
    Process []string `json:"process"`
}

func LoadConfig(cfg *Config,cfgpath string) error {
    if cfgpath == "" {
        return errors.New("the config path is nil")
    }

    filedata,err := os.ReadFile(cfgpath)

    if err != nil {
        return err
    }

    err = json.Unmarshal(filedata,cfg)

    if err != nil {
        return err
    }

    return nil
}


