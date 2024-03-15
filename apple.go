package pay

import (
	"context"
	"github.com/go-pay/gopay/apple"
	"github.com/gogf/gf/v2/errors/gerror"
	"strings"
	"sync"
	"time"
)

// ApplePay 苹果支付
// 这是一个用于处理苹果支付的结构体。
type ApplePay struct {
	pass string       // pass 是用于苹果支付过程中的密钥。
	lock sync.RWMutex // lock 用于确保在并发访问或修改 pass 时的安全性。
}

// Init 是ApplePay类型的初始化函数。
//
// @Description: 对ApplePay对象进行初始化，将传入的数据存储到对象中。
// @receiver p: ApplePay对象的指针，用于接收初始化操作。
// @param data: 一个字节切片，包含需要初始化的数据。
func (p *ApplePay) Init(data []byte) {
	p.lock.Lock()         // 加锁以保证在多线程环境下的线程安全
	defer p.lock.Unlock() // 确保在函数执行完毕退出时自动解锁，避免死锁
	p.pass = string(data) // 将传入的字节切片数据转换为字符串，并赋值给pass字段
}

// VerifyPay 验证苹果支付
//
// @Description: 验证苹果支付的收据信息，以确认支付的有效性。
// @receiver p *ApplePay: ApplePay对象，用于执行验证支付的操作。
// @param userId uint64: 用户ID。
// @param OrderId string: 订单ID。
// @param package1 string: 付费产品的包装名称。
// @param subscriptionID string: 订阅ID。
// @param purchaseToken string: 购买令牌，用于苹果服务器的收据验证。
// @param isDebug bool: 是否为调试模式，决定使用哪个验证URL。
// @param cb func(string) error: 回调函数，用于处理验证成功后的产品ID。
// @return error: 返回错误信息，如果验证过程中出现错误，则返回相应的错误信息。
func (p *ApplePay) VerifyPay(userId uint64, OrderId, package1, subscriptionID, purchaseToken string, isDebug bool, cb func(string) error) error {
	p.lock.RLock()         // 加读锁，保证并发安全
	defer p.lock.RUnlock() // 解读锁，确保函数执行完毕后释放锁
	// 根据是否为调试模式选择验证URL
	url := apple.UrlProd
	if isDebug {
		url = apple.UrlSandbox
	}
	// 向苹果服务器验证收据
	info, err := apple.VerifyReceipt(context.Background(), url, p.pass, purchaseToken)
	if err != nil {
		// 如果验证失败，则返回错误
		return err
	}
	// 检查收据验证的状态
	if info.Status == 0 {
		// 检查收据中是否包含内购信息
		if len(info.Receipt.InApp) <= 0 {
			return gerror.Wrap(err, "info.Receipt.InApp = 0")
		}
		// 调用回调函数处理商品ID
		if err := cb(info.Receipt.InApp[0].ProductId); err != nil {
			// 如果回调处理失败，则返回错误
			return err
		}
	} else {
		// 如果收据验证状态异常，则返回状态错误信息
		return gerror.Wrapf(err, "status err = %v", info.Status)
	}
	return nil
}

// VerifyPayV1 验证苹果支付的交易
//
//	@Description:
//	@receiver p
//	@param purchaseToken
//	@param isDebug
//	@param cb
//	@return error
func (p *ApplePay) VerifyPayV1(purchaseToken string, isDebug bool, cb func(string, string) error) error {
	p.lock.RLock()         // 加读锁，确保并发安全
	defer p.lock.RUnlock() // 结束时自动释放读锁
	// 根据调试模式选择验证服务的URL
	url := apple.UrlProd
	if isDebug {
		url = apple.UrlSandbox
	}
	// 向苹果服务器验证收据
	info, err := apple.VerifyReceipt(context.Background(), url, p.pass, purchaseToken)
	if err != nil {
		// 验证失败，返回错误
		return err
	}
	// 检查验证结果状态
	if info.Status == 0 {
		// 验证成功，检查收据中是否有内购信息
		if len(info.Receipt.InApp) <= 0 {
			// 收据中无内购信息，返回错误
			return gerror.Wrap(err, "info.Receipt.InApp = 0")
		}
		// 调用回调函数，处理内购产品信息
		if err := cb(info.Receipt.InApp[0].ProductId, info.Receipt.InApp[0].OriginalTransactionId); err != nil {
			// 回调函数执行失败，返回错误
			return gerror.Wrap(err, "回调函数执行失败")
		}
	} else {
		// 验证结果状态异常，返回错误
		return gerror.Wrapf(err, "status err = %v", info.Status)
	}
	// 验证成功，返回nil
	return nil
}

// VerifyPayTest 用于验证苹果支付的测试购买。
//
//	@Description:
//	@receiver p
//	@param purchaseToken
//	@return interface{}
//	@return error
func (p *ApplePay) VerifyPayTest(purchaseToken string) (interface{}, error) {
	// 使用沙箱环境的URL进行验证
	url := apple.UrlSandbox
	// 调用apple.VerifyReceipt进行收据验证
	return apple.VerifyReceipt(context.Background(), url, p.pass, purchaseToken)
}

// GetTime 根据提供的 timer 字符串解析时间，格式为 "YYYY-MM-DD HH:MM:SS ZZZ"，若解析失败则返回当前时间
//
//	@Description: 根据指定格式解析时间字符串，如果解析失败或者格式不正确，则返回当前时间。
//	@param timer 时间字符串，格式为 "YYYY-MM-DD HH:MM:SS ZZZ"，其中 ZZZ 为时区标识。
//	@return time.Time 解析得到的时间，若失败则返回当前时间。
func GetTime(timer string) time.Time {
	// 将 timer 字符串按空格分割为年月日和时分秒两部分
	ts := strings.Split(timer, "")
	// 如果分割后的数组长度不为3，则说明格式不正确，返回当前时间
	if len(ts) != 3 {
		return time.Now()
	}
	// 尝试加载指定时区信息
	location, err := time.LoadLocation(ts[2])
	// 如果加载时区失败，则返回当前时间
	if err != nil {
		return time.Now()
	}
	// 使用指定时区解析时间字符串
	t, err := time.ParseInLocation("2006-01-02 15:04:05 MST", ts[0]+" "+ts[1], location)
	// 如果解析失败，则返回当前时间
	if err != nil {
		return time.Now()
	}
	// 将解析得到的时间转换为本地时间并返回
	return t.In(time.Local)
}
