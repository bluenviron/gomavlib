run-example:
	docker run --rm -it \
	--privileged \
	--network=host \
	-v $(PWD):/s -w /s \
	$(BASE_IMAGE) \
	sh -c "go run examples/$(E).go"
