GOCMD=Go111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all: build
build:
	rm -rf target/
	mkdir target/
	cp config/develop.yaml target/develop.yaml
	$(GOBUILD) -o target/clover cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

run:
	nohup target/clover --conf_file ./config/develop.yaml > clover.log 2>&1 &


stop:
	pkill -f target/clover

.PHONY: build test clean run stop