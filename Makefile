build:
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build

.PHONY: build
