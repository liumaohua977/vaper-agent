#!/bin/bash
tag=`date "+%Y-%m-%d-%H:%M:%S"`
tmpdir=/tmp/nginx/vaper-agent/
rm -rf $tmpdir
mkdir -p $tmpdir

# cp -f /home/hxn/go/src/github.com/vaper/conf/config.default.ini $tmpdir/conf/config.default.ini

cp -f /home/hxn/go/src/github.com/vaper/vaper_agent $tmpdir/vaper_agent
cp -f /home/hxn/go/src/github.com/vaper/vaper_agent.ini $tmpdir/vaper_agent.ini

cd $tmpdir
tar czf vaper_agent.tar.gz ./*

cd /home/hxn/go/src/github.com/vaper
sh docker_container_install.sh
