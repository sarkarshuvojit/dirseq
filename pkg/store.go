package dirseq

type Store[T any] interface {
	SetupDatabase() (*T, error)
	GetSequence(path string) (int, error)
	GetSequenceAndPadding(path string) (int, int, error)
	UpdateSequence(path string, newSeq int) error
	UpdatePadding(path string, padding int) error
}
