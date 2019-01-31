#!/usr/bin/env bash

set -x -e

openwrt_tag=v18.06.1

git clone https://github.com/openwrt/openwrt.git
cd openwrt
git checkout $openwrt_tag
patch -Np1 < ../0001-Reduce-max-SPI-clock-frequency-on-flash.patch
grep -E '^src-git packages ' feeds.conf.default > feeds.conf
echo "src-git bus https://github.com/svenschwermer/atb-bus-sign.git" >> feeds.conf
./scripts/feeds update -a
./scripts/feeds install -a
cp ../openwrt-$openwrt_tag-config .config
make -j$(nproc)
