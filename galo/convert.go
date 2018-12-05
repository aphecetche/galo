package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/svg"
	"github.com/aphecetche/galo/yaml"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert cluster(s) from YAML to SVG",
	Run: func(cmd *cobra.Command, args []string) {
		from, err := os.Open(src)
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()
		if len(dest) == 0 {
			fmt.Println("src", src)
			dest = strings.Replace(src, path.Ext(src), ".html", -1)
			fmt.Println("dest", dest)
		}
		to, err := os.Create(dest)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()
		input := yaml.NewClusterDecoder(from)
		output := svg.NewClusterEncoder(to)
		var cluster galo.Cluster
		for {
			err := input.Decode(&cluster)
			if err != nil {
				log.Fatal(err.Error())
				break
			}
			fmt.Println("decoded cluster")
			err = output.Encode(&cluster)
			if err != nil {
				log.Fatal(err.Error())
				break
			}
			fmt.Println("encoded cluster")
		}
	},
}

var src string
var dest string

func init() {
	clusterCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&src, "from", "f", "", "Source file")
	convertCmd.Flags().StringVarP(&dest, "to", "t", "", "Destination file")
}
