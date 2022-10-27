# protoc-gen-meta

一个生成 proto JSON化信息的插件

### build
```shell
go build .
```

### run the plugin
```shell
protoc --plugin protoc-gen-meta --meta_out=./ api.proto
```