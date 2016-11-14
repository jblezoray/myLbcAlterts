SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=lbcAlerts

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} ${SOURCES}

run: $(BINARY)
	./$(BINARY)

.PHONY: clean
clean:
	rm -f ${BINARY}
