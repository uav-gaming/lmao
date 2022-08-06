.PHONY: build clean deploy

BUILD_TIME=$(shell date +%s)

GOFLAGS+=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
LDFLAGS+=-X main.BuildTime=${BUILD_TIME}

build:
	${GOFLAGS} go build -o bin/lmao -v -ldflags "${LDFLAGS}" lmao.go 

clean:
	rm -rf ./bin

deploy: build
	sls deploy function --verbose -f lmao

deploy-full: clean build
	sls deploy --verbose
