package fetchclient

import "github.com/cdvelop/model"

func (h fetchClient) SendAllRequests(endpoint string, data []model.Response, response func([]model.Response, error)) {

	response(nil, model.Error("error SendAllRequests no implementado en fetchClient"))

}
