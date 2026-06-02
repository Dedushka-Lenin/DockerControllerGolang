package domain

type Container struct {
	Id   string `db:"container_id"`
	Name string `db:"container_name"`
}

type ContainerId struct {
	Id string
}

type ContainerLogsData struct {
	Id   string
	Tail int
}

type ContainerCreate struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ExecRequest struct {
	ID  string `json:"Id" binding:"required"`
	Cmd string `json:"cmd" binding:"required"`
}
