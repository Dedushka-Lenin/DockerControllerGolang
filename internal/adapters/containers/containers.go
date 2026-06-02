package containers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/domain"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/moby/go-archive"
	"github.com/moby/moby/api/pkg/stdcopy"
)

type ContainersRepo interface {
	Create(login, container_id, container_name string) (int, error)
	Delete(login string) error
	Check(login string, id string) (bool, error)
	GetList(login string) ([]domain.Container, error)
	GetById(id string) (*domain.Container, error)
}

type Containers struct {
	ctx    context.Context
	cfg    *config.Config
	client *client.Client
	cr     ContainersRepo
}

func NewContainers(ctx context.Context, cfg *config.Config, cr ContainersRepo) *Containers {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Ошибка инициализации Docker клиента: %v", err)
	}

	return &Containers{ctx: ctx, cfg: cfg, client: client, cr: cr}
}

func (c *Containers) Create(login, name, url string) error {
	path := fmt.Sprintf("repo/%s/%s", login, name)

	if err := c.clone(url, path); err != nil {
		log.Println("Create. clone. err: " + err.Error())
		return err
	}

	defer func() {
		if err := os.RemoveAll(path); err != nil {
			log.Println("Create. RemoveAll. err: " + err.Error())
		}
	}()

	containerId, err := c.createContainer(path, name)
	if err != nil {
		log.Println("Create. createContainer. err: " + err.Error())
		return err
	}

	log.Println(fmt.Sprintf("login - %s, containerId - %s, name - %s", login, containerId, name))
	if _, err := c.cr.Create(login, containerId, name); err != nil {
		log.Println("Create. Create. err: " + err.Error())
		return err
	}

	log.Println("Create. err: nil")
	return nil
}

func (c *Containers) Delete(login string, id string) error {
	containerData, err := c.get(login, id)
	if err != nil {
		log.Println("Delete. get. err: " + err.Error())
		return err
	}

	err = c.client.ContainerRemove(c.ctx, containerData.Id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Println("Delete. ContainerRemove. err: " + err.Error())
		return err
	}

	_, err = c.client.ImageRemove(c.ctx, containerData.Name, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})

	if err != nil {
		log.Println("Delete. ImageRemove. err: " + err.Error())
		return err
	}

	if c.cr.Delete(id); err != nil {
		log.Println("Delete. ImageRemove. err: " + err.Error())
		return err
	}

	log.Println("Delete. err: nil")
	return nil
}

func (c *Containers) GetStatus(login string, id string) (string, error) {
	containerData, err := c.get(login, id)
	if err != nil {
		log.Println("GetStatus. get. err: " + err.Error())
		return "", err
	}

	info, err := c.client.ContainerInspect(c.ctx, containerData.Id)
	if err != nil {
		log.Println("GetStatus. ContainerInspect. err: " + err.Error())
		return "", err
	}

	log.Println("GetStatus. err: nil")
	return info.State.Status, nil
}

func (c *Containers) GetList(login string) ([]domain.Container, error) {
	Cont, err := c.cr.GetList(login)
	return Cont, err
}

func (c *Containers) Logs(login string, data domain.ContainerLogsData) (string, error) {
	containerData, err := c.get(login, data.Id)
	if err != nil {
		log.Println("Logs. get. err: " + err.Error())
		return "", err
	}

	readCloser, err := c.client.ContainerLogs(c.ctx, containerData.Id,
		container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       strconv.Itoa(data.Tail),
		})
	if err != nil {
		log.Println("Logs. ContainerLogs. err: " + err.Error())
		return "", err
	}

	defer readCloser.Close()

	content, err := io.ReadAll(readCloser)
	if err != nil {
		log.Println("Logs. ReadAll. err: " + err.Error())
		return "", err
	}

	log.Println("Logs. err: nil")
	return string(content), nil
}

func (c *Containers) Start(login string, id string) error {
	containerData, err := c.get(login, id)
	if err != nil {
		log.Println("Start. get. err: " + err.Error())
		return err
	}

	err = c.client.ContainerStart(c.ctx, containerData.Id, container.StartOptions{})
	return err
}

func (c *Containers) Stop(login string, id string) error {
	containerData, err := c.get(login, id)
	if err != nil {
		log.Println("Stop. get. err: " + err.Error())
		return err
	}

	err = c.client.ContainerStop(c.ctx, containerData.Id, container.StopOptions{})
	return err
}

func (c *Containers) Restart(login string, id string) error {
	containerData, err := c.get(login, id)
	if err != nil {
		log.Println("Restart. get. err: " + err.Error())
		return err
	}

	err = c.client.ContainerRestart(c.ctx, containerData.Id, container.StopOptions{})
	return err
}

func (c *Containers) Exec(login string, id string, cmd string) (string, error) {
	// Настройка создания exec-процесса
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", cmd},
	}

	// 1. Создаем конфигурацию exec команды
	execIDResp, err := c.client.ContainerExecCreate(c.ctx, id, execConfig)
	if err != nil {
		return "", err
	}

	// 2. ИСПОЛЬЗУЕМ ATTACH ВМЕСТО START: это запускает команду и дает нам поток HijackedResponse
	attachConfig := container.ExecStartOptions{
		Detach: false,
		Tty:    false,
	}
	resp, err := c.client.ContainerExecAttach(c.ctx, execIDResp.ID, attachConfig)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	// Буферы для записи вывода
	var stdout, stderr bytes.Buffer

	// 3. Копируем данные из потока и очищаем их от системных заголовков Docker
	_, err = stdcopy.StdCopy(&stdout, &stderr, resp.Reader)
	if err != nil {
		return "", err
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += "[STDERR]: " + stderr.String()
	}

	if strings.TrimSpace(output) == "" {
		output = "Команда выполнена успешно (вывод пуст)"
	}

	return output, nil
}

func (c *Containers) clone(url, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{URL: url})
	return err
}

func (c *Containers) createContainer(path, name string) (string, error) {
	buildCtx, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		log.Println("createContainer. TarWithOptions. err: " + err.Error())
		return "", err
	}

	buildResponse, err := c.client.ImageBuild(c.ctx, buildCtx, build.ImageBuildOptions{
		Tags:       []string{name},
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		log.Println("createContainer. ImageBuild. err: " + err.Error())
		return "", err
	}

	defer buildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		log.Println("createContainer. ReadBuildBody. err: " + err.Error())
		return "", err
	}

	resp, err := c.client.ContainerCreate(c.ctx, &container.Config{
		Image: name,
		Tty:   true,
	}, nil, nil, nil, name)
	if err != nil {
		log.Println("createContainer. ContainerCreate. err: " + err.Error())
		return "", err
	}

	log.Println("createContainer. err: nil")
	return resp.ID, nil
}

func (c *Containers) get(login string, id string) (*domain.Container, error) {
	if exists, err := c.cr.Check(login, id); !exists {
		log.Println("get. Check. err: container does not exist")
		return nil, fmt.Errorf("container does not exist")
	} else if err != nil {
		log.Println("get. Check. err: " + err.Error())
		return nil, err
	}

	container, err := c.cr.GetById(id)
	if err != nil {
		log.Println("get. GetById. err: " + err.Error())
		return nil, err
	}

	log.Println("get. err: nil")
	return container, nil
}
