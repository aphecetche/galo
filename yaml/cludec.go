package yaml

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
	yaml "gopkg.in/yaml.v2"
)

type yamlClusterDecoder struct {
	r io.Reader
}

func NewClusterDecoder(src io.Reader) *yamlClusterDecoder {
	return &yamlClusterDecoder{r: src}
}

var _ galo.DEClustersDecoder = (*yamlClusterDecoder)(nil)

// Decode reads the next YAML-encoded value of a DECluster from its input
// and stores it in the value pointed to by clu.
func (ya *yamlClusterDecoder) Decode(declu *galo.DEClusters) error {
	if declu == nil {
		return fmt.Errorf("Cannot decode into nil")
	}
	data, err := ioutil.ReadAll(ya.r)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no more data")
	}
	var ydeclu yaDEClusters
	err = yaml.Unmarshal([]byte(data), &ydeclu)
	if err != nil {
		return err
	}
	(*declu).DeID = mapping.DEID(ydeclu.DeID)
	seg := galo.SegCache.Segmentation((*declu).DeID)
	for _, yc := range ydeclu.Clusters {
		pre := preCluster(yc.Pre, seg)
		pos := clusterPos(yc.Pos)
		clu := galo.Cluster{Pre: pre, Pos: pos, Q: galo.ClusterCharge(yc.Charge)}
		(*declu).Clusters = append((*declu).Clusters, clu)
	}
	return nil
}

func (ya *yamlClusterDecoder) Close() {
}

func digit(yd yaDigit, finder mapping.PadByFEEFinder) galo.Digit {
	id, err := finder.FindPadByFEE(mapping.DualSampaID(yd.Dsid), mapping.DualSampaChannelID(yd.Dsch))
	if err != nil {
		panic(err)
	}
	return galo.Digit{
		ID: id,
		Q:  float64(yd.Charge),
	}
}

func preCluster(ypre yaPre, finder mapping.PadByFEEFinder) galo.PreCluster {
	var pre galo.PreCluster
	var dg galo.DigitGroup

	dg.RefTime = ypre.DigitGroup.RefTime

	for _, yd := range ypre.DigitGroup.Digits {
		dg.Digits = append(dg.Digits, digit(yd, finder))
	}
	pre.DigitGroup = dg
	return pre
}

func clusterPos(ypos yaPos) galo.ClusterPos {
	var pos galo.ClusterPos
	pos.X = float64(ypos.X)
	pos.Y = float64(ypos.Y)
	return pos
}

func rectangle(x, y, dx, dy float64) geo.Polygon {
	return geo.Polygon{
		{X: x - dx, Y: y - dy},
		{X: x + dx, Y: y - dy},
		{X: x + dx, Y: y + dy},
		{X: x - dx, Y: y + dy},
		{X: x - dx, Y: y - dy}}
}

func (step yaStep) getPixelPolygons() []geo.Polygon {
	var polygons []geo.Polygon
	pixels := step.Pixels
	for i := 0; i < len(pixels); i++ {
		pix := pixels[i]
		x := float64(pix.X)
		y := float64(pix.Y)
		dx := float64(pix.DX)
		dy := float64(pix.DY)
		polygons = append(polygons, rectangle(x, y, dx, dy))
	}
	return polygons
}

func (clu yaCluster) getPixelPolygons(stepNumber int) []geo.Polygon {
	return clu.Steps[stepNumber].getPixelPolygons()
}

func (clu yaCluster) getPadPolygons(bendingPlane bool) []geo.Polygon {
	return clu.Pre.getPadPolygons(bendingPlane)
}

func (clu yaCluster) getPadCharges(bendingPlane bool) []float32 {
	return clu.Pre.getPadCharges(bendingPlane)
}

func (pre yaPre) getPadCharges(bendingPlane bool) []float32 {
	var charges []float32
	for i := 0; i < len(pre.DigitGroup.Digits); i++ {
		digit := pre.DigitGroup.Digits[i]
		manuid := int(digit.Dsid)
		isBending := (manuid < 1024)
		if isBending != bendingPlane {
			continue
		}
		charges = append(charges, digit.Charge)
	}
	return charges
}

func (pre yaPre) getPadPolygons(bendingPlane bool) []geo.Polygon {
	var polygons []geo.Polygon
	for i := 0; i < len(pre.DigitGroup.Digits); i++ {
		digit := pre.DigitGroup.Digits[i]
		deid := digit.Deid
		manuid := mapping.DualSampaID(digit.Dsid)
		isBending := (manuid < 1024)
		if isBending != bendingPlane {
			continue
		}
		seg := galo.SegCache.CathodeSegmentation(mapping.DEID(deid), isBending)
		manuchannel := int(digit.Dsch)
		paduid, err := seg.FindPadByFEE(mapping.DualSampaID(manuid), mapping.DualSampaChannelID(manuchannel))
		if seg.IsValid(paduid) == false || err != nil {
			log.Fatalf("got invalid pad for DE %v MANU %v CH %v : %v -> paduid %v", deid, manuid, manuchannel, err, paduid)
		}
		x := seg.PadPositionX(paduid)
		y := seg.PadPositionY(paduid)
		dx := seg.PadSizeX(paduid) / 2
		dy := seg.PadSizeY(paduid) / 2
		polygons = append(polygons, rectangle(x, y, dx, dy))
	}
	return polygons
}
