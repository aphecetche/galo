package cmd

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/aphecetche/galo/dataformats/run2"
	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var clustersCmd = &cobra.Command{
	Use:   "clusters [file to dump]",
	Short: "Dump clusters",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		dumpRun2Clusters(f)
	},
}

func dumpRun2Clusters(r io.Reader) {
	nevents := 0
	for nevents < 10 {
		sb := make([]byte, 4)
		nb, err := r.Read(sb)
		if nb != 4 {
			break
		}
		nevents++
		size := binary.LittleEndian.Uint32(sb)
		fmt.Println("size=", size)
		buf := make([]byte, size)
		nb, err = r.Read(buf)
		if uint32(nb) != size {
			panic(err)
		}
		event := run2.GetRootAsEvent(buf, 0)
		fmt.Println(event)
	}
}

func init() {
	dumpCmd.AddCommand(clustersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
