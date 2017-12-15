all: deps build-release

.PHONY: deps
deps:
	$(GOPATH)/bin/glide install

.PHONY: build-debug
build-debug:
	$(GOPATH)/bin/go-bindata -debug -prefix public public/... && go build -o editor .

.PHONY: build-release
build-release:
	$(GOPATH)/bin/go-bindata -prefix public public/... && GOOS=linux go build -o editor .

.PHONY: run
run:
	./editor
