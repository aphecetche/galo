package galo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type TaggedClusters struct {
	declu      *DEClusters
	isDup      []bool
	dupIndex   []int
	isSplit    []bool
	splitIndex []int
}

// TaggedClusterSelector selects or discard a cluster from a tagged group
// based on some criteria.
type TaggedClusterSelector interface {
	// Select decides if a cluster is to be kept or not.
	Select(tagclu *TaggedClusters, i int) bool
	// Name of the selector.
	Name() string
}

func NewClusterSelector(selector string) TaggedClusterSelector {
	name := strings.ToUpper(selector)
	if name == "ANY" || name == "ALL" {
		return TaggedClusterSelectorTrue{}
	}
	if name == "SPLIT" {
		return TaggedClusterSelectorSplit{}
	}
	if strings.HasPrefix(name, "MULT") {
		p := strings.Split(name, ":")
		if len(p) == 2 && p[0] == "MULT" {
			n, err := strconv.Atoi(p[1])
			if err == nil {
				return TaggedClusterSelectorMult{mult: n}
			}
		}
	}
	return nil
}

func (tc *TaggedClusters) Clusters() []Cluster {
	return tc.declu.Clusters
}

func getSelected(declu *DEClusters, tcsel TaggedClusterSelector) (*TaggedClusters, []int) {
	tc := GetTaggedClusters(declu)
	var selected []int
	for i, _ := range tc.Clusters() {
		if tcsel.Select(tc, i) {
			selected = append(selected, i)
		}
	}
	return tc, selected
}

func DumpClusters(index int, tc *TaggedClusters, selected []int) {
	printHeader(index)
	tc.PrintSelected(os.Stdout, selected)
}

type TaggedClusterSelectorTrue struct{}

func (sel TaggedClusterSelectorTrue) Select(tc *TaggedClusters, i int) bool {
	return true
}

func (sel TaggedClusterSelectorTrue) Name() string {
	return "ANY"
}

type TaggedClusterSelectorSplit struct{}

func (sel TaggedClusterSelectorSplit) Select(tc *TaggedClusters, i int) bool {
	return tc.IsSplit(i)
}

func (sel TaggedClusterSelectorSplit) Name() string {
	return "SPLIT"
}

type TaggedClusterSelectorMult struct {
	mult int
}

func (sel TaggedClusterSelectorMult) Select(tc *TaggedClusters, i int) bool {
	return len(tc.Clusters()[i].Pre.Digits) > sel.mult
}

func (sel TaggedClusterSelectorMult) Name() string {
	return "MULT:" + strconv.Itoa(sel.mult)
}

func GetTaggedClusters(declu *DEClusters) *TaggedClusters {
	nclu := len(declu.Clusters)
	tc := TaggedClusters{
		declu:      declu,
		isDup:      make([]bool, nclu),
		dupIndex:   make([]int, nclu),
		isSplit:    make([]bool, nclu),
		splitIndex: make([]int, nclu)}

	for i := 0; i < nclu; i++ {
		ci := declu.Clusters[i]
		for j := i + 1; j < nclu; j++ {
			cj := declu.Clusters[j]
			if SameCluster(ci, cj) {
				tc.isDup[i] = true
				tc.isDup[j] = true
				tc.dupIndex[i] = i
				tc.dupIndex[j] = i
			} else {
				pi := ci.Pre
				pj := cj.Pre
				if SamePreCluster(pi, pj) {
					tc.isSplit[i] = true
					tc.isSplit[j] = true
					tc.splitIndex[i] = i
					tc.splitIndex[j] = i
				}
			}
		}
	}
	return &tc
}

func (tc *TaggedClusters) IsSimple(i int) bool {
	return tc.IsDup(i) == false && tc.IsSplit(i) == false
}

func (tc *TaggedClusters) IsSplit(i int) bool {
	return tc.isSplit[i]
}

func (tc *TaggedClusters) IsDup(i int) bool {
	return tc.isDup[i]
}

func (tc *TaggedClusters) Label(i int) string {
	if tc.isDup[i] {
		return "D" + strconv.Itoa(tc.dupIndex[i])
	}

	if tc.isSplit[i] {
		return "S" + strconv.Itoa(tc.splitIndex[i])
	}
	return "N"
}

func (tc *TaggedClusters) NDuplicates() int {
	c := make(map[int]bool)
	for i := range tc.isDup {
		if tc.isDup[i] {
			c[tc.dupIndex[i]] = true
		}
	}
	return len(c)
}

func (tc *TaggedClusters) PrintSelected(w io.Writer, sel []int) {
	s := ""
	padloc := SegCache.Segmentation(tc.declu.DeID)
	for _, i := range sel {
		clu := tc.declu.Clusters[i]
		s += fmt.Sprintf("%6s Q %7.2f POS %v", tc.Label(i), clu.Q, clu.Pos)
		s += fmt.Sprintf(" PREQ: %7.2f %4d digits\n", clu.Pre.Charge(), len(clu.Pre.Digits))
		for _, d := range clu.Pre.Digits {
			dsid := padloc.PadDualSampaID(d.ID)
			dsch := padloc.PadDualSampaChannel(d.ID)
			s += fmt.Sprintf("%sQ %7.2f ID %6d DS %4d CH %2d\n", strings.Repeat(" ", 10),
				d.Q, d.ID, dsid, dsch)
		}
	}
	w.Write([]byte(s))
}

func (tc *TaggedClusters) String() string {
	buf := new(bytes.Buffer)
	var all []int
	for i, _ := range tc.Clusters() {
		all = append(all, i)
	}
	tc.PrintSelected(buf, all)
	return buf.String()
}

func printHeader(nevents int) {
	fmt.Printf("Event %6d\n", nevents)
}
