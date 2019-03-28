
help:
	@echo "usage: make [action] [args...]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy            run go mod tidy."
	@echo "  format              format source files."
	@echo "  test                run all available tests."
	@echo "  gen-dialects        generate dialects."
	@echo "  run-example [name]  run example with given name."
	@echo ""

# do not treat arguments as targets
%:
	@[ "$(word 1, $(MAKECMDGOALS))" != "$@" ] || { echo "unrecognized command."; exit 1; }

ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))

mod-tidy:
	docker run --rm -it -v $(PWD):/src amd64/golang:1.11 \
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
		COPY *.go ./ \n\
		RUN go test -v . \n\
		COPY dialgen ./dialgen \n\
		RUN go install ./dialgen" | docker build . -f - -t gomavlib-test

gen-dialects:
	echo "FROM amd64/golang:1.11-stretch \n\
		WORKDIR /src \n\
		COPY go.mod go.sum ./ \n\
		RUN go mod download \n\
		COPY *.go ./ \n\
		COPY dialgen ./dialgen \n\
		RUN go install ./dialgen" | docker build -q . -f - -t gomavlib-gen-dialects
	docker run --rm -it \
		-v $(PWD):/src \
		gomavlib-gen-dialects \
		make gen-dialects-nodocker

gen-dialects-nodocker:
	$(eval COMMIT = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/commits/master \
		| grep -o '"sha": ".\+"' | sed 's/"sha": "\(.\+\)"/\1/' | head -n1))
	$(eval DIALECTS = $(shell curl -s -L https://api.github.com/repos/mavlink/mavlink/contents/message_definitions/v1.0?ref=$(COMMIT) \
		| grep -o '"name": ".\+\.xml"' | sed 's/"name": "\(.\+\)\.xml"/\1/'))
	@for DIALECT in $(DIALECTS); do \
		dialgen --output=dialects/$$DIALECT/dialect.go \
			--preamble="Generated from revision https://github.com/mavlink/mavlink/tree/$(COMMIT)" \
			https://raw.githubusercontent.com/mavlink/mavlink/$(COMMIT)/message_definitions/v1.0/$$DIALECT.xml \
			|| exit 1; \
	done

run-example:
	$(eval EXAMPLE := $(word 1, $(ARGS)))
	@docker run --rm -it \
		--privileged \
		--network=host \
		-v $(PWD):/src \
		amd64/golang:1.11-stretch \
		sh -c "cd /src && go run example/$(EXAMPLE).go"
