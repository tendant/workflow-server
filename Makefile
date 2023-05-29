objects = main

all: $(objects)

$(objects):
	go build -ldflags $(LDFLAGS)  -o $@ $@.go

dep:
	go mod tidy

vendor:
	go mod vendor

clean:
	go clean
	rm -f $(objects)

.PHONY: clean
