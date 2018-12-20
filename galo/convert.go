package main

import (
	"log"
	"os"

	"github.com/aphecetche/galo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
			log.Fatalf("Could not get decoder for input file %v", src)
		}
		defer from.Close()
		to := NewClusterEncoder(output, dest)
		if to == nil {
			log.Fatalf("Could not get encoder for output file %v", dest)
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
