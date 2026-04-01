package httpbin

import (
	"context"
	"github.com/mono83/transport/http/filters"
	"github.com/mono83/transport/http/json"
)

type PostCall struct {
	Name  string
	Value int
}

func (p PostCall) Execute(ctx context.Context, t Client) (*Response, error) {
	return json.ReadJSON[Response](
		filters.Require200(t.ExecuteRequest(ctx, "POST", "/post", json.WriteIndent(p))),
	)
}
