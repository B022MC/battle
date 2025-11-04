package server

import (
	"battle-tiles/internal/conf"
	"battle-tiles/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	asynq2 "github.com/hibiken/asynq"
	"github.com/tx7do/kratos-transport/transport/asynq"
)

func NewAsyNQServer(c *conf.Server, logger log.Logger, service *service.AsyNQService) (*asynq.Server, error) {
	srv := asynq.NewServer(
		asynq.WithLogger(log.NewHelper(logger)),
		asynq.WithAddress(c.Asynq.Addr),
		asynq.WithRedisPassword(c.Asynq.Password),
		asynq.WithRedisDatabase(int(c.Asynq.Db)),
	)
	helper := log.NewHelper(logger)
	for _, config := range c.Asynq.Subscriber {
		name := config.GetName()
		schedule := config.GetSchedule()
		if handler, ok := service.SubscriberMap[name]; ok {
			if err := asynq.RegisterSubscriber(srv, name, handler); err != nil {
				return nil, err
			}
			if schedule != "" {
				if _, err := srv.NewPeriodicTask(schedule, name, asynq2.Task{}); err != nil {
					return nil, err
				}
			}
			helper.Infow("action", "register subscriber", "subscriber", name, "schedule", schedule)
		}

	}
	return srv, nil
}
