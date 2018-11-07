// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aphecetche/galo/dataformats/run2"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump clusters",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		nevents := run2.ForEachEvent(f, dumpEvent, maxEvents)
		fmt.Println(nevents, " events processed")
	},
}

func dumpEvent(e *run2.Event) {

	fmt.Printf("BC %d Ntracklets %d isMB %v Nclusters %d\n", e.Bc(), e.Ntracklets(), e.IsMB(), e.ClustersLength())

	var clu run2.Cluster
	for i := 0; i < e.ClustersLength(); i++ {
		b := e.Clusters(&clu, i)
		if b == false {
			log.Fatalf("could not get cluster %d", i)
		}

		pos := clu.Pos(nil)
		fmt.Printf("%s X %7.2f Y %7.2f", strings.Repeat(" ", 4), pos.X(), pos.Y())

		pre := clu.Pre(nil)
		fmt.Printf("%4d digits\n", pre.DigitsLength())
	}

}

var outputFileName string

func init() {
	clusterCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output filename")
}
