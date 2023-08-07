define DOCKERFILE_TEST
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make gcc musl-dev
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
endef
export DOCKERFILE_TEST

test-cmd:
	go build -o /dev/null ./cmd/...
	go test -v -race -coverprofile=coverage-cmd.txt ./cmd/...

test-examples:
	go build -o /dev/null ./examples/...

test-pkg:
	go test -v -race -coverprofile=coverage-pkg.txt ./pkg/...

test-root:
	go test -v -race -coverprofile=coverage-root.txt .

test-nodocker: test-cmd test-examples test-pkg test-root

test:
	echo "$$DOCKERFILE_TEST" | docker build . -f - -t temp
	docker run --rm -it temp make test-nodocker
