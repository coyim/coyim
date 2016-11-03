IMAGE_NAME=coyim/reproducible
REFERENCE_DATETIME=$(shell git log -1 --format=%ci)
PWD=$(shell pwd)

default: create-image

create-image:
	docker build -t $(IMAGE_NAME) -f Dockerfile .

console:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -i -t $(IMAGE_NAME) bash

build:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -t $(IMAGE_NAME) /root/build
