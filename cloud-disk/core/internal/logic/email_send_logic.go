package logic

import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
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
	err = util.SendEmail(req.Email, "123456")

	if err != nil {
		return nil, err
	}
	resp.Message = "验证码已发送您的邮箱"
	return
}
