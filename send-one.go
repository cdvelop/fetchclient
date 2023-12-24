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

	fetchOptions := js.Global().Get("Object").New()
	fetchOptions.Set("method", method)

	// si envían un tipo objeto javascript
	if form, ok := body_rq.(js.Value); ok {
		// h.Log("ENVIANDO TIPO multipart/form-data")
		fetchOptions.Set("body", form)
	} else {
		// h.Log("ENCODE MAPS?")
		body_byte, err := h.EncodeMaps(body_rq)
		if err != "" {
			response(nil, err)
			return
		}

		fetchOptions.Set("body", string(body_byte))

		fetchOptions.Set("headers", js.ValueOf(map[string]interface{}{
			"Content-Type": content_type,
		}))
	}

	// h.Log("API endpoint:", endpoint)
	// Crear una función JavaScript que se llamará cuando se complete la solicitud
	h.onComplete = js.FuncOf(func(this js.Value, res []js.Value) interface{} {

		// Extraer el cuerpo de la respuesta usando el método text()
		bodyPromise := res[0].Call("text")

		// Manejar la promesa para obtener el cuerpo real
		bodyPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// args[0] contiene el cuerpo de la respuesta
			var err string
			var data []map[string]string
			body := args[0].String()
			// h.Log("bodyPromises 2:", body)

			// statusText := res[0].Get("statusText").String() // Not Found
			status_code := res[0].Get("status").String() // <number: 404>
			// h.Log("status_code:", status_code)

			if status_code != "<number: 200>" {
				err = body
			} else {
				data, err = h.DecodeMaps([]byte(body))
			}

			// h.Log("SALIDA FETCH:", out, err)

			if err != "" {
				response(nil, err)
			} else {
				response(data, "")
			}

			// Liberar la función JavaScript
			h.onComplete.Release()

			return nil
		}), js.Null())

		return nil
	})

	// Realizar la solicitud Fetch en JavaScript
	js.Global().Get("fetch").Invoke(endpoint, fetchOptions).Call("then", h.onComplete, js.Null())

}
