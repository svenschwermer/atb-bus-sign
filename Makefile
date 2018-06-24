build:
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w"

.PHONY: build
