package main
// This is vaper agent.
import (
    "flag"
    "os"
    "encoding/json"
    log "github.com/sirupsen/logrus"
    "github.com/robfig/cron"
    "github.com/satori/go.uuid"

)

// shell command 命令行参数
type DefaultFlags struct{
    configFilePath string
    action string
}

func initDefaultFlags() DefaultFlags{
    
    defaultFlags :=DefaultFlags{}
    // config file location
    flag.StringVar(&defaultFlags.configFilePath,"f", "./vaper_agent.ini", "The path of config file.")
    // action
    flag.StringVar(&defaultFlags.action,"a", "nothing", "action: init/start")
    
    flag.Parse()
    return defaultFlags
}

type Actions struct{

}
func check(e error) {
    if e != nil {
        panic(e)
    }
}

//generate a new uid
func (this *Actions) generateUid(config *Config, hostname string) string{
    uuidNew := uuid.NewV5(uuid.NewV1(), hostname).String()
    uuid_path := config.GetValue("basic", "uuid_path")
    log.Info("The new uuid is : " + uuidNew)
    file, error := os.OpenFile(uuid_path, os.O_RDWR|os.O_CREATE, 0622)
    if error != nil {
        log.Error("Open uuid file in "+ uuid_path +" failed." +  error.Error())
    }
    _,err := file.WriteString(uuidNew)
    if err != nil {
        log.Error("Save uuid file in "+ uuid_path +" failed." +  err.Error())
    }
    file.Close()
    return uuidNew
}

// init Config
func (this *Actions) Init( config *Config) bool{
    hostname := getHostname()
    this.generateUid(config, hostname)
    return true
}


//start the agent
func (this *Actions) Start(config *Config){
    uuid_exist := getUuid(config.GetValue("basic", "uuid_path"))
    if( uuid_exist == "" && config.GetValue("basic", "auto_generate_uid") == "1"){
        this.Init(config)
    }
    println("start")
    c := cron.New()

    //host meta job
    hostJob := NewHostJob(config)
    hostInfoFrequency := config.GetValue("performance", "hostInfoFrequency")
    log.Info("The hostInfoFrequency is :"+ hostInfoFrequency)
    c.AddJob("@every "+ hostInfoFrequency +"s", hostJob)
    host_bt,_ := json.Marshal(getHostMeta(config))
    log.Info("The Host Meta :"+ string(host_bt))

    //network flows job
    networkflowsJob := NewNetworkflowsJob(config)
    networkFlowFrequency :=  config.GetValue("performance", "networkFlowFrequency")
    log.Info("The networkFlowFrequency is :"+ networkFlowFrequency)
    c.AddJob("@every "+ networkFlowFrequency +"s", networkflowsJob)
    log.Info("The Interfaces list :"+ InterfacesToString(getAllInterfaces()))
    
    c.Start()
    
    server_url := config.GetValue("server","server_url")
    log.Info("The server url is : " + server_url)
    println("Running......")

    select{}
}

func checkConfigFile(filepath string)bool{
    file, error := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0622)
    if error != nil {
        log.Panic("Open config file in "+ filepath +" failed." +  error.Error())
        return false
    }else{
        file.Close()
        return true
    }
}

func main() {
    defaultFlags := initDefaultFlags()
    actionName := defaultFlags.action
    checkConfigFile(defaultFlags.configFilePath)
    config := NewConfig(defaultFlags.configFilePath)

    
    version := config.GetValue("basic","version")
    log.Info("VaperAgent - v"+ version + " " +actionName)
    
    actions := Actions{}
    switch actionName {
    case "init":
        actions.Init(config)
    case "start":
        actions.Start(config)
    case "nothing" :
        println("VaperAgent nothing to do, use -h for help.")
    default:
        println("VaperAgent nothing to do, use -h for help.")
    }
}