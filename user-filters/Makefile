proto:
	rm -f pkg/pb/*.go
	protoc --proto_path=api/userfilters \
		--go_out=pkg/pb \
		--go_opt=paths=source_relative \
        --go-grpc_out=pkg/pb \
		--go-grpc_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
        api/userfilters/*.proto

test:
	go test -v -cover ./...