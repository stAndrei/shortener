.PHONY: dev build image test clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
APP=shortener
REPO?=prologic/$(APP)
TAG?=latest
BUILD?=-dev

all: dev

dev: build
	@./$(APP)

build: clean
	@echo " -> Building $(TAG)$(BUILD)"
	@go build -tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X github.com/$(REPO).GitCommit=$(COMMIT) -X github.com/$(REPO).Build=$(BUILD)" .
	@echo "Built $$(./$(APP) -v)"

image:
	@docker build --build-arg TAG=$(TAG) --build-arg BUILD=$(BUILD) -t $(REPO):$(TAG) .
	@echo "Image created: $(REPO):$(TAG)"

profile:
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench .

bench:
	@go test -v -bench .

test:
	@go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic .

clean:
	@git clean -f -d -X
