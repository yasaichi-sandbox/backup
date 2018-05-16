package backup

// By convention, one-method interfaces are named by the method name plus an -er
// suffix or similar modification to construct an agent noun: Reader, Writer,
// Formatter, CloseNotifier etc.
// ref. https://golang.org/doc/effective_go.html#interface-names
type Archiver interface {
	Archive(src, dest string) error
}
