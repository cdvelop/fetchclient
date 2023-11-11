package fetchclient

import (
	"github.com/cdvelop/model"
)

func AddFetchAdapter(h *model.Handlers) {

	f := fetchClient{
		DataConverter: h,
		Logger:        h,
	}

	h.FetchAdapter = f

}
