all: build run

build:
	@echo "Building client..."
	@cd my-react-app && npm install && npm run build

run: build
	@echo "Running client..."
	@cd my-react-app && npm start

clean:
	@echo "Cleaning client build..."
	@cd my-react-app && npm run clean

.PHONY: all build run clean