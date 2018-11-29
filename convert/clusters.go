package convert

import (
	"io"
	"log"
	"strconv"

	"github.com/aphecetche/pigiron/geo"
)

var (
	cssStyle string = `
.cluster {
        stroke: rgba(250,224,159,0.5);
        stroke-width: 0.075px;
        fill: none;
}
.bending-pads {
        stroke: black;
        stroke-width: 0.0150px;
        fill:none;
}
.non-bending-pads {
        stroke: rgba(200,200,200,1.0);
        stroke-width: 0.0125px;
        fill:rgba(200,200,200,0.25);
}
.pixels {
        stroke: red;
        stroke-width: 0.005px;
        fill:rgba(255,255,255,0.0)
}

/* http://colorbrewer2.org/?type=sequential&scheme=YlGnBu&n=5 */
.q0-5{fill:rgba(255,255,204,0.5)} 
.q1-5{fill:rgba(161,218,180,0.5)} 
.q2-5{fill:rgba(65,182,196,0.5)} 
.q3-5{fill:rgba(44,127,184,0.5)} 
.q4-5{fill:rgba(37,52,148,0.5)}

.overflow{fill:black}
`
)

func getClass(q float32) string {
	color := []float32{0, 10, 20, 30, 40, 60}
	index := len(color) + 1
	for i := 0; i < len(color)-1; i++ {
		if q >= color[i] && q < color[i+1] {
			return "q" + strconv.Itoa(i) + "-5"
		}
	}
	if index >= len(color) {
		return "overflow"
	}
	return "wrong"
}

func convertPadsToSVG(svg *geo.SVGWriter, bending bool, pads []geo.Polygon, charges []float32) {
	if charges != nil && len(charges) != len(pads) {
		log.Fatalf("Got %d charges for %d pads", len(charges), len(pads))
	}
	groupName := "bending-pads"
	if bending == false {
		groupName = "non-bending-pads"
	}
	svg.GroupStart(groupName)
	if charges == nil {
		for _, p := range pads {
			svg.Polygon(&p)
		}
	} else {
		for i, p := range pads {
			svg.PolygonWithClass(&p, getClass(charges[i]))
		}
	}
	svg.GroupEnd()
}

func convertPixelsToSVG(svg *geo.SVGWriter, pixels []geo.Polygon) {
	svg.GroupStart("pixels")
	for _, p := range pixels {
		svg.Polygon(&p)
	}
	svg.GroupEnd()
}

// func convertToSVG(svg *geo.SVGWriter, bpads, nbpads []geo.Polygon, pixels []geo.Polygon) {
//
// 	var pads []geo.Polygon
// 	pads = append(pads, bpads...)
// 	pads = append(pads, nbpads...)
// 	if len(pads) > 0 {
// 		svg.GroupStart("cluster")
// 		c, err := geo.NewContour(pads)
// 		if err != nil {
// 			panic(err)
// 		}
// 		svg.Contour(&c)
// 		svg.GroupEnd()
// 	}
//
// 	svg.GroupStart("pixels")
// 	for _, p := range pixels {
// 		svg.Polygon(&p)
// 	}
// 	svg.GroupEnd()
// }

func Cluster(src io.Reader, dest io.Writer) {

	clu, err := newCluster(src)
	if err != nil {
		log.Fatal(err)
	}

	var pads []geo.Polygon

	bpads := clu.Pre.getPadPolygons(true)
	nbpads := clu.Pre.getPadPolygons(false)
	pads = append(pads, bpads...)
	pads = append(pads, nbpads...)

	c, err := geo.NewContour(pads)
	if err != nil {
		panic(err)
	}

	b := c.BBox()
	box, _ := geo.NewBBox(b.Xmin(), b.Ymin(), b.Xmax()+b.Width()*1.5, b.Ymax())

	svg := geo.NewSVGWriter(1000, box, true)

	const bending bool = true
	// show first non-bending and then bending so bending SVG objects
	// are on top
	convertPadsToSVG(svg, !bending, nbpads, nil)
	convertPadsToSVG(svg, bending, bpads, nil)

	pixels := clu.getPixelPolygons(0)
	convertPixelsToSVG(svg, pixels)

	xshift := b.Width() * 1.1
	yshift := 0.0

	tbpads := geo.Translate(bpads, xshift, yshift)
	tnbpads := geo.Translate(nbpads, xshift, yshift)

	convertPadsToSVG(svg, !bending, tnbpads, clu.getPadCharges(!bending))
	convertPadsToSVG(svg, bending, tbpads, clu.getPadCharges(bending))

	tpixels := geo.Translate(pixels, xshift, yshift)
	convertPixelsToSVG(svg, tpixels)

	svg.Style(cssStyle)

	svg.WriteHTML(dest)

}
