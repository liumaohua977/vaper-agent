#!/bin/bash
runVaper(){
  docker exec -d $1  /bin/sh -c "ps -ef|grep vaper_agent|grep -v grep|awk \"{print \$2}\"|xargs kill"
  docker exec -d $1  /bin/sh -c "mkdir -p /tmp/vaper && cd /tmp/vaper && curl -o vaper_agent.tar.gz http://10.0.2.15/vaper-agent/vaper_agent.tar.gz && tar xzf vaper_agent.tar.gz && chmod +x ./vaper_agent &&nohup ./vaper_agent -a start >>./vaper_agent.log 2>&1 &"
  echo "container "$1":reinstall vaper_agnet"
}

runVaper 'es01'
runVaper 'es02'
runVaper 'es03'
runVaper 'es04'
runVaper 'es05'
runVaper 'kibana'
runVaper 'kafka'
runVaper 'logstash'
runVaper 'metricbeat01'
runVaper 'metricbeat02'
runVaper 'metricbeat03'