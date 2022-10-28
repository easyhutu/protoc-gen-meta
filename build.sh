GOOS=linux GOARCH=amd64 go build -o ./dist/linux/protoc-gen-meta
GOOS=windows GOARCH=amd64 go build -o ./dist/windows/protoc-gen-meta.exe
GOOS=darwin GOARCH=amd64 go build -o ./dist/darwin/protoc-gen-meta