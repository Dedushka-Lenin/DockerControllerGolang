package interactors

import (
	"log"
	"net/http"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"
	tokenM "github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"

	"github.com/gin-gonic/gin"
)

type Containers interface {
	Create(login, name, url string) error
	Delete(login string, id string) error

	GetStatus(login string, id string) (string, error)
	GetList(login string) ([]domain.Container, error)

	Logs(login string, data domain.ContainerLogsData) (string, error)

	Start(login string, id string) error
	Stop(login string, id string) error
	Restart(login string, id string) error

	Exec(login string, id string, cmd string) (string, error)
}

type ContainersHandlers struct {
	cfg       *config.Config
	container Containers
	tkn       Token
}

func NewContainersHandlers(cfg *config.Config, container Containers, tkn Token) *ContainersHandlers {
	return &ContainersHandlers{cfg: cfg, container: container, tkn: tkn}
}

func (ch *ContainersHandlers) Create(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Create. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Create. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var cc domain.ContainerCreate
	if err := c.ShouldBindJSON(&cc); err != nil {
		log.Println("Create. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if err = ch.container.Create(login, cc.Name, cc.Url); err != nil {
		log.Println("Create. Create. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.CreateError})
		return
	}

	log.Println("Create. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (ch *ContainersHandlers) Delete(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Delete. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Delete. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var containerData domain.ContainerId
	if err := c.ShouldBindJSON(&containerData); err != nil {
		log.Println("Delete. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if err = ch.container.Delete(login, containerData.Id); err != nil {
		log.Println("Delete. Delete. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.DeleteError})
		return
	}

	log.Println("Delete. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (ch *ContainersHandlers) GetStatus(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("GetStatus. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("GetStatus. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var containerData domain.ContainerId
	if err := c.ShouldBindJSON(&containerData); err != nil {
		log.Println("GetStatus. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	status, err := ch.container.GetStatus(login, containerData.Id)
	if err != nil {
		log.Println("GetStatus. GetStatus. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.GetByIdError})
		return
	}

	log.Println("GetStatus. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil, "status": status})
}

func (ch *ContainersHandlers) GetList(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("GetList. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("GetList. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	containersList, err := ch.container.GetList(login)
	if err != nil {
		log.Println("GetList. GetList. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.GetListError})
		return
	}

	log.Println("GetList. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil, "containers_list": containersList})
}

func (ch *ContainersHandlers) Logs(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Logs. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Logs. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var ContainerLogsData domain.ContainerLogsData
	if err := c.ShouldBindJSON(&ContainerLogsData); err != nil {
		log.Println("Logs. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	containerLogs, err := ch.container.Logs(login, ContainerLogsData)
	if err != nil {
		log.Println("Logs. Logs. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.LogsError})
		return
	}

	log.Println(containerLogs)

	log.Println("Logs. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil, "logs": containerLogs})
}

func (ch *ContainersHandlers) Start(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Start. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Start. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var containerData domain.ContainerId
	if err := c.ShouldBindJSON(&containerData); err != nil {
		log.Println("Start. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if err = ch.container.Start(login, containerData.Id); err != nil {
		log.Println("Start. Start. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.StartError})
		return
	}

	log.Println("Start. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (ch *ContainersHandlers) Stop(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Stop. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Stop. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var containerData domain.ContainerId
	if err := c.ShouldBindJSON(&containerData); err != nil {
		log.Println("Stop. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if err = ch.container.Stop(login, containerData.Id); err != nil {
		log.Println("Stop. Stop. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.StopError})
		return
	}

	log.Println("Stop. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (ch *ContainersHandlers) Restart(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Restart. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Restart. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var containerData domain.ContainerId
	if err := c.ShouldBindJSON(&containerData); err != nil {
		log.Println("Restart. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ReceivingError})
		return
	}

	if err = ch.container.Restart(login, containerData.Id); err != nil {
		log.Println("Restart. Restart. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.RestartError})
		return
	}

	log.Println("Restart. err: nil")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func (ch *ContainersHandlers) Exec(c *gin.Context) {
	token, err := ch.tkn.GetToken(c)
	if err != nil {
		log.Println("Exec. GetToken. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	login, err := ch.tkn.GetLogin(token)
	if err != nil {
		log.Println("Exec. GetLogin. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenM.InvalidToken})
		return
	}

	var req domain.ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Exec. ShouldBindJSON. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса: " + err.Error()})
		return
	}

	output, err := ch.container.Exec(login, req.ID, req.Cmd)
	if err != nil {
		log.Println("Exec. Exec. err: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  nil,
		"output": output,
	})
}
