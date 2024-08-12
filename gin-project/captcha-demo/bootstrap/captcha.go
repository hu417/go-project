package bootstrap

import (
	"image/color"
	"sync"

	"captcha-demo/utils"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

var (
	// 参数调整参考: https://captcha.mojotv.cn/.netlify/functions/captcha
	// 图片高度
	height = 80
	// 图片宽度
	width = 240
	// 干扰线选项
	showLineOptions = 4 // 2 4 8
	// 干扰数(新增干扰的字符)
	noiseCount = 2
	// 验证码长度
	length = 4
	// 字体
	fonts = []string{"wqy-microhei.ttc"}
	// 字符
	stringSource = "1234567890QAZWSXCDERFVBGTYHNMJUIKLOP"
	// 中文
	chineseSource = "你,好,真,香,哈,哈,消,费,者,狗,仔,北京烤鸭"
	// 音频语言
	language = "zh"
	// 最大绝对偏斜系数为个位数
	// maxSkew = 0.7
	// 圆圈的数量
	// dotCount = 60

	once sync.Once
)

//创建store,保存验证码的位置,默认为mem(内存中)单机部署,如果要布置多台服务器,则可以设置保存在redis中
//var store = base64Captcha.DefaultMemStore

func InitCaptcha(captchatype string, r *redis.Client) (captcha *base64Captcha.Captcha) {
	// 配置RedisStore, RedisStore实现base64Captcha.Store接口
	var store base64Captcha.Store = utils.NewCaptchaStore(r)

	once.Do(func() {

		switch {
		// 字符串
		case captchatype == "string":
			// 字符串验证码
			driver := &base64Captcha.DriverString{
				Height:          height,
				Width:           width,
				ShowLineOptions: showLineOptions,
				NoiseCount:      noiseCount,
				Source:          stringSource,
				Length:          length,
				BgColor: &color.RGBA{ // 背景颜色
					R: 3,
					G: 102,
					B: 214,
					A: 125,
				},
				Fonts: fonts,
			}
			driver = driver.ConvertFonts()
			captcha = base64Captcha.NewCaptcha(driver, store)
		// 算术
		case captchatype == "math":
			driver := &base64Captcha.DriverMath{
				Height:          height,
				Width:           width,
				NoiseCount:      noiseCount,
				ShowLineOptions: showLineOptions,
				Fonts:           fonts,
			}
			driver = driver.ConvertFonts()
			captcha = base64Captcha.NewCaptcha(driver, store)
		// 中文
		case captchatype == "cn":
			driver := &base64Captcha.DriverChinese{
				Height:          height,
				Width:           width,
				NoiseCount:      noiseCount,
				ShowLineOptions: showLineOptions,
				Length:          length,
				Source:          chineseSource,
				Fonts:           []string{"wqy-microhei.ttc"},
			}
			driver = driver.ConvertFonts()
			captcha = base64Captcha.NewCaptcha(driver, store)
		// 音频
		case captchatype == "":
			driver := &base64Captcha.DriverAudio{
				Length: length,
				// "en", "ja", "ru", "zh".
				Language: language,
			}
			captcha = base64Captcha.NewCaptcha(driver, store)
		// 数字
		case captchatype == "":
			driver := &base64Captcha.DriverDigit{
				Height: height,
				Width:  width,
				Length: length,
				// 最大绝对偏斜系数为个位数
				MaxSkew: 0,
				// 背景圆圈的数量。
				DotCount: 0,
			}
			captcha = base64Captcha.NewCaptcha(driver, store)
		default:
			// 字符串验证码
			driver := &base64Captcha.DriverString{
				Height:          height,
				Width:           width,
				ShowLineOptions: showLineOptions,
				NoiseCount:      noiseCount,
				Source:          stringSource,
				Length:          length,
				Fonts:           fonts,
			}
			driver = driver.ConvertFonts()
			captcha = base64Captcha.NewCaptcha(driver, store)
		}

	})
	return captcha
}
