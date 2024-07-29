
> [参考文章](https://juejin.cn/post/7241938328839618597)

在 Go 语言中，函数选项模式是一种优雅的设计模式，用于处理函数的可选参数。它提供了一种灵活的方式，允许用户在函数调用时传递一组可选参数，而不是依赖于固定数量和顺序的参数列表。

函数选项模式的好处
- 易于使用：调用者可以选择性的设置函数参数，而不需要记住参数的顺序和类型；
- 可读性强：函数选项模式的代码有着自文档化的特点，调用者能够直观地理解代码的功能；
- 扩展性好：通过添加新的 Option 参数选项，函数可以方便地扩展功能，无需修改函数的签名；
- 函数选项模式可以提供默认参数值，以减少参数传递的复杂性。


函数选项模式的实现一般包含以下几个部分：
- 选项结构体：用于存储函数的配置参数
- 选项函数类型：接收选项结构体参数的函数
- 定义功能函数：接收 0 个或多个固定参数和可变的选项函数参数
- 设置选项的函数：定义多个设置选项的函数，用于设置选项

存在一些缺点:
- 复杂性：函数选项模式引入了更多的类型和概念，需要更多的代码和逻辑来处理。这增加了代码的复杂性和理解的难度，尤其是对于初学者来说。
- 可能存在错误的选项组合：由于函数选项模式允许在函数调用中指定多个选项，某些选项之间可能存在冲突或不兼容的情况。这可能导致意外的行为或错误的结果。
- 不适用于所有情况：函数选项模式适用于有大量可选参数或者可配置选项的函数，但对于只有几个简单参数的函数，使用该模式可能过于复杂和冗余。在这种情况下，简单的命名参数可能更直观和易于使用。

## 场景一: redis参数

```go

const (
	MaxTotal = 8
)

// RedisPoolConfig 连接池配置
type RedisPoolConfig struct {
	name     string // 名称
	maxTotal int    // 最大连接数
}

func (c *RedisPoolConfig) check() error {
	if c.maxTotal == 0 {
		c.maxTotal = MaxTotal
	}
	if len(c.name) == 0 {
		return errors.New("name 为空")
	}

	return nil
}

type Option func(option *RedisPoolConfig)

func WithMaxTotal(maxTotal int) Option {
	return func(options *RedisPoolConfig) {
		options.maxTotal = maxTotal
	}
}

func WithName(name string) Option {
	return func(options *RedisPoolConfig) {
		options.name = name
	}
}

func NewConfig(opts ...Option) (*RedisPoolConfig, error) {
	c := &RedisPoolConfig{}
	for _, opt := range opts {
		opt(c)
	}

	return c, c.check()
}

func TestOption(t *testing.T) {
	c, err := NewConfig(WithName("test"), WithMaxTotal(8))
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
}

```

解释下上面这段代码：
1、  首先定义 Option 变量，类型是func(option *RedisPoolConfig)。
2、  定义若干个高阶函数 WithName(name string) 、WithMaxTotal(maxTotal int)...，返回值都是 Option。
3、  高阶函数返回值 Option 作为 NewConfig(opts ...Option) 函数参数，该函数遍历 opts 分别调用 WithXX 方法给 RedisPoolConfig 设置字段值。


## 场景二: db Update 方法

```go

type UserRepo interface {
	UpdateMap(ctx context.Context, query *UserQueryOption, opts ...UpdateOption) (int64, error) // 更新数据
}

type UserRepoImpl struct {
}

func NewUserRepoImpl() *UserRepoImpl {
	return &UserRepoImpl{}
}

func (d *UserRepoImpl) UpdateMap(ctx context.Context, query *UserQueryOption, opts ...UpdateOption) (int64, error) {
	if len(opts) == 0 {
		return 0, errors.New("opts 为空暂无数据更新")
	}

	m := make(map[string]interface{})
	for _, v := range opts {
		v(m) // 调用 opt 为 map 设置值
	}

	// 省略数据库更新操作，query 对象字段用于更新操作的检索条件 比如：where id = query.ID
	return 10, nil
}

type UpdateOption func(op map[string]interface{})

// WithFinishTime 更新完成时间
func WithFinishTime() UpdateOption {
	return func(opt map[string]interface{}) {
		opt["finish_time"] = time.Now().UnixMilli()
	}
}

// WithStatus 更新状态
func WithStatus(status int) UpdateOption {
	return func(opt map[string]interface{}) {
		opt["status"] = status
	}
}

type UserQuery struct {
	ID  string
	IDs []string
}

// 测试
func TestUserRepo(t *testing.T) {
	row, err := NewUserRepoImpl().UpdateMap(context.TODO(), UserQuery{ID: "xxx"}, WithFinishTime(), WithStatus(1))
	if err != nil {
		panic(err)
	}
	fmt.Printf("本次更新行数 row =%d", row)
}

```