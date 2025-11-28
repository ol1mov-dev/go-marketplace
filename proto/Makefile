USER_PROTO_DIR=user/v1
USER_OUT_DIR=pkg/user/v1

.PHONY: proto

proto:
	@echo "🧩 Generating user gRPC code from proto files..."
	protoc \
		--proto_path=$(USER_PROTO_DIR) \
		--go_out=$(USER_OUT_DIR) \
		--go-grpc_out=$(USER_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$$(find $(USER_PROTO_DIR) -name "*.proto")
	@echo "✅ Generation complete!"

.PHONY: clean
clean:
	rm -rf pkg/user
