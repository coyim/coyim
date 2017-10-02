IMAGE_NAME=coyim/coyim-ci

default: image build

image:
	docker build -t $(IMAGE_NAME) .

build:
	docker run $(IMAGE_NAME)

run:
	docker run -i -t $(IMAGE_NAME) bash
