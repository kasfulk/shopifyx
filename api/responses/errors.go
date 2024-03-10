package responses

func ErrorBadRequest(m string) (int, map[string]interface{}) {
	return 422, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}

func ErrorServer(m string) (int, map[string]interface{}) {
	return 500, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}
