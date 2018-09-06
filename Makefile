DOCKER_IMAGE := statusteam/statusd-bots

dependencies:
	dep ensure

image: AUTHOR = $(shell echo $$USER)
image: GIT_COMMIT = $(shell tag=`git describe --exact-match --tag 2>/dev/null`; \
	if [ $$? -eq 0 ]; then echo $$tag | sed 's/^v\(.*\)$$/\1/'; \
	else git rev-parse --short HEAD; fi)
image:
	docker build . \
		--label "commit=$(GIT_COMMIT)" \
		--label "author=$(AUTHOR)" \
		-t $(DOCKER_IMAGE):$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):latest
