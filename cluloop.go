package galo

func ClusterLoop(dec DEClustersDecoder, tcsel TaggedClusterSelector, firstEvent int, maxEvents int, worker func(index int, tc *TaggedClusters, selected []int)) (int, int) {
	nevents := 0
	nsel := 0
	for {
		var declu DEClusters
		err := dec.Decode(&declu)
		if err != nil {
			break
		}
		if nevents >= firstEvent {
			tc, selected := getSelected(&declu, tcsel)
			if len(selected) > 0 {
				nsel++
				worker(nevents, tc, selected)
			}
		}
		nevents++
		if nevents > maxEvents {
			break
		}
	}
	return nevents, nsel
}
