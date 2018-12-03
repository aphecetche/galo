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

var (
	segcache mapping.SegCache
)

type yaDigit struct {
	Deid        int
	Manuid      int
	Manuchannel int
	Adc         int
	Charge      float32
}

type yaPixel struct {
	X  float32
	Y  float32
	DX float32
	DY float32
}

type yaStep struct {
	Pixels []yaPixel
	Ncalls int `yaml:"ncalls,omitempty`
}

type yaPre struct {
	Digits []yaDigit
}

type yaPos struct {
	X float32
	Y float32
	Z float32
}

type yaCluster struct {
	Pre    yaPre
	Pos    yaPos
	Charge float32
	Steps  []yaStep `yaml:"steps,omitempty`
}

type yamlClusterDecoder struct {
	r io.Reader
}

func Toto() {
}

func NewClusterDecoder(src io.Reader) *yamlClusterDecoder {
	return &yamlClusterDecoder{r: src}
}

func (ya *yamlClusterDecoder) Decode(clu *galo.Cluster) error {
	if clu == nil {
		return fmt.Errorf("Cannot decode into nil")
	}
	data, err := ioutil.ReadAll(ya.r)
	if err != nil {
		return err
	}
	var yaclu yaCluster
	err = yaml.Unmarshal([]byte(data), &yaclu)
	if err != nil {
		return err
	}
	(*clu).Pos.X = float64(yaclu.Pos.X)
	(*clu).Pos.Y = float64(yaclu.Pos.Y)
	(*clu).Pos.Z = float64(yaclu.Pos.Z)
	return nil
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
	for i := 0; i < len(pre.Digits); i++ {
		digit := pre.Digits[i]
		manuid := int(digit.Manuid)
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
	for i := 0; i < len(pre.Digits); i++ {
		digit := pre.Digits[i]
		deid := digit.Deid
		manuid := mapping.DualSampaID(digit.Manuid)
		isBending := (manuid < 1024)
		if isBending != bendingPlane {
			continue
		}
		seg := segcache.CathodeSegmentation(int(deid), isBending)
		manuchannel := int(digit.Manuchannel)
		paduid, err := seg.FindPadByFEE(manuid, manuchannel)
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
