package service

import (
	"battle-tiles/internal/biz"
	cloudRepo "battle-tiles/internal/dal/repo/cloud"
	pdb "battle-tiles/pkg/plugin/dbx"
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

type HandlerWithCtx func(context.Context) error

type TaskPayload struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Handler func(string, *TaskPayload) error

type AsyNQService struct {
	log           *log.Helper
	SubscriberMap map[string]Handler
	uc            *biz.AsyNQUseCase
	cloudRepo     cloudRepo.BasePlatformRepo
}

func NewAsyNQService(
	logger log.Logger,
	uc *biz.AsyNQUseCase,
	cloudRepo cloudRepo.BasePlatformRepo,
) *AsyNQService {
	s := &AsyNQService{
		log:       log.NewHelper(log.With(logger, "module", "service/asynq")),
		uc:        uc,
		cloudRepo: cloudRepo,
	}
	s.initSubscriber()
	return s
}

func (s *AsyNQService) initSubscriber() {
	s.SubscriberMap = map[string]Handler{
		"refresh:test:job": func(msg string, payload *TaskPayload) error {
			s.log.Infof("refresh:test:job start running [%s]", payload.Message)
			return nil
		},
	}
}
func (s *AsyNQService) AutoPlatformExec(funcWithCtx HandlerWithCtx, taskType string, taskPayload TaskPayload) {
	if funcWithCtx == nil {
		s.log.Error("执行函数不能为nil,请检查")
		return
	}
	paltformList, err := s.cloudRepo.ListPlatform()
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	for _, platform := range paltformList {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx := context.Background()
			// 给上下文添加平台机构信息
			ctx = context.WithValue(ctx, pdb.CtxDBKey, platform.Platform)
			//fmt.Println("ctx: ", ctx)

			if err = funcWithCtx(ctx); err != nil {
				s.log.Errorw("action", "AutoPlatform", "cloud.Code", platform.Platform, "cloud.Name", platform.Name, "cloud.DBName", platform.DBName, "err", err)
			}
		}()
	}
	wg.Wait()
}
