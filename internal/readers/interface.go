package readers

type IReaderModule interface {
	ReadCard() ([]byte, error)
	Close() error
	Buzz() error
}
