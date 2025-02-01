package tools

import "net/http"

func ShowError(writer http.ResponseWriter, message string) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>BAD REQUEST!</title>
</head>
<body>
<h1>400, BadRequest!</h1><br>
<span style="color: red;font-style: italic;">` + message + `</span>
</body>
</html>`))
}
