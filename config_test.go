package main

import (
    "testing"
    "regexp"
)

func Test_NewConfig_1(t *testing.T) {
    config := NewConfig("./conf/config.ini")
    version := config.Ini.GetValue("basic","version")
    isMatch, _:= regexp.MatchString(`^(\d+)\.(\d+)\.(\d+)$`, version)
    if(isMatch == false){
        t.Error("error")
    }else{
        t.Log("PASS")
    }
}
