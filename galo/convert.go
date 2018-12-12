package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/run2"
	"github.com/aphecetche/galo/svg"
	"github.com/aphecetche/galo/yaml"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert cluster(s) from one format to another",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := os.Open(src)
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()
		output, err := os.Create(dest)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()
		from := NewClusterDecoder(input, src)
		if from == nil {
			log.Fatalf("Could not get decoder for input file %v", input)
		}
		defer from.Close()
		to := NewClusterEncoder(output, dest)
		if to == nil {
			log.Fatalf("Could not get encoder for output file %v", input)
		}
		defer to.Close()
		n := 0
		for {
			err = convert(from, to)
			if err != nil {
				break
			}
			n++
			if n > maxEvents {
				break
			}
		}
	},
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

func convert(from galo.DEClustersDecoder, to galo.DEClustersEncoder) error {
	var clusters galo.DEClusters
	for {
		err := from.Decode(&clusters)
		if err != nil {
			return errors.Wrap(err, "Decoding error")
		}
		if len(clusters.Clusters) == 0 {
			continue
		}
		err = to.Encode(&clusters)
		if err != nil {
			return errors.Wrap(err, "Encoding error")
		}
		return nil
	}
}

var src string
var dest string

func init() {
	clusterCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&src, "input", "i", "", "Source file")
	convertCmd.Flags().StringVarP(&dest, "output", "o", "", "Destination file")
	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")
}
