#!/bin/sh

cd release
tar xf terraform-provider-zabbix*.tar.gz
cd terraform-provider-zabbix

for arch in *;
  do tar czf ${arch}.tar.gz $arch/*
done