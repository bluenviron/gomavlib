BASE_IMAGE = golang:1.21-alpine3.18
LINT_IMAGE = golangci/golangci-lint:v1.55.0

.PHONY: $(shell ls)

help:
	@echo "usage: make [action]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy              run go mod tidy"
	@echo "  format                format source files"
	@echo "  test                  run tests"
	@echo "  lint                  run linter"
	@echo "  dialects              generate dialects"
	@echo "  run-example E=[name]  run example by name"
	@echo ""

blank :=
define NL

$(blank)
endef

include scripts/*.mk
