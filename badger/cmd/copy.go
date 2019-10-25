package cmd

import (
	"github.com/dgraph-io/badger"
	"github.com/spf13/cobra"
)

var outputDir string

// backupCmd represents the backup command
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
	// Open the source DB.
	source, err := badger.Open(
		badger.DefaultOptions(sstDir).
			WithValueDir(vlogDir).
			WithTruncate(true))
	if err != nil {
		return err
	}
	defer source.Close()

	// Open the destination DB.
	dest, err := badger.Open(badger.DefaultOptions(outputDir))
	if err != nil {
		return err
	}
	defer dest.Close()

	// Do the copy.
	return source.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.IteratorOptions{
			AllVersions: false,
		})
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			batch := dest.NewWriteBatch()
			defer batch.Cancel()

			val, err := iter.Item().ValueCopy(nil)
			if err != nil {
				return err
			}

			batch.Set(iter.Item().Key(), val)
			err = batch.Flush()
			if err != nil {
				return err
			}
		}
		return nil
	})
}
