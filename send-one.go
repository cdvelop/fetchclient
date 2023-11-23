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
	if body_form, ok := body_rq.(js.Value); !ok {
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

	// h.Log("API endpoint:", endpoint)

	// Crear una función JavaScript que se llamará cuando se complete la solicitud
	h.onComplete = js.FuncOf(func(this js.Value, server []js.Value) interface{} {
		// argumento 0 es el cuerpo de la respuesta de la solicitud Fetch, que debería ser una cadena de texto JSON.
		// argumento 1 indica si la promesa se resolvió o se rechazó.

		h.Log("RESPUESTA:")
		h.Log("TAMAÑO:", len(server))
		h.Log(server[0])

		msg := server[0].Get("statusText").String() //Not Found

		status := server[0].Get("status").String() //<number: 404>
		if status == "<number: 404>" {
			msg += " 404"
		}

		ok := server[0].Get("ok").String() //<boolean: false>
		if ok == "<boolean: false)>" {
			ok = "false"
		}

		// status := res.Header.Get("Status")
		// fmt.Println("ESTATUS GET:", status)

		// if res.StatusCode != 200 {
		// 	response(nil, status))
		// 	return
		// }

		// Decode := res.Header.Get("Decode")

		h.Log("RESP OK:", ok, "status:", status, "text", msg)

		// if len(server) != 2 {
		// 	return msg)
		// }

		// Decodificar la respuesta
		// responseData := h.DecodeResponses([]byte(server[0].String()))

		// Llamar a la función de respuesta de Go con los datos decodificados
		// clientReturn(responseData, nil)

		// Liberar la función JavaScript
		h.onComplete.Release()

		return nil
	})

	// Realizar la solicitud Fetch en JavaScript
	js.Global().Get("fetch").Invoke(endpoint, js.ValueOf(map[string]interface{}{
		"method": method,
		"body":   body,
		// "body":    js.ValueOf(string(body)),
		"headers": js.ValueOf(map[string]interface{}{"Content-Type": content_type}),
	})).Call("then", h.onComplete, js.Null())

}
