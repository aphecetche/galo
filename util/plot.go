package util

import (
	"log"
	"path/filepath"
	"strings"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

func SavePlot(p *hplot.Plot, outputFileName string, name string) {
	fname := strings.TrimSuffix(outputFileName, filepath.Ext(outputFileName))
	fname += "_"
	fname += strings.Replace(name, "/", "_", -1) + ".pdf"
	fname = strings.Replace(fname, "__", "_", -1)
	err := p.Save(20*vg.Centimeter, -1, fname)
	if err != nil {
		log.Fatalf("Cannot save histogram:%s", err)
	}
}
