package main

func getDefaultSetting(key1 string,key2 string) string{

    var DEFAULT_SETTINGS map[string]string = make(map[string]string)
    //basic
    DEFAULT_SETTINGS["basic.version"] = "0.0.1"
    DEFAULT_SETTINGS["basic.uuid_path"] = "./vaper_agent.uid"
    DEFAULT_SETTINGS["basic.auto_generate_uid"] = "1"

    //log
    DEFAULT_SETTINGS["log.level"] = "error"
    DEFAULT_SETTINGS["log.path"] = "./vaper_agent.log"

    //server
    DEFAULT_SETTINGS["server.server_url"] = "http://127.0.0.1:3000"

    //performance
    DEFAULT_SETTINGS["performance.hostInfoFrequency"] = "60"
    DEFAULT_SETTINGS["performance.networkFlowFrequency"] = "10"
    DEFAULT_SETTINGS["performance.packages_limit"] = "10"


    //api
    DEFAULT_SETTINGS["api.host_add_or_update"] = "/host/add_or_update"
    DEFAULT_SETTINGS["api.netflow_add"] = "/netflow/add"

    var value string = DEFAULT_SETTINGS[key1+"."+key2]
    return value
}

