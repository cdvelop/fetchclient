package fetchclient

import (
	"github.com/cdvelop/model"
)

func AddFetchAdapter(h *model.Handlers) (err string) {
	const e = "error fetchclient nil "
	if h.Logger == nil {
		return e + "Logger"
	}
	if h.DataConverter == nil {
		return e + "DataConverter"
	}

	f := fetchClient{
		DataConverter: h,
		Logger:        h,
	}

	h.FetchAdapter = f

	return ""
}
