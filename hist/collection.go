package hist

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"go-hep.org/x/hep/hbook"
)

type Collection struct {
	Name string
	bags map[string]*objectBag
}

func NewCollection(n string) *Collection {
	return &Collection{
		Name: n,
		bags: make(map[string]*objectBag),
	}
}

func correctPath(sid string) string {
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

func (hc *Collection) Add(path string, h hbook.Object) error {
	ob := hc.getObjectBagBy(path)
	if ob == nil {
		ob = &objectBag{}
		if hc.bags == nil {
			hc.bags = make(map[string]*objectBag)
		}
		hc.bags[path] = ob
	}
	return ob.Add(h)
}

func (hc *Collection) getObjectBagBy(path string) *objectBag {
	if hc.bags == nil {
		return nil
	}
	sid := correctPath(path)
	bag, ok := hc.bags[sid]
	if ok == false {
		return nil
	}
	return bag
}

func (hc *Collection) NObjects() int {
	n := 0
	for h, b := range hc.bags {
		fmt.Println(h)
		n += b.NObjects()
	}
	return n
}

func (hc *Collection) NKeys() int {
	return len(hc.bags)
}

func (hc *Collection) SortAllPaths() []string {
	ids := make([]string, 0, len(hc.bags))
	for sid := range hc.bags {
		ids = append(ids, sid)
	}
	sort.Strings(ids)
	return ids
}

func decodeFullPath(fullpath string) (string, string) {
	path := correctPath(FullIdToId(fullpath))
	hname := FullIdToObjectName(fullpath)
	return path, hname
}

func (hc *Collection) H1D(fullpath string) *hbook.H1D {
	path, hname := decodeFullPath(fullpath)
	ob := hc.getObjectBagBy(path)
	if ob == nil {
		fmt.Errorf("Could not get objectBag for path=%s\n", path)
		return nil
	}
	o, err := ob.Get(hname)
	if err != nil {
		fmt.Errorf("Could not get object %s\n", hname)
		return nil
	}
	h, ok := (*o).(*hbook.H1D)
	if !ok {
		fmt.Errorf("Object %s is not a H1D\n", hname)
		return nil
	}
	return h
}

func FullIdToId(fullpath string) string {
	return path.Dir(fullpath)
}

func FullIdToObjectName(fullpath string) string {
	return path.Base(fullpath)
}

func (hc *Collection) Print(out io.Writer) {
	paths := hc.SortAllPaths()
	fmt.Fprintf(out, "Number of paths %d\n", len(paths))
	for _, sid := range paths {
		fmt.Fprintf(out, "KEY %s\n", sid)
		b := hc.getObjectBagBy(sid)
		b.Print(out)
	}

}
