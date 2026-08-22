define DOCKERFILE_DIALECTS
FROM $(BASE_IMAGE)
RUN apk add --no-cache git make
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
endef
export DOCKERFILE_DIALECTS

dialects:
	echo "$$DOCKERFILE_DIALECTS" | docker build . -f - -t temp
	docker run --rm -it -v $(shell pwd):/s temp \
	make dialects-nodocker
	make format

dialects-nodocker:
	$(eval export CGO_ENABLED = 0)
	go run ./cmd/dialects-gen
