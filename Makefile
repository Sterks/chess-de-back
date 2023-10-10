all: build push
build:
	docker build -t sterks/chess . --platform=linux/amd64
#--no-cache
.PHONY: build
push:
	docker push sterks/chess
