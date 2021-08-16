generate: # Generates loaders for config proto.
	@echo "Generating the loaders"
	@go get ./...
	@go run main.go -o $(LOADERS_OUT_DIR) $(PROTO_FILE)
	@echo "Done."
