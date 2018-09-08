# AtB Bus Sign

_Scrolling LED bus sign for Trondheim based AtB bus stops._

## Hardware
* [Mediatek LinkIt Smart 7688](https://docs.labs.mediatek.com/resource/linkit-smart-7688)
* 32×8 LED matrix display with 4 daisy-chained [MAX7219](https://datasheets.maximintegrated.com/en/ds/MAX7219-MAX7221.pdf)-compatible controllers

```
LinkIt Smart 7688       Display
+---------------+       +-----+
|            5V |<----->| VCC |
|           GND |<----->| GND |
| P22, SPI_MOSI |<----->| DIN |
|  P25, SPI_CS1 |<----->| CS  |
|  P24, SPI_CLK |<----->| CLK |
+---------------+       +-----+
```

## Software
This software has been tested with OpenWrt v18.06.1. To build the entire image, do the following:
```
$ git clone https://github.com/openwrt/openwrt.git
$ cd openwrt
$ git checkout v18.06.1
$ cp feeds.conf{.default,}
$ echo "src-git atb-bus-sign https://github.com/svenschwermer/atb-bus-sign.git" >> feeds.conf
$ ./scripts/feeds update -a
$ ./scripts/feeds install -a
$ curl -o .config https://raw.githubusercontent.com/svenschwermer/atb-bus-sign/cpp/openwrt-v18.06.1-config
$ make
```
