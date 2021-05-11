#!/bin/bash

ip=$(docker ps | grep likakuli | awk '{print $1}' | xargs docker inspect | grep \"IPAddress\" | head -n 1 | awk -F '\"' '{print $4}')

sed s/ip/$ip/g ./conf/config.toml.tmpl > ./conf/config.toml
