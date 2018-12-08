package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/svg"
	"github.com/aphecetche/galo/yaml"
	"github.com/aphecetche/pigiron/mapping"
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
	return false
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
		_ = convert(from, to)
		// TODO: handle err (in particular EOF which is not really an
		// error but the happy ending...
	},
}

func NewClusterDecoder(r io.Reader, filename string) galo.DEClustersDecoder {
	if isYAML(filename) {
		return yaml.NewClusterDecoder(r)
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
			return err
		}
		err = to.Encode(&clusters)
		if err != nil {
			return err
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
