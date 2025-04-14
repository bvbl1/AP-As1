generate:
  protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/inventory.proto