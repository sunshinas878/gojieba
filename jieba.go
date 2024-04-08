//go:generate go-bindata -o=./dict/bindata.go -pkg=dict ./dict/
package gojieba

/*
#cgo CXXFLAGS: -I./deps -DLOGGING_LEVEL=LL_WARNING -O3 -Wall  -Wno-deprecated
#include <stdlib.h>
#include "jieba.h"
*/
import "C"
import (
	"os"
	"runtime"
	"unsafe"

	"github.com/sunshinas878/gojieba/dict"
)

type TokenizeMode int

const (
	DefaultMode TokenizeMode = iota
	SearchMode
)

type Word struct {
	Str   string
	Start int
	End   int
}

type Jieba struct {
	jieba C.Jieba
}

func NewJieba(paths ...string) *Jieba {
	dictpaths := getDictPaths(paths...)
	dpath, hpath, upath, ipath, spath := C.CString(dictpaths[0]), C.CString(dictpaths[1]), C.CString(dictpaths[2]), C.CString(dictpaths[3]), C.CString(dictpaths[4])
	defer C.free(unsafe.Pointer(dpath))
	defer C.free(unsafe.Pointer(hpath))
	defer C.free(unsafe.Pointer(upath))
	defer C.free(unsafe.Pointer(ipath))
	defer C.free(unsafe.Pointer(spath))
	jieba := &Jieba{
		C.NewJieba(
			dpath,
			hpath,
			upath,
			ipath,
			spath,
		),
	}
	runtime.SetFinalizer(jieba, (*Jieba).Free)
	return jieba
}

func NewJiebaWithContent(userDictPath string) *Jieba {

	dpath := C.CString(dict.MustAssetString("dict/jieba.dict.utf8"))
	hpath := C.CString(dict.MustAssetString("dict/hmm_model.utf8"))
	upath := C.CString("")
	if len(userDictPath) == 0 {
		upath = C.CString(dict.MustAssetString("dict/user.dict.utf8"))
	} else {
		upath = C.CString(loadFile2String(userDictPath))
	}
	ipath := C.CString(dict.MustAssetString("dict/idf.utf8"))
	spath := C.CString(dict.MustAssetString("dict/stop_words.utf8"))

	defer C.free(unsafe.Pointer(dpath))
	defer C.free(unsafe.Pointer(hpath))
	defer C.free(unsafe.Pointer(upath))
	defer C.free(unsafe.Pointer(ipath))
	defer C.free(unsafe.Pointer(spath))
	jieba := &Jieba{
		C.NewJiebaWithContent(
			dpath,
			hpath,
			upath,
			ipath,
			spath,
		),
	}
	runtime.SetFinalizer(jieba, (*Jieba).Free)
	return jieba
}

func (x *Jieba) Free() {
	C.FreeJieba(x.jieba)
}

func (x *Jieba) Cut(s string, hmm bool) []string {
	c_int_hmm := 0
	if hmm {
		c_int_hmm = 1
	}
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words **C.char = C.Cut(x.jieba, cstr, C.int(c_int_hmm))
	defer C.FreeWords(words)
	res := cstrings(words)
	return res
}

func (x *Jieba) CutAll(s string) []string {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words **C.char = C.CutAll(x.jieba, cstr)
	defer C.FreeWords(words)
	res := cstrings(words)
	return res
}

func (x *Jieba) CutForSearch(s string, hmm bool) []string {
	c_int_hmm := 0
	if hmm {
		c_int_hmm = 1
	}
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words **C.char = C.CutForSearch(x.jieba, cstr, C.int(c_int_hmm))
	defer C.FreeWords(words)
	res := cstrings(words)
	return res
}

func (x *Jieba) Tag(s string) []string {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words **C.char = C.Tag(x.jieba, cstr)
	defer C.FreeWords(words)
	res := cstrings(words)
	return res
}

func (x *Jieba) AddWord(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.AddWord(x.jieba, cstr)
}

func (x *Jieba) AddWordEx(s string, freq int, tag string) {
	cstr := C.CString(s)
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	defer C.free(unsafe.Pointer(cstr))
	C.AddWordEx(x.jieba, cstr, C.int(freq), ctag)
}

func (x *Jieba) RemoveWord(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.RemoveWord(x.jieba, cstr)
}

func (x *Jieba) Tokenize(s string, mode TokenizeMode, hmm bool) []Word {
	c_int_hmm := 0
	if hmm {
		c_int_hmm = 1
	}
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words *C.Word = C.Tokenize(x.jieba, cstr, C.TokenizeMode(mode), C.int(c_int_hmm))
	defer C.free(unsafe.Pointer(words))
	return convertWords(s, words)
}

type WordWeight struct {
	Word   string
	Weight float64
}

func (x *Jieba) Extract(s string, topk int) []string {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var words **C.char = C.Extract(x.jieba, cstr, C.int(topk))
	res := cstrings(words)
	defer C.FreeWords(words)
	return res
}

func (x *Jieba) ExtractWithWeight(s string, topk int) []WordWeight {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	words := C.ExtractWithWeight(x.jieba, cstr, C.int(topk))
	p := unsafe.Pointer(words)
	res := cwordweights((*C.struct_CWordWeight)(p))
	defer C.FreeWordWeights(words)
	return res
}

func cwordweights(x *C.struct_CWordWeight) []WordWeight {
	var s []WordWeight
	eltSize := unsafe.Sizeof(*x)
	for (*x).word != nil {
		ww := WordWeight{
			C.GoString(((C.struct_CWordWeight)(*x)).word),
			float64((*x).weight),
		}
		s = append(s, ww)
		x = (*C.struct_CWordWeight)(unsafe.Pointer(uintptr(unsafe.Pointer(x)) + eltSize))
	}
	return s
}

func loadFile(filepath string) []byte {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}
	return buf
}

func loadFile2String(filepath string) string {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return ""
	}
	return string(buf)
}

func loadByBytes(filepath string) {
	buf := loadFile(filepath)
	cstr := C.CString(string(buf))
	defer C.free(unsafe.Pointer(cstr))
	C.printStr(cstr)
}
