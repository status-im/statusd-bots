DOCKER_IMAGE := statusteam/statusd-bots

dependencies:
	dep ensure

image: AUTHOR = $(shell echo $$USER)
image: GIT_COMMIT = $(shell git rev-parse --short HEAD)
image:
	docker build . \
		--label "commit=$(GIT_COMMIT)" \
		--label "author=$(AUTHOR)" \
		-t $(DOCKER_IMAGE):$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):latest

build: bin/pubchats bin/bench-mailserver bin/x-check-mailserver

bin/pubchats:
	go build -o ./bin/pubchats ./cmd/pubchats
bin/bench-mailserver:
	go build -o ./bin/bench-mailserver ./cmd/bench-mailserver
bin/x-check-mailserver:
	go build -o ./bin/x-check-mailserver ./cmd/x-check-mailserver

clean:
	rm bin/*
