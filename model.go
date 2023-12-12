package fetchclient

import (
	"syscall/js"

	"github.com/cdvelop/model"
)

type fetchClient struct {
	model.DataConverter
	model.Logger
	onComplete js.Func

	// campos nuevos
	// controller js.Value
	abortController *js.Value

	model.HeaderAuthenticationAdapter
}
