package run2

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type EventClusters struct {
	E        *Event
	isDup    []bool
	dupIndex []int
	isSplit  []bool
}

func GetEventClusters(e *Event) *EventClusters {
	ec := EventClusters{
		E:        e,
		isDup:    make([]bool, e.ClustersLength()),
		dupIndex: make([]int, e.ClustersLength()),
		isSplit:  make([]bool, e.ClustersLength())}

	var ci, cj Cluster
	for i := 0; i < e.ClustersLength(); i++ {
		e.Clusters(&ci, i)
		for j := i + 1; j < e.ClustersLength(); j++ {
			e.Clusters(&cj, j)
			if SameCluster(ci, cj) {
				ec.isDup[i] = true
				ec.isDup[j] = true
				ec.dupIndex[i] = i
				ec.dupIndex[j] = i
			}
		}
	}
	return &ec
}

func (ec *EventClusters) IsSimple(i int) bool {
	return ec.IsDup(i) == false && ec.IsSplit(i) == false
}

func (ec *EventClusters) IsSplit(i int) bool {
	return ec.isSplit[i]
}

func (ec *EventClusters) IsDup(i int) bool {
	return ec.isDup[i]
}

func (ec *EventClusters) Label(i int) string {
	if ec.isDup[i] {
		return "D" + strconv.Itoa(ec.dupIndex[i])
	}

	if ec.isSplit[i] {
		return "S"
	}
	return "N"
}

func (ec *EventClusters) NDuplicates() int {
	c := make(map[int]bool)
	for i := range ec.isDup {
		if ec.isDup[i] {
			c[ec.dupIndex[i]] = true
		}
	}
	return len(c)
}

func (ec *EventClusters) dumpHeader() {
	fmt.Printf("BC %d Ntracklets %d isMB %v Nclusters %d", ec.E.Bc(), ec.E.Ntracklets(), ec.E.IsMB(), ec.E.ClustersLength())

	nd := ec.NDuplicates()
	if nd > 0 {
		fmt.Printf(" (%d duplicates)", nd)
	}
	fmt.Printf("\n")
}

func DumpEventClusters(ec *EventClusters) {

	ec.dumpHeader()

	var clu Cluster
	var digit Digit

	for i := 0; i < ec.E.ClustersLength(); i++ {
		b := ec.E.Clusters(&clu, i)
		if b == false {
			log.Fatalf("could not get cluster %d", i)
		}

		pos := clu.Pos(nil)
		fmt.Printf("%6s X %7.4f Y %7.4f", ec.Label(i), pos.X(), pos.Y())

		pre := clu.Pre(nil)
		fmt.Printf("%4d digits [DE,MANU,CH]:", pre.DigitsLength())
		n := 0
		for id := 0; id < pre.DigitsLength(); id++ {
			if n == 0 {
				fmt.Printf("\n%s", strings.Repeat(" ", 15))
			}
			bd := pre.Digits(&digit, id)
			if bd == false {
				log.Fatalf("could not get digit %d", i)
			}
			fmt.Printf(" {%4d,%4d,%2d}", digit.Deid(), digit.Manuid(), digit.Manuchannel())
			n++
			if n == 5 {
				n = 0
			}
		}
		fmt.Println("")
	}
}
