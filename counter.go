package iterate

type Counter interface {
	Seen() int64
	IsIndexing() bool
}
