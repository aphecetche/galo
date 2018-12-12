package main

import (
	"io"
	"path"
	"strings"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/run2"
	"github.com/aphecetche/galo/svg"
	"github.com/aphecetche/galo/yaml"
	"github.com/aphecetche/pigiron/mapping"
)

func isYAML(filename string) bool {
	return strings.ToLower(path.Ext(filename)) == ".yaml"
}

func isSVG(filename string) bool {
	return strings.ToLower(path.Ext(filename)) == ".svg"
}

func isHMTL(filename string) bool {
	return strings.ToLower(path.Ext(filename)) == ".html"
}

func isRun2(filename string) bool {
	//FIXME: that's a very poor way of checking the file
	//is actually containing flatbuffers ...
	return strings.ToLower(path.Ext(filename)) == ".dat"
}

func NewClusterDecoder(r io.Reader, filename string) galo.DEClustersDecoder {
	if isYAML(filename) {
		return yaml.NewClusterDecoder(r)
	}
	if isRun2(filename) {
		return run2.NewClusterDecoder(r.(io.ReaderAt),
			func(deid mapping.DEID) mapping.PadByFEEFinder {
				return galo.SegCache.Segmentation(deid)
			}, 0)
	}
	return nil
}

func NewClusterEncoder(w io.Writer, filename string) galo.DEClustersEncoder {
	if isSVG(filename) || isHMTL(filename) {
		output := svg.NewClusterEncoder(w)
		output.DefaultStyle()
		output.MoveToOrigin()
		if isSVG(filename) {
			output.SVGOnly()
		}
		return output
	}
	if isYAML(filename) {
		return yaml.NewClusterEncoder(w, func(deid mapping.DEID) mapping.PadFEELocator {
			return galo.SegCache.Segmentation(deid)
		})
	}
	return nil
}
