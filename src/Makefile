.PHONY: build container

build: container
	docker run \
		-v $(shell pwd):/usr/local/src \
		-w /usr/local/src \
		mipsel-build \
			mipsel-linux-gnu-g++ \
			-lssl -lcrypto -lboost_system \
			-O2 \
			-o bus-sign \
			max7219.cpp

container:
	docker build -t mipsel-build .