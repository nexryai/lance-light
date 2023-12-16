#!/usr/bin/env bash

set -e

required_commands=("wget" "curl" "unzip" "nft")

for cmd in "${required_commands[@]}"; do
    if command -v "$cmd" > /dev/null 2>&1; then
        echo "Checking $cmd >> OK!"
    else
        echo "$cmd is not installed. Please install it first."
        exit 1
    fi
done

ARCHITECTURE=$(uname -m)

case $ARCHITECTURE in
    x86_64)
        ARCH="amd64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    mips64)
        ARCH="mips64"
        ;;
    *)
        echo "Unsupported ARCHITECTURE >_< : $ARCHITECTURE"
        exit 1
        ;;
esac

ZIPNAME="LanceLight-linux-$ARCH.zip"

wget "https://github.com/nexryai/lance-light/releases/latest/download/$ZIPNAME"
unzip $ZIPNAME
rm $ZIPNAME
mv llfctl /usr/bin/

curl https://raw.githubusercontent.com/nexryai/lance-light/main/systemd/lance.service > /etc/systemd/system/lance.service

if [ ! -f /etc/lance.yml ]; then
  curl https://raw.githubusercontent.com/nexryai/lance-light/main/config.default.yml > /etc/lance.yml
fi

# restoreconがあるならSELinuxに怒られないようにする
if command -v restorecon >/dev/null 2>&1; then
    restorecon -R /usr/bin/llfctl
fi

systemctl daemon-reload

echo "Lance Light has been installed successfully. enjoy! :)"