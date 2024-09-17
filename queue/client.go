package queue

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type Client struct {
	log    *log.Helper
	Client *asynq.Client
}

type Conf struct {
	RedisAddr string
}

func NewClient(c Conf, logger log.Logger) *Client {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAddr})
	return &Client{
		Client: asynqClient,
		log:    log.NewHelper(logger),
	}
}

// PublishTask 发布队列任务
func (r *Client) PublishTask(data []byte, taskType string) error {
	info, err := r.Client.Enqueue(asynq.NewTask(taskType, data))
	if err != nil {
		r.log.Errorf("could not schedule task: %v", err)
		return err
	}
	r.log.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}
