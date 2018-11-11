package hist

import (
	"fmt"
	"io"
	"log"
	"path"
	"sort"
	"strings"

	"go-hep.org/x/hep/hbook"
)

type objectBag = map[string]hbook.Object

type Collection struct {
	Name string
	m    map[string]*objectBag
}

func NewCollection(n string) *Collection {
	return &Collection{
		Name: n,
		m:    make(map[string]*objectBag),
	}
}

func correctIdentifier(sid string) string {

	cid := sid
	if len(sid) > 0 {

		if !strings.HasSuffix(sid, "/") {
			cid = sid + "/"
		}

		if !strings.HasPrefix(sid, "/") {
			cid = "/" + sid
		}

		cid = strings.Replace(cid, "//", "/", -1)

	}
	return cid
}

func (hc *Collection) Add(identifier string, h hbook.Object) {
	sid := correctIdentifier(identifier)
	_, ok := hc.m[sid]
	if ok == false {
		hc.m[sid] = &objectBag{}
		(*hc.m[sid])[h.Name()] = h
	}

	(*hc.m[sid])[h.Name()] = h
}

func (hc *Collection) NObjects() int {
	n := 0
	for _, nb := range hc.m {
		n += len(*nb)
	}
	return n
}

func (hc *Collection) NKeys() int {
	return len(hc.m)
}

func (hc *Collection) SortAllIdentifiers() []string {
	ids := make([]string, 0, len(hc.m))
	for sid := range hc.m {
		ids = append(ids, sid)
	}
	sort.Strings(ids)
	return ids
}

func (hc *Collection) H1D(fullIdentifier string) *hbook.H1D {
	sid := correctIdentifier(FullIdToId(fullIdentifier))
	oname := FullIdToObjectName(fullIdentifier)
	nb := hc.m[sid]
	if nb == nil {
		log.Fatalf("Could not get objectBag for sid=%s\n", sid)
		return nil
	}
	h, _ := (*nb)[oname].(*hbook.H1D)
	return h
}

func FullIdToId(fullIdentifier string) string {
	return path.Dir(fullIdentifier)
}

func FullIdToObjectName(fullIdentifier string) string {
	return path.Base(fullIdentifier)
}

func (hc *Collection) Print(out io.Writer) {
	identifiers := hc.SortAllIdentifiers()
	fmt.Fprintf(out, "Number of identifiers %d\n", len(identifiers))
	for _, sid := range identifiers {
		fmt.Fprintf(out, "KEY %s\n", sid)
		nb := hc.m[sid]
		for _, obj := range *nb {
			fmt.Fprintf(out, "%sOBJ:%s\n", strings.Repeat(" ", 8), obj.Name())
		}
	}

}
