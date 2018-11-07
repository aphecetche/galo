package run2

import "fmt"

func (e *Event) String() string {

	return fmt.Sprintf("BC %d Ntracklets %d isMB %v Nclusters %d", e.Bc(), e.Ntracklets(), e.IsMB(), e.ClustersLength())
}
