package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// NOTE: By convention, one-method interfaces are named by the method name plus
// an -er suffix or similar modification to construct an agent noun: Reader,
// Writer, Formatter, CloseNotifier etc.
// ref. https://golang.org/doc/effective_go.html#interface-names
type Archiver interface {
	Archive(src, dest string) error
	DestFmt() string
}

type zipper struct{}

// NOTE: Type conversion from nil to *zipper
var ZIP Archiver = (*zipper)(nil)

func (*zipper) Archive(src, dest string) error {
	// NOTE: `filepath.Dir` acts like "dirname" command in UNIX, and `os.MkdirAll`
	// does like "mkdir" with -p option.
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	// NOTE: `filepath.Walk` returns an error that the callback returns.
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// NOTE: Create adds a file to the zip file using the provided name.
		// It returns a Writer to which the file contents should be written.
		f, err := w.Create(path)
		if err != nil {
			return err
		}

		io.Copy(f, in)
		return nil
	})
}

func (*zipper) DestFmt() string {
	return "%d.zip"
}
