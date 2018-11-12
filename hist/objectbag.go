package hist

import (
	"errors"
	"fmt"
	"io"

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
	ob.m[h.Name()] = h
	return nil
}

func (ob *objectBag) Get(hname string) (*hbook.Object, error) {
	val, err := (*ob)[hname]
	if err != nil {
		return Error.Newf("Could not get an object with name %s", hname)
	}
	return val
}

func (ob *objectBag) Print(out io.Writer) {
	for _, hname := range ob.m {
		fmt.Fprintf(out, "%s\n", hname)
	}
}
