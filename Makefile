receiver:
	@echo "Starting receiver..."
	@go build -o bin/receiver ./number_receiver
	@./bin/receiver

generator:
	@echo "Starting numbers generator..."
	@go build -o bin/generator ./numbers_generator
	@./bin/generator

checker:
	@echo "Starting result checker..."
	@go build -o bin/checker ./result_checker
	@./bin/checker

annoucer:
	@echo "Starting annoucer checker..."
	@go build -o bin/annoucer ./result_annoucer
	@./bin/annoucer