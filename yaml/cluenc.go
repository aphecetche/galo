package yaml

import (
	"io"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
	goyaml "gopkg.in/yaml.v2"
	yaml "gopkg.in/yaml.v2"
)

type PadFEELocatorFunc func(deid mapping.DEID) mapping.PadFEELocator

type yamlClusterEncoder struct {
	w          io.Writer
	yenc       *yaml.Encoder
	padlocfunc PadFEELocatorFunc
}

// NewClusterEncoder returns a new encoder that writes to dest.
// The returned encoder should be closed after use to flush
// of its data to dest.
func NewClusterEncoder(dest io.Writer, padlocfunc PadFEELocatorFunc) *yamlClusterEncoder {
	return &yamlClusterEncoder{w: dest,
		yenc:       goyaml.NewEncoder(dest),
		padlocfunc: padlocfunc}
}

func (enc *yamlClusterEncoder) Encode(clusters *galo.DEClusters) error {

	// Note we are not simlpy using something like :
	// buf, err := goyaml.Marshal(clusters) (which would technically
	// correctly generate yaml) because we want to have in the YAML file
	// a humand readable version of the digits,
	// with ID decoded, i.e. with dsid and dsch

	deid := clusters.DeID

	var yadeclus yaDEClusters

	yadeclus.DeID = int(deid)

	padloc := enc.padlocfunc(deid)

	for _, clu := range clusters.Clusters {
		pos := yaPos{
			X: float32(clu.Pos.X),
			Y: float32(clu.Pos.Y)}
		var pre yaPre
		for _, d := range clu.Pre.Digits {
			paduid := mapping.PadUID(d.ID)
			dsid := padloc.PadDualSampaID(paduid)
			dsch := padloc.PadDualSampaChannel(paduid)
			yd := yaDigit{
				Deid:   int(deid),
				Dsid:   int(dsid),
				Dsch:   int(dsch),
				Charge: float32(d.Q)}
			pre.DigitGroup.Digits = append(pre.DigitGroup.Digits, yd)
		}
		yaClu := yaCluster{Pos: pos, Charge: float32(clu.Q), Pre: pre}
		yadeclus.Clusters = append(yadeclus.Clusters, yaClu)
	}

	buf, err := goyaml.Marshal(yadeclus)

	_, err = enc.w.Write(buf)

	return err
}

func (enc *yamlClusterEncoder) Close() {
	enc.yenc.Close()
}

var _ galo.DEClustersEncoder = (*yamlClusterEncoder)(nil)
