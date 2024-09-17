package queue

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type Client struct {
	Client *asynq.Client
}

type Conf struct {
	RedisAddr string
}

func NewClient(c Conf, logger log.Logger) *Client {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAddr})
	return &Client{
		Client: asynqClient,
	}
}

// PublishTask 发布队列任务
func (r *Client) PublishTask(data []byte, taskType string) (*asynq.TaskInfo, error) {
	return r.Client.Enqueue(asynq.NewTask(taskType, data))
}
