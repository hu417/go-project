package logic

import (
	"context"
	"time"

	"cloud-disk/core/define"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListReq, UserIdentity string) (resp *types.UserFileListResp, err error) {
	// todo: add your logic here and delete this line

	// 定义list,count
	uf := make([]*types.UserFile, 0)
	resp = new(types.UserFileListResp)
	//
	// type UserFile struct {
	// 	Id                 int    `json:"id"`
	// 	Identity           string `json:"identity"`
	// 	RepositoryIdentity string `json:"repository_identity"`
	// 	Name               string `json:"name"`
	// 	Ext                string `json:"ext"`
	// 	Path               string `json:"path"`
	// 	Size               string `json:"size"`
	// }

	// type UserRepository struct {
	// 	Id                 int
	// 	Identity           string
	// 	UserIdentity       string
	// 	ParentId           int
	// 	RepositoryIdentity string
	// 	Ext                string
	// 	Name               string
	// 	CreatedAt          time.Time `xorm:"created"`
	// 	UpdatedAt          time.Time `xorm:"updated"`
	// 	DeletedAt          time.Time `xorm:"deleted"`
	// }

	// 分页参数
	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = define.Page
	}
	offsize := (page - 1) * size

	// 查询用户文件列表
	err = l.svcCtx.Engine.Table("user_repository").Where("parent_id = ? AND user_identity = ? ", req.Id, UserIdentity).Select("user_repository.id, user_repository.identity, user_repository.repository_identity, user_repository.ext, user_repository.name, "+" repository_pool.path, repository_pool.size").Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.Datetime)).Limit(size, offsize).Find(&uf)

	if err != nil {
		return
	}

	// 查询用户文件总数
	cnt, err := l.svcCtx.Engine.Where("parent_id = ? AND user_identity = ?", req.Id, UserIdentity).Count(new(models.UserRepository))
	if err != nil {
		return
	}
	resp.List = uf
	resp.Count = cnt

	return
}
