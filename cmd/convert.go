package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aphecetche/galo/convert"
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
		if len(dest) == 0 {
			fmt.Println("src", src)
			dest = strings.Replace(src, path.Ext(src), ".html", -1)
			fmt.Println("dest", dest)
		}
		to, err := os.Create(dest)
		if err != nil {
			log.Fatal(err)
		}
		convert.Cluster(from, to)
	},
}

var src string
var dest string

func init() {
	clusterCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&src, "from", "f", "", "Source file")
	convertCmd.Flags().StringVarP(&dest, "to", "t", "", "Destination file")
}
