receiver:
	@echo "Starting receiver..."
	@go build -o bin/receiver ./number_receiver
	@./bin/receiver