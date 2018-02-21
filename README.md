# Vaper-agnet
Vaper is a tool to get the topology of any systems by collecting the network traffic information among the hosts.
# Install
yum install tcpdump unzip -y 

docker exec -it es01 /bin/bash
rm -rf /tmp/vaper-agent
mkdir /tmp/vaper-agent
cd /tmp/vaper-agent
curl http://10.0.2.15/vaper-agent/2017-12-16-14%3A19%3A33/vaper.zip -o vaper.zip
unzip vaper.zip
./vaper_agent -f ./conf/config.ini -a init
nohup ./vaper_agent -f ./conf/config.ini -a start >/dev/null 2>&1 &
exit
