package casbin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	model2 "github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/wonderivan/logger"
	"rbac-v1/common"
	"rbac-v1/common/constants"
	"rbac-v1/model/po"
	"rbac-v1/model/vo"
	"rbac-v1/service"
	"sync"
)

var (
	m sync.Mutex
	onceInit sync.Once //单例模式，只会执行一次
	cb *Casbin
)

type Casbin struct {
	enforcer *casbin.Enforcer
	policy []byte
	*service.Service
}

type casbinPolicy struct {
	PType string `json:"PType"`
	V0 string `json:"V0"`
	V1 string `json:"V1"`
	V2 string `json:"V2"`
}

//刷新
func (c *Casbin) Refresh() error {
	//获取model
	model, err := model2.NewModelFromString(constants.RBAC_MODEL)
	if err != nil {
		return err
	}
	//获取policy
	b, err := c.generatePolicy(context.Background())
	if err != nil {
		return err
	}
	//生成enforcer
	a := jsonadapter.NewAdapter(&b)
	e, err := casbin.NewEnforcer(model, a)
	if err != nil {
		return err
	}
	m.Lock()
	defer m.Unlock()
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	cb.enforcer = e
	cb.policy = b

	return nil
}

//初始化
// xx.NewCasbin().check
func NewCasbin() *Casbin {
	onceInit.Do(func() {
		cb = &Casbin{
			Service:  service.Srv(),
		}
		if err := cb.Refresh(); err != nil {
			panic(err)
		}
	})

	return cb
}

//生成policy
func (c *Casbin) generatePolicy(ctx context.Context) (ret []byte, err error) {
	var cps = []*casbinPolicy{}

	dataP, err := c.generateP(ctx)
	if err != nil {
		return nil, err
	}
	dataG, err := c.generateG(ctx)
	if err != nil {
		return nil, err
	}

	cps = append(cps, dataP...)
	cps = append(cps, dataG...)
	b, err := json.Marshal(cps)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//生成p，权限-操作
func (c *Casbin) generateP(ctx context.Context) (ret []*casbinPolicy, err error) {
	//获取powerlist
	ret = []*casbinPolicy{}
	powerList, err := c.GetPowerList(ctx, &vo.PowerListRequest{})
	if err != nil {
		return nil, err
	}
	if powerList.Total == 0 {
		return ret, nil
	}
	//组装数据，生成casbinPolicy
	for _, power := range powerList.List {
		if len(power.Operations) == 0 {
			continue
		}
		for _, op := range power.Operations {
			ret = append(ret, &casbinPolicy{
				PType: "p",
				V0:    power.Code,
				V1:    op.Path,
				V2:    op.Method,
			})
		}
	}

	return ret, nil
}

//生成g，用户-角色-权限
func (c *Casbin) generateG(ctx context.Context) (ret []*casbinPolicy, err error) {
	var (
		wg = sync.WaitGroup{}
		process = 5
		errs = []error{}
	)
	ret = []*casbinPolicy{}

	//查询用户
	userList, err := c.GetUserList(ctx, &vo.UserListRequest{})
	if err != nil {
		return nil, err
	}
	if userList.Total == 0 {
		return ret, nil
	}

	wg.Add(process)
	total := int(userList.Total)
	size := total / process //每个协程处理的user数量
	// total=11 process=5 size=2
	// i=0 startId=0 endId=1
	// i=1 startId=2 endId=3
	// i=4 startId=8 endId=10
	for i := 0; i < process; i++ {
		startId := i * size
		endId := (i+1)*size -1
		if i == process -1 {
			endId = (i+1)*size -1 + (total % process)
		}

		go func(i, startId, endId int) {
			for j := startId; j <= endId; j++ {
				user := userList.List[j]
				poMp := map[uint]*po.Power{}
				//处理权限
				for _, power := range user.Powers {
					poMp[power.Id] = power
				}
				//处理角色
				roleIds := []uint{}
				for _, role := range user.Roles {
					roleIds = append(roleIds, role.Id)
				}
				roleList, err := c.GetRoleList(ctx, &vo.RoleListRequest{
					Ids:      common.RemoveDuplicates(roleIds),
				})
				if err != nil {
					errs = append(errs, err)
					continue
				}
				//把角色转换成权限
				for _, role := range roleList.List {
					if len(role.Powers) == 0 {
						continue
					}
					for _, power := range role.Powers {
						poMp[power.Id] = power
					}
				}
				for _, power := range poMp {
					ret = append(ret, &casbinPolicy{
						PType: "g",
						V0:    user.Username,
						V1:    power.Code,
					})
				}
			}
			wg.Done()
		}(i, startId, endId)
	}
	wg.Wait()
	if len(errs) != 0 {
		err := errors.New(fmt.Sprintf("generateG failed: %v", errs))
		return nil, err
	}
	return ret, nil
}

//权限校验
func (c *Casbin) Auth(param *vo.CasbinAuthRequest) (err error) {
	ok, err := c.enforcer.Enforce(param.Username, param.Path, param.Method)
	if err != nil {
		logger.Error("casbin enforce failed: ", err)
		return err
	}
	if !ok {
		err = errors.New("no permission")
		logger.Error(err.Error())
		return err
	}

	return nil
}

//获取policy
func (c *Casbin) GetPolicy() string {
	return string(c.policy)
}