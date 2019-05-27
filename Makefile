
help:
	@echo "usage: make [action] [args...]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy              run go mod tidy"
	@echo "  format                format source files"
	@echo "  test                  run all available tests"
	@echo "  gen-dialects          generate dialects"
	@echo "  run-example E=[name]  run example by name"
	@echo ""

blank :=
define NL

$(blank)
endef

mod-tidy:
	docker run --rm -it -v $(PWD):/src \
	amd64/golang:1.11 \
	sh -c "cd /src && go get -m ./... && go mod tidy"

format:
	@docker run --rm -it \
	-v $(PWD):/src \
	amd64/golang:1.11-stretch \
	sh -c "cd /src \
	&& find . -type f -name '*.go' | xargs gofmt -l -w -s"

test:
	echo "FROM amd64/golang:1.11-stretch \n\
	WORKDIR /src \n\
	COPY go.mod go.sum ./ \n\
	RUN go mod download \n\
	COPY Makefile *.go ./ \n\
	COPY dialgen ./dialgen \n\
	COPY dialects ./dialects \n\
	COPY example ./example" | docker build . -f - -t gomavlib-test
	docker run --rm -it gomavlib-test make test-nodocker

test-nodocker:
	go test -v .
	go build -o /dev/null ./dialgen
	$(foreach f, $(shell ls example/*), go build -o /dev/null $(f)$(NL))

gen-dialects:
	echo "FROM amd64/golang:1.11-stretch \n\
	WORKDIR /src \n\
	COPY go.mod go.sum ./ \n\
	RUN go mod download \n\
	COPY *.go ./ \n\
	COPY dialgen ./dialgen" | docker build -q . -f - -t gomavlib-gen-dialects
	docker run --rm -it \
	-v $(PWD):/src \
	gomavlib-gen-dialects \
	make gen-dialects-nodocker

gen-dialects-nodocker:
	$(eval COMMIT = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/commits/master \
	| grep -o '"sha": ".\+"' | sed 's/"sha": "\(.\+\)"/\1/' | head -n1))
	$(eval DIALECTS = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/contents/message_definitions/v1.0?ref=$(COMMIT) \
	| grep -o '"name": ".\+\.xml"' | sed 's/"name": "\(.\+\)\.xml"/\1/'))
	rm -rf dialects/*
	echo "package dialects" > dialects/dialects.go
	$(foreach d, $(DIALECTS), go run ./dialgen -q --output=dialects/$(d)/dialect.go \
	--preamble="Generated from revision https://github.com/mavlink/mavlink/tree/$(COMMIT)" \
	https://raw.githubusercontent.com/mavlink/mavlink/$(COMMIT)/message_definitions/v1.0/$(d).xml$(NL))
	find ./dialects -type f -name '*.go' | xargs gofmt -l -w -s

run-example:
	@docker run --rm -it \
	--privileged \
	--network=host \
	-v $(PWD):/src \
	amd64/golang:1.11-stretch \
	sh -c "cd /src && go run example/$(E).go"
