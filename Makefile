
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

DIALECTS = ASLUAV ardupilotmega autoquad common icarous matrixpilot minimal \
	paparazzi slugs standard test uAvionix ualberta
gen-dialects:
	echo "FROM amd64/golang:1.11-stretch \n\
		WORKDIR /src \n\
		COPY go.mod go.sum ./ \n\
		RUN go mod download \n\
		COPY *.go ./ \n\
		COPY dialgen ./dialgen \n\
		RUN go install ./dialgen" | docker build -q . -f - -t gomavlib-gen-dialects
	@for DIALECT in $(DIALECTS); do \
		docker run --rm -it \
		-v $(PWD):/src \
		gomavlib-gen-dialects \
		dialgen --output=dialects/$$DIALECT/dialect.go \
		https://raw.githubusercontent.com/mavlink/mavlink/master/message_definitions/v1.0/$$DIALECT.xml \
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
