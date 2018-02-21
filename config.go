package main
/* This is a ini file configure read function */
import (
    "os"
    "strings"
    "github.com/widuu/goini"
    log "github.com/sirupsen/logrus"
)

type Config struct {
    version string
    Ini *goini.Config
    defaultIni *goini.Config // Someday maybe implement.
}

func NewConfig (conf_path string) *Config{
    config := &Config{}
    config.Ini = goini.SetConfig(conf_path)
    config.initLogger()
    return config
}

func (config *Config) initLogger(){
    level := strings.ToLower(config.Ini.GetValue("log","level"))
    logpath := config.Ini.GetValue("log","path")
    file, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY, 0666)
    if err == nil {
        log.SetOutput(file)
    } else {
        log.Error("Failed to log to file:"+logpath)
    }
    switch {
        case level == "debug":
            log.SetLevel(log.DebugLevel)
        case level == "info":
            log.SetLevel(log.InfoLevel)
        case level == "warn":
            log.SetLevel(log.WarnLevel)
        case level == "error":
            log.SetLevel(log.ErrorLevel)
        case level == "fatal":
            log.SetLevel(log.FatalLevel)
        case level == "panic":
            log.SetLevel(log.PanicLevel)
        default:
            log.SetLevel(log.ErrorLevel)
    }
    log.Info("Log output to :"+logpath)
    log.Info("Log Level is :"+level)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}