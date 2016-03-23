IMAGE_NAME=coyim/reproducible

default:
	docker build -t $(IMAGE_NAME) -f Dockerfile .

build:
	docker run -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" $(IMAGE_NAME)

run:
	docker run -e REFERENCE_DATETIME="$(REFERENCE_DATETIME)" -i -t $(IMAGE_NAME) bash

