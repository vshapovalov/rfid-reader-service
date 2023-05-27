package readers

type IReaderModule interface {
	ReadCards() ([][]byte, error)
	Close() error
	Buzz() error
	GetReaderInfo() string
}
