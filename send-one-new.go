package fetchclient

import "syscall/js"

func (h *fetchClient) SendOneRequestNEW(method, endpoint, object string, body_rq interface{}, response func(result []map[string]string, err string)) {

	// Liberar la función JavaScript onComplete si ya existe
	if h.abortController != nil {
		// Abortar la solicitud en curso
		h.abortController.Call("abort")
	}

	// h.onComplete.Release()
	// Crear un nuevo controlador AbortController
	abortController := js.Global().Get("AbortController").New()
	h.abortController = &abortController

	var back string
	if object != "" {
		back = "/"
	}

	endpoint = endpoint + back + object

	var content_type = "application/json"

	var body string

	// si envían un tipo objeto JavaScript
	if body_form, ok := body_rq.(js.Value); ok {
		body = body_form.String()
		content_type = "multipart/form-data"
	} else {
		body_byte, err := h.EncodeMaps(body_rq, object)
		if err != "" {
			response(nil, err)
			return
		}
		body = string(body_byte)
	}

	// Variable para realizar un seguimiento del número de solicitudes pendientes
	var pendingRequests int

	// Crear una función JavaScript que se llamará cuando se complete la solicitud Fetch
	h.onComplete = js.FuncOf(func(this js.Value, res []js.Value) interface{} {
		// argumento 0 es el cuerpo de la respuesta de la solicitud Fetch
		// argumento 1 indica si la promesa se resolvió o se rechazó.

		h.Log("RES:", res[0])
		h.Log("RES STRING:", res[0].String())
		h.Log("BODY:", res[0].Get("body"))

		// Extraer el cuerpo de la respuesta usando el método text()
		bodyPromise := res[0].Call("text")

		// Manejar la promesa para obtener el cuerpo real
		bodyPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			h.Log("ARGS", args)
			// args[0] contiene el cuerpo de la respuesta

			body := args[0].String()
			h.Log("BODY:", body)

			// statusText := res[0].Get("statusText").String() // Not Found
			status_code := res[0].Get("status").String() // <number: 404>
			h.Log("status_code:", status_code)

			if status_code != "<number: 200>" {
				response(nil, "error "+body)
			}

			// Liberar la función JavaScript y el controlador AbortController después de ambos bloques then y catch
			pendingRequests--
			if pendingRequests == 0 {
				h.onComplete.Release()
			}

			return nil
		}), js.Null())

		return nil
	})

	// Incrementar el contador de solicitudes pendientes antes de realizar la solicitud
	pendingRequests++

	// Realizar la solicitud Fetch en JavaScript con el controlador AbortController
	fetchOptions := js.Global().Get("Object").New()
	fetchOptions.Set("method", method)

	if method != "GET" {
		fetchOptions.Set("body", body)
	}

	fetchOptions.Set("headers", js.ValueOf(map[string]interface{}{"Content-Type": content_type}))
	fetchOptions.Set("signal", h.abortController.Get("signal"))

	fetchPromise := js.Global().Get("fetch").Invoke(endpoint, fetchOptions).Call("then", h.onComplete, js.Null())

	fetchPromise.Call("then", h.onComplete, js.Null()).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// args[0] contiene el error
		errorMessage := args[0].Get("message").String()
		// h.Log("Error en la solicitud Fetch:", errorMessage)

		// Liberar la función JavaScript y el controlador AbortController
		// h.onComplete.Release()
		// Decrementar el contador de solicitudes pendientes y liberar la función si no hay más solicitudes pendientes

		// Llamar a la función de respuesta con el error
		response(nil, "Error en la solicitud Fetch: "+errorMessage)

		pendingRequests--
		if pendingRequests == 0 {
			h.onComplete.Release()
		}

		return nil
	}), js.Null())

}
