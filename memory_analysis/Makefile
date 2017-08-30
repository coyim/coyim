IMAGE_NAME=coyim/memory_analysis
REFERENCE_DATETIME=$(shell git log -1 --format=%ci)
PWD=$(shell pwd)

default: create-image

create-image:
	docker build -t $(IMAGE_NAME) -f Dockerfile .

console:
	docker run -v $(PWD)/..:/src -i -t $(IMAGE_NAME) bash

build:
	docker run -v $(PWD)/..:/src -t $(IMAGE_NAME) /root/build-coyim
