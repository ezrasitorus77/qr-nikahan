package domain

type (
	Response struct {
		Data  interface{}   `json:"data"`
		Error ErrorResponse `json:"error"`
	}

	ErrorResponse struct {
		Err  interface{} `json:"errData"`
		Desc string      `json:"description"`
	}
)
