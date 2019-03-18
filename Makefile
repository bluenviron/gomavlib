
.PHONY: help
help:
	@echo "usage: make [action] [args...]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  format              format source files."
	@echo ""
	@echo "  test                run available tests."
	@echo ""
	@echo "  gen-dialects        generate dialects."
	@echo ""
	@echo "  run-example [name]  run example with given name."
	@echo ""


.PHONY: format
format:
	@docker run --rm -it \
		-v $(PWD):/src \
		amd64/golang:1.11-stretch \
		sh -c "cd /src \
		&& find . -type f -name '*.go' | xargs gofmt -l -w -s"


.PHONY: test
test:
	@docker run --rm -it \
		-v $(PWD):/src \
		amd64/golang:1.11-stretch \
		sh -c "cd /src \
		&& go test -v . \
		&& go install ./dialgen"


.PHONY: gen-dialects
define GEN_DIALECTS_DOCKERFILE
FROM amd64/golang:1.11-stretch
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go install ./dialgen
endef
export GEN_DIALECTS_DOCKERFILE
gen-dialects:
	@for DIALECT in ASLUAV ardupilotmega autoquad common icarous matrixpilot minimal paparazzi slugs standard test uAvionix ualberta; do \
		echo "$$GEN_DIALECTS_DOCKERFILE" | docker build -q . -f - -t gomavlib-gen-dialects \
			&& docker run --rm -it \
			-v $(PWD):/src \
			gomavlib-gen-dialects \
			dialgen --output=dialects/$$DIALECT/dialect.go \
			https://raw.githubusercontent.com/mavlink/mavlink/master/message_definitions/v1.0/$$DIALECT.xml \
			|| exit 1; \
	done


.PHONY: run-example
ifeq (run-example, $(firstword $(MAKECMDGOALS)))
  $(eval %:;@:) # do not treat arguments as targets
  ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  EXAMPLE := $(word 1, $(ARGS))
endif
run-example:
	@docker run --rm -it \
		--privileged \
		-v $(PWD):/src \
		amd64/golang:1.11-stretch \
		sh -c "cd /src && go run example/$(EXAMPLE).go"
