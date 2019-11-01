package impl

import (
	"github.com/dgraph-io/badger"
)

func Copy(
	fromSSTDir string,
	fromVlogDir string,
	toDir string,
) error {
	// Open the source DB.
	source, err := badger.Open(
		badger.DefaultOptions(fromSSTDir).
			WithValueDir(fromVlogDir).
			WithTruncate(true))
	if err != nil {
		return err
	}
	defer source.Close()

	// Open the destination DB.
	dest, err := badger.Open(badger.DefaultOptions(toDir))
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
			err := func() error { // Provides tighter scope for the defer inside.
				batch := dest.NewWriteBatch()
				defer batch.Cancel()

				val, err := iter.Item().ValueCopy(nil)
				if err != nil {
					return err
				}

				err = batch.Set(iter.Item().Key(), val)
				if err != nil {
					return err
				}
				return batch.Flush()
			}()
			if err != nil {
				return err
			}
		}
		return nil
	})
}
