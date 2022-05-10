package logic

import (
	"context"
	"errors"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"cloud-disk/core/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailSendLogic {
	return &EmailSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailSendLogic) EmailSend(req *types.EmailSendRequest) (resp *types.EmailSendReply, err error) {
	resp = &types.EmailSendReply{}
	// 判断邮箱是否注册
	count, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(models.User))
	if err != nil {
		return
	}
	
	if count > 0 {
		return nil, errors.New("该邮箱已被注册")
	}

	//发送邮件
	code := util.RandCode()
	err = util.SendEmail(req.Email, code)
	if err != nil {
		return nil, err
	}

	// redis存储
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))

	resp.Message = "验证码已发送您的邮箱"
	return
}
