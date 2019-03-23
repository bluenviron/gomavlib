
help:
	@echo "usage: make [action] [args...]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  format              format source files."
	@echo "  test                run all available tests."
	@echo "  gen-dialects        generate dialects."
	@echo "  run-example [name]  run example with given name."
	@echo ""

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
		RUN go mod download" | docker build . -f - -t gomavlib-test \
		&& docker run --rm -it \
		-v $(PWD):/src \
		gomavlib-test \
		sh -c "go test -v . \
		&& go install ./dialgen"

DIALECTS = ASLUAV ardupilotmega autoquad common icarous matrixpilot minimal \
	paparazzi slugs standard test uAvionix ualberta
gen-dialects:
	echo "FROM amd64/golang:1.11-stretch \n\
		WORKDIR /src \n\
		COPY go.mod go.sum ./ \n\
		RUN go mod download \n\
		COPY . ./ \n\
		RUN go install ./dialgen" | docker build -q . -f - -t gomavlib-gen-dialects
	@for DIALECT in $(DIALECTS); do \
		docker run --rm -it \
		-v $(PWD):/src \
		gomavlib-gen-dialects \
		dialgen --output=dialects/$$DIALECT/dialect.go \
		https://raw.githubusercontent.com/mavlink/mavlink/master/message_definitions/v1.0/$$DIALECT.xml \
		|| exit 1; \
	done

ifeq (run-example, $(firstword $(MAKECMDGOALS)))
  $(eval %:;@:) # do not treat arguments as targets
  ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  EXAMPLE := $(word 1, $(ARGS))
endif
run-example:
	@docker run --rm -it \
		--privileged \
		--network=host \
		-v $(PWD):/src \
		amd64/golang:1.11-stretch \
		sh -c "cd /src && go run example/$(EXAMPLE).go"
