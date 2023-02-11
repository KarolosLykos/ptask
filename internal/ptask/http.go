package ptask

import (
	"net/http"
)

type Handlers interface {
	List() func(w http.ResponseWriter, r *http.Request)
}
