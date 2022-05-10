package logic

import (
	"context"
	"errors"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"cloud-disk/core/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginReply, err error) {
	// 1. 从数据库查询用户
	user := new(models.User)
	has, err := l.svcCtx.Engine.Where("username = ? AND password = ?", req.UserName, util.Md5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名或密码错误")
	}
	// 2. 生成token
	token, err := util.GenerateToken(user.Id, user.UserName, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	resp = new(types.LoginReply)
	resp.Token = token
	return
}
