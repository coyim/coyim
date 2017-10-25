IMAGE_NAME=coyim/reproducible
REFERENCE_DATETIME=$(shell git log -1 --format=%ci)
PWD=$(shell pwd)

default: create-image-vivid

create-image-buster:
	docker build -t $(IMAGE_NAME):buster -f Dockerfile.buster .

create-image-vivid:
	docker build -t $(IMAGE_NAME):vivid -f Dockerfile.vivid .

create-image: create-image-vivid

console-vivid:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -i -t $(IMAGE_NAME):vivid bash

console-buster:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -i -t $(IMAGE_NAME):buster bash

console: console-buster

build-buster:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -t $(IMAGE_NAME):buster /root/build

build-vivid:
	docker run -v $(PWD)/../..:/src -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -t $(IMAGE_NAME):vivid /root/build

build: build-vivid
