package run2

//
// import (
// 	"fmt"
// 	"log"
// 	"strings"
//
// 	"github.com/aphecetche/galo"
// 	"github.com/aphecetche/pigiron/mapping"
// )
//
// type WeightModel int
//
// const (
// 	NoWeight WeightModel = iota + 1
// 	RegularWeight
// 	SquaredWeight
// )
//
// // ClusterPosFunc computes the (x,y) position of a cluster
// type ClusterPosFunc struct {
// 	F    func(*Cluster) (float64, float64)
// 	Name string
// }
//
// type ClusterSelFunc struct {
// 	F    func(*EventClusters, int) bool
// 	Name string
// }
//
// func chargeClusters(ec *EventClusters, i int) bool {
// 	var clu Cluster
// 	var dig Digit
// 	ec.E.Clusters(&clu, i)
// 	pos := clu.Pos(nil)
// 	pre := clu.Pre(nil)
//
// 	bad := false
//
// 	n := pre.DigitsLength()
// 	if n < 50 {
// 		qTot := 0.0
// 		qCheck := 0.0
// 		for d := 0; d < n; d++ {
// 			pre.Digits(&dig, d)
// 			qCheck += computeDigitCharge(dig, pos.X(), pos.Y())
// 			qTot += float64(dig.Charge())
// 		}
// 		qCheck *= float64(clu.Charge())
// 		fmt.Println(strings.Repeat("-", 20))
// 		qcv := newValueCheck("Qclu", qTot, qCheck)
// 		fmt.Printf("%-6s Cluster %2d Mult %3d %v\n", ec.Label(i), i, n, qcv)
//
// 		for d := 0; d < n; d++ {
// 			pre.Digits(&dig, d)
// 			manuid := dig.Manuid()
// 			manuch := dig.Manuchannel()
// 			digCharge := float64(clu.Charge()) * computeDigitCharge(dig, pos.X(), pos.Y())
// 			qdv := newValueCheck("Qd", float64(dig.Charge()), digCharge)
// 			fmt.Printf("    Manu %4d ch %2d ADC %4d %v\n", manuid, manuch, dig.Adc(), qdv)
// 			if qdv.IsExtremelyBad() {
// 				bad = true
// 			}
// 		}
// 	}
//
// 	if bad {
// 		fmt.Println("would output cluster to html for visual inspection")
// 		cluster2SVG(ec, i, "inspectcluster", false)
// 	}
// 	return false
// }
//
// func strangeClusters(ec *EventClusters, i int) bool {
// 	silent := true
// 	var clu Cluster
// 	ec.E.Clusters(&clu, i)
// 	n := clu.Pre(nil).DigitsLength()
// 	if n > 200 {
// 		cluster2SVG(ec, i, "bigcluster", true)
// 		if !silent {
// 			fmt.Println("WARNING", "cluster with", n, "digits")
// 			DumpEventClusters(ec)
// 			fmt.Println("")
// 		}
// 	}
// 	return false
// }
//
// func allClusters(ec *EventClusters, i int) bool {
// 	return true
// }
//
// func simpleClusters(ec *EventClusters, i int) bool {
// 	return ec.IsSimple(i)
// }
//
// func splitClusters(ec *EventClusters, i int) bool {
// 	return ec.IsSplit(i)
// }
//
// func dupClusters(ec *EventClusters, i int) bool {
// 	return ec.IsDup(i)
// }
//
// func cogNoWeight(clu *Cluster) (float64, float64) {
// 	return cog(clu.Pre(nil), NoWeight)
// }
//
// func cogWithRegularWeight(clu *Cluster) (float64, float64) {
// 	return cog(clu.Pre(nil), RegularWeight)
// }
//
// func cogWithSquaredWeight(clu *Cluster) (float64, float64) {
// 	return cog(clu.Pre(nil), SquaredWeight)
// }
//
// // cog compute the center of gravity of the digits within precluster
// func cog(pre *PreCluster, weight WeightModel) (float64, float64) {
// 	var digit Digit
// 	var x, y, sumw float64
// 	for i := 0; i < pre.DigitsLength(); i++ {
// 		pre.Digits(&digit, i)
// 		deid := digit.Deid()
// 		manuid := mapping.DualSampaID(digit.Manuid())
// 		cseg := galo.SegCache.CathodeSegmentation(mapping.DEID(deid), manuid < 1024)
// 		manuchannel := int(digit.Manuchannel())
// 		paduid, err := cseg.FindPadByFEE(mapping.DualSampaID(manuid), mapping.DualSampaChannelID(manuchannel))
// 		if cseg.IsValid(paduid) == false || err != nil {
// 			log.Fatalf("got invalid pad for DE %v MANU %v CH %v : %v -> paduid %v", deid, manuid, manuchannel, err, paduid)
// 		}
// 		var w float64
// 		if weight == NoWeight {
// 			w = 1.0
// 		} else {
// 			w = float64(digit.Adc())
// 		}
// 		sumw += w
//
// 		x += cseg.PadPositionX(paduid) * w
// 		y += cseg.PadPositionY(paduid) * w
// 	}
// 	x /= sumw
// 	y /= sumw
// 	return x, y
// }
