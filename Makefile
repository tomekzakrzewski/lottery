receiver:
	@echo "Starting receiver..."
	@go build -o bin/receiver ./number_receiver
	@./bin/receiver


generator:
	@echo "Starting numbers generator..."
	@go build -o bin/generator ./numbers_generator
	@./bin/generator