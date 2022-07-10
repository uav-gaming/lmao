.PHONY: build clean deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/lmao lmao.go

clean:
	rm -rf ./bin

deploy: build
	sls deploy --verbose function -f lmao

deploy-full: clean build
	sls deploy --verbose
