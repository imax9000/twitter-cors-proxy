.PHONY: all build deploy

gcloud := gcloud
ifdef PROJECT
	gcloud := gcloud --project=$(PROJECT)
endif

all: build

build: twitter-cors-proxy

twitter-cors-proxy: main.go config.go go.mod go.sum
	go build -v

# deployment doesn't actually require us to build a binary locally, but we
# depend on it here to make sure that the code actually compiles.
deploy: build
	$(gcloud) --quiet app deploy
