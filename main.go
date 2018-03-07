package main
// This is vaper agent.
import (
    "flag"
    "net"
    "strconv"
    "os" 
    // "os/exec"
    // "bytes"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "regexp"
    // "reflect"
    log "github.com/sirupsen/logrus"
    "github.com/satori/go.uuid"
    "github.com/robfig/cron"
    "github.com/google/gopacket"
    // "github.com/google/gopacket/pcap"
)

// shell command 命令行参数
type DefaultFlags struct{
    configFilePath string
    action string
}

func initDefaultFlags() DefaultFlags{
    defaultFlags :=DefaultFlags{}
    // config file location
    flag.StringVar(&defaultFlags.configFilePath,"f", "./conf/config.default.ini", "The path of config file.")
    // action
    flag.StringVar(&defaultFlags.action,"a", "nothing", "action: init/start")
    
    flag.Parse()
    return defaultFlags
}

// get all net interfaces 所有网卡
func getAllInterfaces() []net.Interface {
    //Get all interface info
    interfaces,err := net.Interfaces()
    if err != nil{
        log.Fatal("Failed to get the interfaces info.")
    }
    for k, interf := range interfaces {
        log.Debug("Interface " +strconv.Itoa(k)+":"+interf.Name)
    }
    return interfaces
}

// Net relation 网络流量关系
type NetRelation struct{
    SendIp string
    SendPort string
    ReceiverIp string
    RecevierPort string
}

// func getNetRelations() []NetRelation {
//     // Run tcpdump with parameters shell command
//     command := exec.Command("tcpdump", "-i", "any", "-c","10","-f","-q","-nn")
//     stdOutBuff := bytes.NewBuffer(nil)
//     stdErrBuff := bytes.NewBuffer(nil)
//     command.Stderr = stdErrBuff
//     command.Stdout = stdOutBuff
//     if err := command.Run(); err != nil {
//         log.Error("The tcpdump command can not exec.Run returns: %s\n", err)
//     }
//     CommandOutPut := string(stdOutBuff.Bytes())
//     // log.Debug("tcpdump Stdout:" + CommandOutPut)
//     reg := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\.(\d+)\s*\>\s*(\d+\.\d+\.\d+\.\d+)\.(\d+)`)
//     relations := reg.FindAllString(CommandOutPut, -1)
//     var netRalations []NetRelation
//     for i := 0; i < len(relations); i++ {
//         relaStr := relations[i]
//         netRalation := string2NetRelation(relaStr)
//         netRalations = append(netRalations, netRalation)
//         // log.Debug(netRalation.SendIp+":"+netRalation.sendPort +"->-"+netRalation.receiverIp+":"+netRalation.recevierPort)
//     }
//     return netRalations
// }

func string2NetRelation(relaStr string) NetRelation {
    reg := regexp.MustCompile(`(?P<ip1>\d+\.\d+\.\d+\.\d+)\.(?P<port1>\d+)\s*\>\s*(?P<ip2>\d+\.\d+\.\d+\.\d+)\.(?P<port2>\d+)`)
    match := reg.FindStringSubmatch(relaStr)
    netRelation := NetRelation{}
    netRelation.SendIp = match[1]
    netRelation.SendPort = match[2]
    netRelation.ReceiverIp = match[3]
    netRelation.RecevierPort = match[4]
    return netRelation
}

// get host name 获取主机名称
func getHostname() string{
    host, err := os.Hostname()
    if err != nil {
        log.Error("Get hostname:" + err.Error())
        return "error"
    } else {
        return host
    }
}

// get ip list
func get_internal_ips() []string{
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Error("ipnet ip:" + err.Error())
    }
    var ips []string
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return ips
}

// host info
type Host struct{
    Hostname string
    Uid string // unique id for this host
    Ips []string
}

func getHostMeta(config * Config) Host{
    var host Host
    host.Hostname = getHostname()
    host.Ips = get_internal_ips()
    host.Uid = getUuid(config.Ini.GetValue("basic", "uuid_path"))
    return host
}



func sendHost(url string, host Host) bool {
    host_bt,_ := json.Marshal(host)
    host_str := string(host_bt)
    payload := strings.NewReader(host_str)
    req, _ := http.NewRequest("POST", url, payload)
    req.Header.Add("content-type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if(err != nil){
        log.Error("send host info fail.detail:"+err.Error())
        return false
    }
    defer res.Body.Close()
    // body, _ := ioutil.ReadAll(res.Body)
    return true
}

// func sendVpMsg(url string, vpMsg VpMsg) bool {
//     vpMsg_bt,_ := json.Marshal(vpMsg)
//     vpMsg_str := string(vpMsg_bt)
//     payload := strings.NewReader(vpMsg_str)

//     req, _ := http.NewRequest("POST", url, payload)
//     req.Header.Add("content-type", "application/json")

//     res, err := http.DefaultClient.Do(req)
//     if(err != nil){
//         log.Error("send msg info fail.detail:"+err.Error())
//     }
//     defer res.Body.Close()
//     // body, _ := ioutil.ReadAll(res.Body)
//     // log.Debug(res)
//     // log.Debug(body)
//     return true
// }

//generate a new uid
func (this *Actions) generateUid(hostname string) string{
    uid := uuid.NewV5(uuid.NewV1(), hostname).String()
    return uid
}

type Actions struct{

}
func check(e error) {
 if e != nil {
  panic(e)
 }
}

func getUuid(filePath string)string{
    uuid_file, error := os.Open(filePath)
    if error != nil {
        log.Error("Read uuid file from '"+ filePath +"' failed." +  error.Error() +
        ".Please run command './vaper_agent -a init ' for generate or refresh the uuid first.")
        return ""
    }
    defer uuid_file.Close()
    uuid_bt,err := ioutil.ReadAll(uuid_file)
    if err != nil {
        log.Error("Read uudi string from '"+ filePath +"' failed." +  error.Error()+
        ".Please run command './vaper_agent -a init ' for generate or refresh the uuid first.")
        return ""
    }
    uid := string(uuid_bt)
    uuid_file.Close()
    return uid
}

// init
func (this *Actions) Init( config *Config) bool{
    hostname := getHostname()
    uuidNew := this.generateUid(hostname)
    uuid_path := config.Ini.GetValue("basic", "uuid_path")
    log.Info("The new uuid is : " + uuidNew)
    file, error := os.OpenFile(uuid_path, os.O_RDWR|os.O_CREATE, 0622)
    if error != nil {
        log.Error("init uuid file in "+ uuid_path +"failed." +  error.Error())
    }
    file.WriteString(uuidNew)
    file.Close()
    return true
}


/*
HostJob
*/
type HostJob struct{
    config *Config
}
 
func NewHostJob (config *Config) HostJob{
    hostJob := HostJob{}
    hostJob.config = config
    return hostJob
}

func (this HostJob)Run(){
    host := getHostMeta(this.config)
    server_http_url := this.config.Ini.GetValue("server","http_url")
    host_url := server_http_url + this.config.Ini.GetValue("url","host")
    sendHost(host_url, host)
}


/*
NetRelationShipJob
*/
// type NetRelationShipJob struct{
//     config *Config
// } 

// func NewNetRelationShipJob (config *Config) NetRelationShipJob{
//     netRelationShipJob := NetRelationShipJob{}
//     netRelationShipJob.config = config
//     return netRelationShipJob
// }

// func (this NetRelationShipJob)Run(){
//     var vpMsg VpMsg = getVpMsg(this.config)
//     server_http_url := this.config.Ini.GetValue("server","http_url")
//     netrelation_url := server_http_url + this.config.Ini.GetValue("url","netrelation")
//     sendVpMsg(netrelation_url, vpMsg)
// }



func postMsg(url string, msg string) bool {
    println("222222222222")

    payload := strings.NewReader(msg)
    req, _ := http.NewRequest("POST", url, payload)
    req.Header.Add("content-type", "application/json")
    println("33333333333333")
    res, err := http.DefaultClient.Do(req)
    if(err != nil){
        log.Error("send msg info fail.detail:"+err.Error())
    }
    println("44444444444")
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    println(res)
    println(body)
    println("55555555555555555")
    return true
}


/*
NetworkflowsJob
*/
type Flow struct{
    SrcIp string
    DstIp string
}

type NetworkFlowMsg struct{
    Uid string
    NetworkFlows []Flow
}

func getNetworkFlowMsg(config *Config) NetworkFlowMsg{
    var(
        networkFlows []gopacket.Flow = tcpcatch(10, 10)
        networkFlowMsg NetworkFlowMsg
        flows []Flow
    )
    for _,netflow := range networkFlows{
        flow := Flow{}
        flow.SrcIp = netflow.Src().String()
        flow.DstIp = netflow.Dst().String()
        flows = append(flows, flow)
    }
    networkFlowMsg.Uid = getUuid(config.Ini.GetValue("basic", "uuid_path"))
    networkFlowMsg.NetworkFlows = flows
    return networkFlowMsg
}

type NetworkflowsJob struct{
    config *Config
} 

func NewNetworkflowsJob (config *Config) NetworkflowsJob{
    networkflowsJob := NetworkflowsJob{}
    networkflowsJob.config = config
    return networkflowsJob
}

func (this NetworkflowsJob)Run(){
    nwfMsg := getNetworkFlowMsg(this.config)
    if nwfMsg.NetworkFlows == nil{
        log.Debug("No flow data was catched.")
        return
    }else{
        server_url := this.config.Ini.GetValue("server","server_url")
        netrelation_url := server_url + this.config.Ini.GetValue("url","netrelation")
        msg_bt,_ := json.Marshal(nwfMsg)
        msg_str := string(msg_bt)
        println(msg_str)
        postMsg(netrelation_url, msg_str)
        println("66666666666666666666")
    }
}


//start the agent
func (this *Actions) Start(config *Config){
    println("start")
    c := cron.New()

    //host meta job
    hostJob := NewHostJob(config)
    hostInfoFrequency := config.Ini.GetValue("server", "hostInfoFrequency")
    log.Info("The hostInfoFrequency is :"+ hostInfoFrequency)
    c.AddJob("@every "+ hostInfoFrequency +"s", hostJob)

    //netRelationShipJob
    // netRelationShipJob := NewNetRelationShipJob(config)
    // netRelationshipFrequency := config.Ini.GetValue("server", "netRelationshipFrequency")
    // log.Info("The netRelationshipFrequency is :"+ netRelationshipFrequency)
    // c.AddJob("@every "+ netRelationshipFrequency +"s", netRelationShipJob)

    //network flows job
    networkflowsJob := NewNetworkflowsJob(config)
    networkFlowFrequency :=  config.Ini.GetValue("server", "networkFlowFrequency")
    log.Info("The networkFlowFrequency is :"+ networkFlowFrequency)
    c.AddJob("@every "+ networkFlowFrequency +"s", networkflowsJob)

    c.Start()
    
    server_http_url := config.Ini.GetValue("server","http_url")
    log.Info("The Server url is : " + server_http_url)
    println("Running......")

    select{}
}


func main() {
    defaultFlags := initDefaultFlags()
    actionName := defaultFlags.action
    config := NewConfig(defaultFlags.configFilePath)
    
    version := config.Ini.GetValue("basic","version")
    log.Info("Vaper-"+ version + " " +actionName)
    
    actions := Actions{}
    switch actionName {
    case "init":
        actions.Init(config)
    case "start":
        actions.Start(config)
    case "nothing" :
        println("nothing to do, use -h for help.")
    default:
        println("nothing to do, use -h for help.")
    }
    // log.Debug(host.ips)
    // log.Debug("hostname:"+host.hostname)
    // log.Debug(host.uid)
    // log.Debug(vpMsg.NetRelations)
}