package main
// This is vaper agent.
import (
    "net"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

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
    host.Uid = getUuid(config.GetValue("basic", "uuid_path"))
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
    body, _ := ioutil.ReadAll(res.Body)
    log.Debug(url + ";request:"+ host_str +";response:" + string(body))
    return true
}



func getUuid(filePath string)string{
    uuid_file, error := os.Open(filePath)
    if error != nil {
        log.Warn("Open uuid file from '"+ filePath +"' failed." +  error.Error() +
        ".Please run command './vaper_agent -a init ' for generate or refresh the uuid first.")
        return ""
    }
    defer uuid_file.Close()
    uuid_bt,err := ioutil.ReadAll(uuid_file)
    if err != nil {
        log.Warn("Read uudi string from '"+ filePath +"' failed." +  error.Error()+
        ".Please run command './vaper_agent -a init ' for generate or refresh the uuid first.")
        return ""
    }
    uid := string(uuid_bt)
    uuid_file.Close()
    return uid
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
    server_url := this.config.GetValue("server","server_url")
    host_url := server_url + this.config.GetValue("api","host_add_or_update")
    sendHost(host_url, host)
}
