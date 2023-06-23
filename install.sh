#!/usr/bin/env bash

set -e

wget https://github.com/nexryai/lance-light/releases/latest/download/LanceLight-linux-amd64.zip
unzip LanceLight-linux-amd64.zip
rm LanceLight-linux-amd64.zip
mv llfctl /usr/bin/

curl https://raw.githubusercontent.com/nexryai/lance-light/main/systemd/lance.service > /etc/systemd/system/lance.service
curl https://raw.githubusercontent.com/nexryai/lance-light/main/config.default.yml > /etc/lance.yml