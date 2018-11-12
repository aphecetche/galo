package hist

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"go-hep.org/x/hep/hbook"
)

type objectBag struct {
	m map[string]hbook.Object
}

var (
	ErrCannotAddNilObject     = errors.New("Cannot add a nil object")
	ErrCannotAddUnnamedObject = errors.New("Cannot add an unnamed object")
)

func (ob *objectBag) NObjects() int {
	return len(ob.m)
}

func (ob *objectBag) Add(h hbook.Object) error {
	if h == nil {
		return ErrCannotAddNilObject
	}
	if len(h.Name()) == 0 {
		return ErrCannotAddUnnamedObject
	}
	if ob.m == nil {
		ob.m = make(map[string]hbook.Object)
	}
	ob.m[h.Name()] = h
	return nil
}

func (ob *objectBag) Get(hname string) (*hbook.Object, error) {
	val, ok := ob.m[hname]
	if !ok {
		return nil, fmt.Errorf("Could not get an object with name %s", hname)
	}
	return &val, nil
}

func (ob *objectBag) Print(out io.Writer) {
	for hname, _ := range ob.m {
		fmt.Fprintf(out, "%sOBJ:%s\n", strings.Repeat(" ", 8), hname)
	}
}
