.PHONY: build container

build: container
	docker run \
		-v $(shell pwd):/usr/local/src \
		-w /usr/local/src \
		mipsel-build \
			mipsel-linux-gnu-g++ \
			-pthread -lssl -lcrypto -lboost_system \
			-O2 -s -static \
			-o bus-sign \
			main.cpp atb.cpp font.cpp frame.cpp max7219.cpp

container:
	docker build -t mipsel-build .
