#!/bin/bash
echo ''>./vaper_agent.log
buildCMD='go build -o vaper_agent'
echo -e "\n[BUILD]"
echo -e "comand:$buildCMD"
$buildCMD
if [ $? -ne 0 ] ; then
  echo ">>> Error in build work, exit."
  exit 1
else
  echo ">>> Build success."
fi

RUNCMD='sudo ./vaper_agent -f ./config.ini -a start'
echo -e "\n[RUN]"
echo -e "comand:$RUNCMD"
$RUNCMD
