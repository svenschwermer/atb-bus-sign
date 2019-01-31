# AtB Bus Sign

_Scrolling LED bus sign for Trondheim based AtB bus stops._

## Hardware
* [Mediatek LinkIt Smart 7688][1]
* 32×8 LED matrix display with 4 daisy-chained [MAX7219][2]-compatible controllers

```
LinkIt Smart 7688                      Display
┌───────────────┐                      ┌─────┐
│            5V ├──────────────────────┤ VCC │
│           GND ├──────────────────────┤ GND │
│           3V3 ├──────────┐           │     │
│               │         ┌┴┐          │     │
│               │     10k │R│          │     │
│               │         └┬┘          │     │
│  P25, SPI_CS1 ├──────────┴───────────┤ CS  │
│  P24, SPI_CLK ├──────────────────────┤ CLK │
│ P22, SPI_MOSI ├──────────────────────┤ DIN │
└───────────────┘                      └─────┘
```

### Gotchas
* The MAX7219 has a LOAD latch pin instead of a classic chip select pin (MAX7221).
* Having long traces on the SPI lines (CLK, MOSI) seems to cause trouble when
  Linux accesses the flash which hangs on the same SPI bus. JFFS errors are the
  result. Reducing the SPI clock frequency from 40 MHz to 1 MHz fixed this issue.
  This is probably overly conservative. A patch that makes the necessary changes
  is provided via the patch file in the root directory of this repository.

## Software
This software has been tested with OpenWrt v18.06.1. To build the entire image, do the following:
```
$ ./build.sh
```

## Housing
A simple housing is available in the `housing` directory. A bit of perfboard
holds the Linkit Smart 7688 module. 

[1]: https://docs.labs.mediatek.com/resource/linkit-smart-7688
[2]: https://datasheets.maximintegrated.com/en/ds/MAX7219-MAX7221.pdf
