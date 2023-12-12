package fetchclient

import (
	"syscall/js"
)

// action ej: "create","update","delete","upload" = "POST",  "file" and "read" == GET
// object ej: patientcare.printdoc
func (h fetchClient) SendOneRequest(method, endpoint, object string, body_rq any, response func(result []map[string]string, err string)) {

	var back string
	if object != "" {
		back = "/"
	}

	endpoint = endpoint + back + object

	var content_type = "application/json"

	var body string

	// si envían un tipo objeto javascript
	if body_form, ok := body_rq.(js.Value); ok {
		body = body_form.String()
		content_type = "multipart/form-data"
	} else {

		// h.Log("ENCODE MAPS?")
		body_byte, err := h.EncodeMaps(body_rq, object)
		if err != "" {
			response(nil, err)
			return
		}

		body = string(body_byte)
	}

	// h.Log("API endpoint:", endpoint)

	// Crear una función JavaScript que se llamará cuando se complete la solicitud
	h.onComplete = js.FuncOf(func(this js.Value, res []js.Value) interface{} {

		// Extraer el cuerpo de la respuesta usando el método text()
		bodyPromise := res[0].Call("text")

		// Manejar la promesa para obtener el cuerpo real
		bodyPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// args[0] contiene el cuerpo de la respuesta

			body := args[0].String()
			// h.Log("bodyPromise:", body)

			// statusText := res[0].Get("statusText").String() // Not Found
			status_code := res[0].Get("status").String() // <number: 404>
			// h.Log("status_code:", status_code)

			if status_code != "<number: 200>" {
				response(nil, "error "+body)
			}

			out, err := h.DecodeMaps([]byte(body))
			if err != "" {
				response(nil, err)
			}

			response(out, "")

			// Liberar la función JavaScript
			h.onComplete.Release()

			return nil
		}), js.Null())

		return nil
	})

	fetchOptions := js.Global().Get("Object").New()
	fetchOptions.Set("method", method)

	if method != "GET" {
		fetchOptions.Set("body", body)
	}

	auth := h.AddHeaderAuthentication()

	fetchOptions.Set("headers", js.ValueOf(map[string]interface{}{
		"Content-Type": content_type,
		auth.Name:      auth.Content,
	}))

	// Realizar la solicitud Fetch en JavaScript
	js.Global().Get("fetch").Invoke(endpoint, fetchOptions).Call("then", h.onComplete, js.Null())

}
