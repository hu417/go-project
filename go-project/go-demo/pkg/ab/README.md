

# 开关


## 降级

降级开关一般需要放到动态配置中心

```go

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// 定义服务接口
type Service interface {
    DoSomething() (string, error)
}

// 正常服务的实现
type NormalService struct{}

func (s *NormalService) DoSomething() (string, error) {
    return "I am normal service", nil
}

// 降级服务的实现
type DegradedService struct{}

func (s *DegradedService) DoSomething() (string, error) {
    // 在这里你可以返回一些默认值或者缓存数据，或者抛出一个已知的异常
    return "I am degraded service", nil
}

var (
    // 服务降级开关
    DegradeService bool
)

func main() {
    normalService := &NormalService{}
    degradedService := &DegradedService{}

    // 模拟系统状态
    go func() {
        for {
            // 随机改变服务降级开关状态
            DegradeService = rand.Intn(2) == 1
            fmt.Println("Service degrade status:", DegradeService)
            time.Sleep(1 * time.Second)
        }
    }()

    for {
        var service Service
        if DegradeService {
            service = degradedService
        } else {
            service = normalService
        }

        result, err := service.DoSomething()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        fmt.Println("Result:", result)
        time.Sleep(1 * time.Second)
    }
}

```


## 灰度

1、随机百分比放量，针对允许用户可一会是走新逻辑，一会又可以走旧逻辑的情况
```go

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// 定义处理请求的接口
type Handler interface {
    HandleRequest()
}

// 新特性处理请求的实现
type NewFeatureHandler struct {
}

func (f *NewFeatureHandler) HandleRequest() {
    fmt.Println("New feature handling request")
}

// 老特性处理请求的实现
type OldFeatureHandler struct {
}

func (f *OldFeatureHandler) HandleRequest() {
    fmt.Println("Old feature handling request")
}

func main() {
    // 设置随机数种子
    rand.Seed(time.Now().UnixNano())

    newHandler := &NewFeatureHandler{}
    oldHandler := &OldFeatureHandler{}

    // 此处设置的是灰度放量的百分比,这个10以及下面的100可以放到配置中心动态配置，逐步加大放量人群
    percentage := 10

    for i := 0; i < 100; i++ {
        if rand.Intn(100) < percentage {
            newHandler.HandleRequest()
        } else {
            oldHandler.HandleRequest()
        }
    }
}

```

2、基于用户百分比放量，一个用户从命中放量访问新特性起就应该一直是出于放量中，能够访问到新特性
```go

package gray

import (
	"context"
	"google.golang.org/appengine/log"
	"hash/fnv"
)

// ControlGrayRule 灰度控制规则, 按字段优先级逐个判断规则, 未命中任何规则默认不灰度
type ControlGrayRule struct {
	GlobalSwitch     bool     `json:"global_switch"`      // 优先级1: 全局灰度开关, false-不允许灰度, 止损时一键关闭灰度
	BlackList        []string `json:"black_list"`         // 优先级2: 灰度黑名单
	WhiteList        []string `json:"white_list"`         // 优先级3: 灰度白名单
	GrayThousandRate uint32   `json:"gray_thousand_rate"` // 优先级3: 灰度比例, 千分比, 0-1000
}

// IsHitGray 基于ID的灰度控制, 命中灰度返回true
func IsHitGray(ctx context.Context, grayID string, rule *ControlGrayRule) bool {
	// 全局灰度开关关闭, 不灰度
	if !rule.GlobalSwitch {
		log.Infof(ctx, "[IsHitGray] gray global switch is close, not gray: grayID=%s", grayID)
		return false
	}
	// 命中黑名单, 不灰度
	if containsString(rule.BlackList, grayID) {
		log.Infof(ctx, "[IsHitGray] hit black list, not gray: grayID=%s", grayID)
		return false
	}
	// 命中白名单, 进行灰度
	if containsString(rule.WhiteList, grayID) {
		log.Infof(ctx, "[IsHitGray] hit white list, will gray: grayID=%s", grayID)
		return true
	}
	// 命中灰度比例, 进行灰度
	grayHash, err := hashStringToInt(grayID)
	if err != nil {
		// 实际上不会触发，因为 hash.Hash32 的 Write 方法不会返回 err
		return false
	}
	log.Infof(ctx, "[IsHitGray] grayID=%s, grayHash=%d, GrayThousandRate=%d",
		grayID, grayHash, rule.GrayThousandRate)
	if grayHash%1000 < rule.GrayThousandRate {
		log.Infof(ctx, "[IsHitGray] hit gray percent, will gray: grayID=%s", grayID)
		return true
	}
	// 未命中任何规则, 默认不灰度
	log.Infof(ctx, "[IsHitGray] does not hit any gray rule, not gray: grayID=%s", grayID)
	return false
}

func hashStringToInt(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func containsString(list []string, userId string) bool {
	for _, val := range list {
		if val == userId {
			return true
		}
	}
	return false
}


```

## ceds


## 