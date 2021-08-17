package routes

type NormalResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ResponseData struct {
	Slug  string `json:"slug"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Key   string `json:"key"`
}

type DataResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}