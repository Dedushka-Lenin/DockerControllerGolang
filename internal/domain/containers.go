package domain

type Container struct {
	Id   string `db:"container_id"`
	Name string `db:"container_name"`
}

type ContainerId struct {
	Id int
}

type ContainerLogsData struct {
	Id   int
	Tail int
}

type ContainerCreate struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
