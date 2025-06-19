package iterate

import (
	"context"
	"iter"
)

type Source interface {
	Walk(context.Context, string) iter.Seq2[*Record, error]
}
