package iterate

type Candidate struct {
	Path string
	Reader io.ReadSeeker
}

type Iterator interface {
	Walk(context.Context, string) iter.Seq[*Candidate, error]
}
