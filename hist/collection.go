package hist

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"go-hep.org/x/hep/hbook"
)

type Collection struct {
	m map[string]hbook.Object
}

var (
	ErrCannotAddNilObject     = errors.New("Cannot add a nil object")
	ErrCannotAddUnnamedObject = errors.New("Cannot add an unnamed object")
)

func (ob *Collection) NObjects() int {
	return len(ob.m)
}

func (ob *Collection) Add(h hbook.Object) error {
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

func (ob *Collection) Get(hname string) (*hbook.Object, error) {
	val, ok := ob.m[hname]
	if !ok {
		return nil, fmt.Errorf("Could not get an object with name %s", hname)
	}
	return &val, nil
}

func (ob *Collection) H1D(hname string) (*hbook.H1D, error) {
	o, err := ob.Get(hname)
	if err != nil {
		return nil, err
	}
	h, ok := (*o).(*hbook.H1D)
	if ok {
		return h, nil
	}
	return nil, fmt.Errorf("Object %s is not a H1D", hname)
}

func (ob *Collection) H2D(hname string) (*hbook.H2D, error) {
	o, err := ob.Get(hname)
	if err != nil {
		return nil, err
	}
	h, ok := (*o).(*hbook.H2D)
	if ok {
		return h, nil
	}
	return nil, fmt.Errorf("Object %s is not a H2D", hname)
}

func (ob *Collection) Print(out io.Writer) {
	for hname, o := range ob.m {
		fmt.Fprintf(out, "%sNAME:%20s ", strings.Repeat(" ", 8), hname)
		switch v := o.(type) {
		case *hbook.H1D:
			fmt.Fprintf(out, "H1D NENTRIES %d", v.Entries())
		case *hbook.H2D:
			fmt.Fprintf(out, "H2D NENTRIES %d", v.Entries())
		default:
			fmt.Fprintf(out, "other stuf")
		}
		fmt.Fprintf(out, "\n")
	}
}