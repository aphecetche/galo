package run2

import (
	"log"
	"os"
	"strconv"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/segcontour"
)

var cluster2SVGindex int

func cluster2SVG(ec *EventClusters, i int) {

	var clu Cluster
	ec.E.Clusters(&clu, i)

	var pads, bpads, nbpads []geo.Polygon

	pre := clu.Pre(nil)

	var digit Digit
	for i := 0; i < pre.DigitsLength(); i++ {
		pre.Digits(&digit, i)
		deid := digit.Deid()
		manuid := int(digit.Manuid())
		isBending := (manuid < 1024)
		seg := segcache.Segmentation(int(deid), isBending)
		manuchannel := int(digit.Manuchannel())
		paduid, err := seg.FindPadByFEE(manuid, manuchannel)
		if seg.IsValid(paduid) == false || err != nil {
			log.Fatalf("got invalid pad for DE %v MANU %v CH %v : %v -> paduid %v", deid, manuid, manuchannel, err, paduid)
		}
		x := seg.PadPositionX(paduid)
		y := seg.PadPositionY(paduid)
		dx := seg.PadSizeX(paduid) / 2
		dy := seg.PadSizeY(paduid) / 2
		p := geo.Polygon{
			{X: x - dx, Y: y - dy},
			{X: x + dx, Y: y - dy},
			{X: x + dx, Y: y + dy},
			{X: x - dx, Y: y + dy},
			{X: x - dx, Y: y - dy}}
		pads = append(pads, p)
		if isBending {
			bpads = append(bpads, p)
		} else {
			nbpads = append(nbpads, p)
		}
	}

	c, err := geo.NewContour(pads)
	if err != nil {
		panic(err)
	}

	seg := segcache.Segmentation(100, true)
	segcontour := segcontour.Contour(seg)

	svg := geo.SVGWriter{Width: 1024, BBox: segcontour.BBox()}

	svg.GroupStart("de")
	svg.Contour(&segcontour)
	svg.GroupEnd()

	svg.GroupStart("non-bending-pads")
	for _, p := range nbpads {
		svg.Polygon(&p)
	}
	svg.GroupEnd()

	svg.GroupStart("bending-pads")
	for _, p := range bpads {
		svg.Polygon(&p)
	}
	svg.GroupEnd()

	svg.GroupStart("clusters")
	svg.Contour(&c)
	svg.GroupEnd()

	svg.GroupStart("clupos")
	pos := clu.Pos(nil)
	x := float64(pos.X())
	y := float64(pos.Y())
	svg.Circle(x, y, 0.025)
	svg.GroupEnd()
	svg.Style(`
.clusters {
  stroke: rgba(250,224,159,0.5);
  stroke-width: 0.075px;
  fill: none;
}
.de {
  stroke: black;
  stroke-width: 0.1px;
  fill: none;
}
.bending-pads {
  stroke: black;
  stroke-width: 0.025px;
  fill:none;
}
.non-bending-pads {
  stroke: rgba(200,200,200,1.0);
  stroke-width: 0.0125px;
  fill:rgba(200,200,200,0.25);
`)

	out, err := os.Create("toto" + strconv.Itoa(cluster2SVGindex) + ".html")
	cluster2SVGindex++
	if err != nil {
		panic(err)
	}

	svg.WriteHTML(out)
}
