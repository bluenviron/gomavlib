define DOCKERFILE_DIALECTS
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make
RUN go install mvdan.cc/gofumpt@v0.5.0
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
endef
export DOCKERFILE_DIALECTS

dialects:
	echo "$$DOCKERFILE_DIALECTS" | docker build . -f - -t temp
	docker run --rm -it -v $(PWD):/s temp \
	make dialects-nodocker

dialects-nodocker:
	$(eval export CGO_ENABLED = 0)
	go run ./cmd/dialects-gen
	find ./pkg/dialects -type f -name '*.go' | xargs gofumpt -l -w
