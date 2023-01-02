package common

import (
	"fmt"

	"github.com/antonholmquist/jason"

	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
)

type SmwQuery struct {
	w      *mwclient.Client
	params params.Values
	resp   *jason.Object
	err    error
}

// Err returns the first error encountered by the Next method.
func (q *SmwQuery) Err() error {
	return q.err
}

// Resp returns the API response retrieved by the Next method.
func (q *SmwQuery) Resp() *jason.Object {
	return q.resp
}

func NewSmwQuery(w *mwclient.Client, p params.Values) *SmwQuery {
	p.Set("action", "askargs")

	return &SmwQuery{
		w:      w,
		params: p,
		resp:   nil,
		err:    nil,
	}
}

func (q *SmwQuery) Next() (done bool) {
	if q.resp == nil {
		// first call to Next
		q.resp, q.err = q.w.Get(q.params)
		return q.err == nil
	}

	cont, err := q.resp.GetNumber("query-continue-offset")
	if err != nil {
		return false
	}

	offset := fmt.Sprintf("|offset=%s", cont)
	currentParams := q.params.Get("parameters")
	q.params.Set("parameters", currentParams+offset)

	q.resp, q.err = q.w.Get(q.params)
	q.params.Set("parameters", currentParams)
	return q.err == nil
}
