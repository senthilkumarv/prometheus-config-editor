all: deps build-release

.PHONY: deps
deps:
	$(GOPATH)/bin/glide install

.PHONY: build-debug
build:
	$(GOPATH)/bin/go-bindata -debug -prefix public public/... && go build -o editor .

.PHONY: build-release
build-prod:
	$(GOPATH)/bin/go-bindata -prefix public public/... && GOOS=linux go build -o editor .

.PHONY: run
run:
	./editor
