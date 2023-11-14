package fetchclient

import (
	"github.com/cdvelop/model"
)

func AddFetchAdapter(h *model.Handlers) error {
	const err = "error fetchclient nil"
	if h.Logger == nil {
		return model.Error(err, "Logger")
	}
	if h.DataConverter == nil {
		return model.Error(err, "DataConverter")
	}

	f := fetchClient{
		DataConverter: h,
		Logger:        h,
	}

	h.FetchAdapter = f

	return nil
}
