package backup

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
)

func DirHash(path string) (string, error) {
	// NOTE: `io.Writer` is embedded in `hash.Hash`
	hash := md5.New()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Fprint(hash, path)
		fmt.Fprint(hash, info.IsDir())
		fmt.Fprint(hash, info.Mode())
		fmt.Fprint(hash, info.ModTime())
		fmt.Fprint(hash, info.Name())
		fmt.Fprint(hash, info.Size())

		return nil
	})
	if err != nil {
		return "", err
	}

	// `%x`: base 16, with lower-case letters for a-f
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
