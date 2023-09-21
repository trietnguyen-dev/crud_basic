build:
	protoc -I . --go_out ./gen/go/ --go-grpc_out ./gen/go/ proto/*.proto
	ls ./gen/go/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'