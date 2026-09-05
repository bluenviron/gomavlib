BASE_IMAGE = golang:1.27-alpine3.24
LINT_IMAGE = golangci/golangci-lint:v2.13.2
NODE_IMAGE = node:24-alpine3.24

.PHONY: $(shell ls)

help:
	@echo "usage: make [action]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  format                format source files"
	@echo "  test                  run tests"
	@echo "  lint                  run linter"
	@echo "  dialects              generate dialects"
	@echo ""

blank :=
define NL

$(blank)
endef

include scripts/*.mk
