APP=gh-art

.PHONY: build install clean

build:
	go build -o bin/$(APP)

install:
	go install

clean:
	rm -rf bin
