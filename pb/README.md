With `protoc` on the `PATH` and `protoc-gen-go` on the `PATH` (usually via `$GOPATH/bin`), from this dir run:

    protoc --go_out=. ast.proto token.proto type.proto