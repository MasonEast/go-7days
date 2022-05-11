package logic

import (
	"context"
	"errors"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"cloud-disk/core/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil{
		return nil, err
	}
	if code != req.Code {
		return nil, errors.New("验证码错误")
	}

	count, err := l.svcCtx.Engine.Where("username = ?", req.UserName).Count(new(models.User))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该用户已被注册")
	}

	_, err = l.svcCtx.Engine.Insert(&models.User{
		Identity: util.UUID(),
		UserName: req.UserName,
		Password: util.Md5(req.Password),
		Email: req.Email,
		CreatedTime: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UserRegisterReply{
		Message: "注册成功",
	}

	return
}
