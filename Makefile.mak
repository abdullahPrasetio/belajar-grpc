create:
	protoc --proto_path=proto proto/*.proto  --go_out=models
	protoc --proto_path=proto proto/*.proto  --go-grpc_out=models
	protoc -I . --grpc-gateway_out ./models/ \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --proto_path=proto proto/*.proto
	protoc -I . --openapiv2_out ./swagger \
    --openapiv2_opt logtostderr=true \
    --proto_path=. proto/*.proto
	protoc \
	-I . \
	--go_out=models \
	--validate_out="lang=go:models" \
	--proto_path=proto proto/*.proto

clean:
	rm models/proto/*.go