package gojieba

import (
	"github.com/sunshinas878/gojieba/deps/cppjieba"
	"github.com/sunshinas878/gojieba/deps/limonp"
	"github.com/sunshinas878/gojieba/dict"
)

func init() {
	dict.Init()
	limonp.Init()
	cppjieba.Init()
}
