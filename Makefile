GO=CGO_ENABLED=0 go
GODOCKER=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
TAG=latest
BIN=bigbrother
IMAGE=dailyhotel/$(BIN)

build:
	glide install
	$(GO) build -a -installsuffix cgo -o bin/$(BIN) .

test: build
	$(GO) test -coverprofile=coverage.txt -covermode=atomic

image: clean
	glide install
	$(GODOCKER) build -a -installsuffix cgo -o bin/$(BIN) .
	docker build -t $(IMAGE):$(TAG) .

deploy: image
	docker push $(IMAGE):$(TAG)

.PHONY: clean

clean:
	rm -rf bin/
	rm -f coverage.txt

cleanall: clean
	rm -rf vendor/

