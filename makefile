BINARY=myLbcAlterts

VERSION=`git describe --tags`
BUILD=`date +%Y-%m-%d/%H:%M:%S`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	go build ${LDFLAGS} -o ${BINARY}

test:
	go test -v

testrun: clean build test 
	./${BINARY} MyConf_v0.2.json --migratedb
	./${BINARY} MyConf_v0.2.json --analyze

install: test
	go install ${LDFLAGS_f1}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
