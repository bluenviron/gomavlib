
BASE_IMAGE = amd64/golang:1.15-alpine3.12

.PHONY: $(shell ls)

help:
	@echo "usage: make [action]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy              run go mod tidy"
	@echo "  format                format source files"
	@echo "  test                  run all available tests"
	@echo "  dialects              generate dialects"
	@echo "  run-example E=[name]  run example by name"
	@echo ""

blank :=
define NL

$(blank)
endef

mod-tidy:
	docker run --rm -it -v $(PWD):/s $(BASE_IMAGE) \
	sh -c "apk add git && cd /s && go get && go mod tidy"

format:
	docker run --rm -it -v $(PWD):/s $(BASE_IMAGE) \
	sh -c "cd /s && find . -type f -name '*.go' | xargs gofmt -l -w -s"

define DOCKERFILE_TEST
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make gcc musl-dev
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
endef
export DOCKERFILE_TEST

test:
	echo "$$DOCKERFILE_TEST" | docker build . -f - -t temp
	docker run --rm -it temp make test-nodocker

test-nodocker:
	go test -race -v ./...
	go build -o /dev/null ./cmd/...
	$(foreach f,$(shell ls examples/*),go build -o /dev/null $(f)$(NL))

define DOCKERFILE_GEN_DIALECTS
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
endef
export DOCKERFILE_GEN_DIALECTS

dialects:
	echo "$$DOCKERFILE_GEN_DIALECTS" | docker build . -f - -t temp
	docker run --rm -it -v $(PWD):/s temp \
	make dialects-nodocker

dialects-nodocker:
	$(eval export CGO_ENABLED = 0)
	go run ./cmd/dialects-gen
	find ./pkg/dialects -type f -name '*.go' | xargs gofmt -l -w -s

run-example:
	docker run --rm -it \
	--privileged \
	--network=host \
	-v $(PWD):/s $(BASE_IMAGE) \
	sh -c "cd /s && go run examples/$(E).go"
