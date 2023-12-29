package handler

import "github.com/julienschmidt/httprouter"

// Handler interface needed for a quick router change
type Handler interface {
	Register(router httprouter.Router)
}
