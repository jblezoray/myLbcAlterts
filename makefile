BINARY=myLbcAlterts

VERSION=`git describe --tags`
BUILD=`date +%Y-%m-%d/%H:%M:%S`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	go build ${LDFLAGS} -o ${BINARY}

test:
	go test

install: test
	go install ${LDFLAGS_f1}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
