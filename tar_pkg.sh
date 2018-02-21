#!/bin/bash
tag=`date "+%Y-%m-%d-%H:%M:%S"`
tmpdir=/tmp/nginx/vaper-agent/$tag
mkdir -p $tmpdir
mkdir -p $tmpdir/conf
mkdir -p $tmpdir/logs

cp -f /home/hxn/go/src/github.com/vaper/vaper_agent $tmpdir/vaper_agent
cp -f /home/hxn/go/src/github.com/vaper/conf/config.default.ini $tmpdir/conf/config.default.ini
cp -f /home/hxn/go/src/github.com/vaper/conf/config.ini $tmpdir/conf/config.ini

cd $tmpdir
zip -r vaper.zip ./*
cd /home/hxn/go/src/github.com/vaper