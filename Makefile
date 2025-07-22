BASE_IMAGE = golang:1.24-alpine3.20
LINT_IMAGE = golangci/golangci-lint:v2.3.0

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
