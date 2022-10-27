package main

import (
	"github.com/easyhutu/protoc-gen-meta/utils/generate"
)

func main() {
	gt := generate.New()
	gt.GenMeta().Done()
}