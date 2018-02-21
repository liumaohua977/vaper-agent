#!/bin/bash
echo ''>./logs/vaper_agent.log
buildCMD='go build -o vaper_agent'
echo -e "\n[BUILD]"
echo -e "comand:$buildCMD"
$buildCMD

RUNCMD='sudo ./vaper_agent -f ./conf/config.ini -a start'
echo -e "\n[RUN]"
echo -e "comand:$RUNCMD"
$RUNCMD

echo -e "\n[LOG]"
cat ./logs/vaper_agent.log

