package main

import (
    "strconv"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "time"
    "io"
    logrus "github.com/sirupsen/logrus"
)


func getPkgsByDeviceName(deviceName string, networkflows chan []gopacket.Flow, limit int, timeoutSecond int ){
    var (
        snapshot_len int32  = 1024
        promiscuous  bool   = false
        err          error
        timeout      time.Duration = time.Duration(timeoutSecond) * time.Second
        handle       *pcap.Handle
        nwfs [] gopacket.Flow
    )
    handle, err = pcap.OpenLive(deviceName, snapshot_len, promiscuous, timeout)
    if err == nil {
        defer handle.Close()
        count := 0
        packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
        for {
            packet, err := packetSource.NextPacket()
            if err == io.EOF {
                break
            } else if err != nil {
                logrus.Debug("No more packages." + err.Error())
                break
            }
            applicationLayer := packet.NetworkLayer()
            if applicationLayer != nil{
                flow := applicationLayer.NetworkFlow()
                nwfs = append(nwfs, flow)

                // fmt.Println(flow.Src())
                // fmt.Println(flow.Dst())
            }
            count += 1
            if count > limit{
                logrus.Debug("get "+ strconv.Itoa(count) +" package.")
                break
            }
        }
    }
    networkflows <- nwfs
}

func tcpcatch (limit int,timeoutSecond int) [] gopacket.Flow{
    var(
        // timeoutSecond int = 10
        // limit int = 10
        networkflowsCh chan []gopacket.Flow = make(chan []gopacket.Flow)
    )
    devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatal(err)
    }
    count := 0
    for _, device := range devices {
        if device.Name == "any"{
            continue
        }else{
            count += 1
            go getPkgsByDeviceName(device.Name, networkflowsCh, limit, timeoutSecond)
        }
    }
    var networkflows [] gopacket.Flow
    for i := 0; i < count; i++ {
        flows := <- networkflowsCh
        networkflows = append(networkflows, flows...)
    }
    return networkflows
    // for _,flow := range networkflows {
    //     fmt.Println(flow)
    // }
}
