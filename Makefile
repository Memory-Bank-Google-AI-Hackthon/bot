NAMESPACE := weichen-lin
REPO_NAME := bot
TAG := $(GIT_SHA)
IMAGE_BASE := ghcr.io/$(NAMESPACE)/$(REPO_NAME)

build-prod:
	docker build -t $(IMAGE_BASE):latest .

push-prod:
	docker push $(IMAGE_BASE):latest