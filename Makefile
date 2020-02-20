.PHONY: image

IMAGE_NAME ?= codeclimate/codeclimate-keepachangelog

image:
	docker build --tag "$(IMAGE_NAME)" .
test: image
	docker run --rm -v "$(PWD):/code" -v "$(PWD)/sample-config.json:/engine.json" -w /code "$(IMAGE_NAME)"
