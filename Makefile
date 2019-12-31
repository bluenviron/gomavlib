
.PHONY: $(shell ls)

BASE_IMAGE = amd64/golang:1.13-alpine3.10

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
RUN apk add --no-cache git make
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
	$(eval export CGO_ENABLED = 0)
	go test -v ./...
	go build -o /dev/null ./dialgen
	$(foreach f,$(shell ls example/*),go build -o /dev/null $(f)$(NL))

define DOCKERFILE_GEN_DIALECTS
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make curl
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY dialgen ./dialgen
endef
export DOCKERFILE_GEN_DIALECTS

dialects:
	echo "$$DOCKERFILE_GEN_DIALECTS" | docker build . -f - -t temp
	docker run --rm -it -v $(PWD):/s temp \
	make dialects-nodocker

dialects-nodocker:
	$(eval export CGO_ENABLED = 0)
	$(eval COMMIT = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/commits/master \
	| grep -o '"sha": ".\+"' | sed 's/"sha": "\(.\+\)"/\1/' | head -n1))
	$(eval DIALECTS = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/contents/message_definitions/v1.0?ref=$(COMMIT) \
	| grep -o '"name": ".\+\.xml"' | sed 's/"name": "\(.\+\)\.xml"/\1/'))
	rm -rf dialects/*
	echo "package dialects" > dialects/dialects.go
	$(foreach d,$(DIALECTS),go run ./dialgen -q --output=dialects/$(subst _,,$(d))/dialect.go \
	--preamble="Generated from revision https://github.com/mavlink/mavlink/tree/$(COMMIT)" \
	https://raw.githubusercontent.com/mavlink/mavlink/$(COMMIT)/message_definitions/v1.0/$(d).xml$(NL))
	find ./dialects -type f -name '*.go' | xargs gofmt -l -w -s

run-example:
	docker run --rm -it \
	--privileged \
	--network=host \
	-v $(PWD):/s $(BASE_IMAGE) \
	sh -c "cd /s && go run example/$(E).go"
