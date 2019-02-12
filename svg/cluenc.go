package svg

import (
	"fmt"
	"io"
	"strconv"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

type svgClusterEncoder struct {
	w      io.Writer
	svgw   geo.SVGWriter
	html   bool
	origin bool
	charge bool
}

// NewClusterEncoder returns a new encoder that writes to dest.
// The returned encoder should be closed after use to flush
// of its data to dest.
func NewClusterEncoder(dest io.Writer) *svgClusterEncoder {
	return &svgClusterEncoder{w: dest, html: true, charge: false}
}

// Encode writes the encoding of clu to the stream.
func (enc *svgClusterEncoder) Encode(clusters *galo.DEClusters) error {
	seg := mapping.NewSegmentation(clusters.DeID)
	if seg == nil {
		return fmt.Errorf("Could not create segmentation for DE=%d", clusters.DeID)
	}
	isBending := true
	for _, clu := range clusters.Clusters {
		enc.encodeCathode(&clu, seg, isBending)
		enc.encodeCathode(&clu, seg, !isBending)
		enc.svgw.GroupStart("position")
		enc.svgw.Circle(clu.Pos.X, clu.Pos.Y, 0.02)
		enc.svgw.GroupEnd()
	}
	return nil
}

func (enc *svgClusterEncoder) SVGOnly() {
	enc.html = false
}

func (enc *svgClusterEncoder) MoveToOrigin() {
	enc.origin = true
}

func (enc *svgClusterEncoder) WithCharge() {
	enc.charge = true
}

func (enc *svgClusterEncoder) DefaultStyle() {
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
.position {
        fill: green;
        stroke: none;
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
	enc.svgw.Style(cssStyle)
}

func (enc *svgClusterEncoder) Style(style string) {
	enc.svgw.Style(style)
}

func (enc *svgClusterEncoder) Close() {
	if enc.origin {
		enc.svgw.MoveToOrigin()
	}
	if enc.html == true {
		enc.svgw.WriteHTML(enc.w)
	} else {
		enc.svgw.WriteSVG(enc.w)
	}
}

func (enc *svgClusterEncoder) encodeCathode(clu *galo.Cluster, seg mapping.Segmentation, isBending bool) {
	if isBending {
		enc.svgw.GroupStart("bending-pads")
	} else {
		enc.svgw.GroupStart("non-bending-pads")
	}
	for _, d := range clu.Pre.Digits {
		paduid := mapping.PadUID(d.ID)
		if seg.IsBendingPad(paduid) != isBending {
			continue
		}
		var xmin, ymin, xmax, ymax float64
		mapping.ComputePadBBox(seg, paduid, &xmin, &ymin, &xmax, &ymax)
		if enc.charge == false {
			enc.svgw.Rect(xmin, ymin, xmax-xmin, ymax-ymin)
		} else {
			enc.svgw.RectWithClass(xmin, ymin, xmax-xmin, ymax-ymin, getClass(float32(d.Q)))
		}
	}
	enc.svgw.GroupEnd()
}

func getClass(q float32) string {
	color := []float32{0, 0.10, 0.20, 0.30, 0.40, 0.60}
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

var _ galo.DEClustersEncoder = (*svgClusterEncoder)(nil)

// TODO: get back the pixels into the game
//
// func convertPixelsToSVG(svg *geo.SVGWriter, pixels []geo.Polygon) {
// 	svg.GroupStart("pixels")
// 	for _, p := range pixels {
// 		svg.Polygon(&p)
// 	}
// 	svg.GroupEnd()
// }
