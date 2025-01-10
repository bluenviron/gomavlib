run-example:
	docker run --rm -it \
	--privileged \
	--network=host \
	-v $(shell pwd):/s -w /s \
	$(BASE_IMAGE) \
	sh -c "go run examples/$(E).go"
