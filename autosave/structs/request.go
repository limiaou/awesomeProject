package structs

type Save struct {
	SaveList []SaveList `json:"save_list"`
}
type SaveList struct {
	ResourceId    []int64 `json:"resource_id"`
	EnvironmentId int64   `json:"environment_id"`
	AppId         int64   `json:"app_id"`
}
