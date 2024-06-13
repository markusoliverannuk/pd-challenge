all: buildapp app

buildapp:
	@echo "Building api..."
	@go build -o ./build/app ./cmd/main.go

app: buildapp
	@echo "Running API..."
	@./build/app

.PHONY: all buildapp app
