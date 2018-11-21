package cmd

import (
	"github.com/aphecetche/galo/mathieson"
	"github.com/spf13/cobra"
)

// mathiesonCmd represents the mathieson command
var mathiesonCmd = &cobra.Command{
	Use:   "mathieson",
	Short: "plot some Mathieson functions",
	Run: func(cmd *cobra.Command, args []string) {
		mathieson.MakePlots("mathieson")
	},
}

func init() {
	rootCmd.AddCommand(mathiesonCmd)
}
