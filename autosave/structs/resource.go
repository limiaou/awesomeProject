package structs

type ResponseResource struct {
	Code  int    `json:"code"`
	Data  Data   `json:"data"`
	Error string `json:"error"`
}

type Resource struct {
	ResourceId int64 `json:"resource_id"`
}

type Data struct {
	ResourceList []Resource `json:"resource_list"`
}
