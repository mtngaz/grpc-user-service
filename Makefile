HOME_DIR = $(HOME)
CURR_DIR = $(PWD)
GOPATH_BIN = $(HOME_DIR)/go/bin
PROTO_DIR = $(CURR_DIR)/proto
SWAGGER_DIR = $(CURR_DIR)/swagger
API_DIR = $(CURR_DIR)/api

.PHONY: all install gen clean check setup

all: check clean install gen

install:
	@echo "📦 Installing tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/easyp-tech/easyp/cmd/easyp@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@echo "✅ Tools installed"

gen:
	@echo "🔧 Generating code..."
	@mkdir -p $(SWAGGER_DIR)
	@mkdir -p $(API_DIR)
	@easyp generate
	@echo "✅ Generation complete"

clean:
	@echo "🧹 Cleaning..."
	@rm -f $(GOPATH_BIN)/easyp
	@rm -f $(GOPATH_BIN)/swag
	@rm -f $(GOPATH_BIN)/protoc-gen-go
	@rm -f $(GOPATH_BIN)/protoc-gen-go-grpc
	@rm -f $(GOPATH_BIN)/protoc-gen-grpc-gateway
	@rm -f $(GOPATH_BIN)/protoc-gen-openapiv2
	@rm -rf $(API_DIR)/*.pb.go $(API_DIR)/*.gw.go
	@rm -rf $(SWAGGER_DIR)/*.json
	@echo "✅ Clean complete"

check:
	@echo "🔍 Checking environment..."
	@if [ ! -d "$(PROTO_DIR)" ]; then \
		echo "❌ Proto directory not found: $(PROTO_DIR)"; \
		exit 1; \
	fi
	@if [ -z "$$(ls -A $(PROTO_DIR)/*.proto 2>/dev/null)" ]; then \
		echo "❌ No .proto files found in $(PROTO_DIR)"; \
		exit 1; \
	fi
	@echo "✅ Found proto files:"
	@ls -la $(PROTO_DIR)/*.proto
	@echo ""
	@echo "📁 GOPATH/bin: $(GOPATH_BIN)"
	@echo "📋 PATH: $$PATH"

setup:
	@echo "🔧 Setting up project..."
	@go mod init grpc-user-service 2>/dev/null || true
	@go mod tidy
	@echo "✅ Setup complete"

help:
	@echo "Available commands:"
	@echo "  make all     - Full cycle: check → clean → install → gen"
	@echo "  make install - Install all plugins"
	@echo "  make gen     - Generate code (with auto-install)"
	@echo "  make clean   - Remove generated files and binaries"
	@echo "  make check   - Check environment and proto files"
	@echo "  make setup   - Initialize go module and tidy"
	@echo "  make help    - Show this help"