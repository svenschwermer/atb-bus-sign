build:
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w"

sim:
	go build -tags=sim -o sim

.PHONY: build sim
