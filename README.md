# Vaper-agent

Agent is a golang project. Collect netflow data and hostmeta info from operating system.
[https://github.com/vapering/vaper-agent](https://github.com/vapering/vaper-agent)

## Deploy in production environment

Vaper-agent need two files to run:

- vaper_agent
- vaper_agent.ini

### Run

`./vaper_agent -a start`

### Run in daemon

`nohup ./vaper_agent -a start >>./vaper_agent.log 2>&1 &`

### Example

```bash
mkdir -p /tmp/vaper
cd /tmp/vaper
curl -o vaper_agent.tar.gz http://10.0.2.15/vaper-agent/vaper_agent.tar.gz
tar xzf vaper_agent.tar.gz
chmod +x ./vaper_agent
nohup ./vaper_agent -a start >>./vaper_agent.log 2>&1 &
```

## Build and run in the development environment

`sh buildRun.sh`
The buid result is vaper_agent

## Something More

Vaper-agent need `libpcap` in the development environment. Nobody want to waste time in install libpcap on every host. So we need to make sure that  compiling the vaper_agent statically.

Find (`#cgo linux LDFLAGS: -lpcap`) in file (pcap.go), and change to something like (`#cgo linux LDFLAGS: -L /tmp/nginx/libpcap-1.8.1 -lpcap`)