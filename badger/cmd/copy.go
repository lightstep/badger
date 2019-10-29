package cmd

import (
	"fmt"

	"github.com/dgraph-io/badger/badger/cmd/impl"
	"github.com/spf13/cobra"
)

var outputDir string

// copyCmd represents the copy command.
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy Badger database.",
	Long:  "Copy Badger database, to eliminate deletion tombstones",
	RunE:  doCopy,
}

func init() {
	RootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVar(
		&outputDir,
		"output-dir",
		"/dev/null",
		"Location of the copied database",
	)
}

func doCopy(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("no extra args allowed: %v", args)
	}
	return impl.Copy(sstDir, vlogDir, outputDir)
}
