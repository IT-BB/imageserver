// Source http parser
package source

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

// Takes the "source" query parameter
type SourceParser struct {
}

func (parser *SourceParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}
	return nil
}
