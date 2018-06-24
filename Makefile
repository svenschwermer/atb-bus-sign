build:
	# TODO: Remove once spi is compatiple with MIPS
	sed -i.org 's/0x40006B00/0x80006B00/' vendor/golang.org/x/exp/io/spi/devfs.go
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w"

.PHONY: build
