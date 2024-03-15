package pay

import (
	"context"

	"github.com/ayflying/pay/playstore"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx = gctx.New()
)

// GooglePay是一个处理Google支付的结构体。
type GooglePay struct {
	c *playstore.Client
}

// Init 初始化GooglePay客户端。
// data: 初始化客户端所需的配置数据。
func (p *GooglePay) Init(data []byte) {
	var err error
	p.c, err = playstore.New(data)
	if err != nil {
		panic(err) // 如果初始化失败，则panic。
	}
}

// VerifyPay 验证用户的支付。
// userId: 用户ID。
// OrderId: 订单ID。
// package1: 应用包名。
// subscriptionID: 订阅ID。
// purchaseToken: 购买凭证。
// cb: 验证结果的回调函数，如果验证成功，会调用此函数。
// 返回值: 执行错误。
func (p *GooglePay) VerifyPay(userId int64, OrderId, package1, subscriptionID, purchaseToken string, cb func(string, string) error) error {
	info, err := p.c.VerifyProduct(context.Background(), package1, subscriptionID, purchaseToken)
	if err != nil {
		return gerror.Cause(err) // 验证产品失败，返回错误。
	}
	if info.PurchaseState == 0 {
		if err := cb(subscriptionID, info.OrderId); err != nil {
			return gerror.Cause(err) // 调用回调函数失败，返回错误。
		}
	} else {
		return nil // 验证结果不为购买状态，直接返回nil。
	}
	return nil
}

// VerifyPayV1 是VerifyPay的另一个版本，用于验证订阅支付。
// package1: 应用包名。
// subscriptionID: 订阅ID。
// purchaseToken: 购买凭证。
// cb: 验证结果的回调函数。
// 返回值: 执行错误。
func (p *GooglePay) VerifyPayV1(package1, subscriptionID, purchaseToken string, cb func(string, string) error) error {
	g.Log().Infof(ctx, "VerifyPayV1: package = %v subscriptionID = %v, purchaseToken = %v", package1, subscriptionID, purchaseToken)
	info, err := p.c.VerifyProduct(context.Background(), package1, subscriptionID, purchaseToken)
	if err != nil {
		return gerror.Cause(err) // 验证产品失败，返回错误。
	}
	if info.PurchaseState == 0 {
		if err := cb(subscriptionID, info.OrderId); err != nil {
			return gerror.Cause(err) // 调用回调函数失败，返回错误。
		}
	} else {
		return nil // 验证结果不为购买状态，直接返回nil。
	}
	return nil
}

// VerifyPayV2 是VerifyPay的另一个版本，支持不同类型产品的验证。
// types: 验证的产品类型。
// package1: 应用包名。
// subscriptionID: 订阅ID。
// purchaseToken: 购买凭证。
// cb: 验证结果的回调函数。
// 返回值: 执行错误。
func (p *GooglePay) VerifyPayV2(types int32, package1, subscriptionID, purchaseToken string, cb func(string, string) error) error {
	g.Log().Infof(ctx, "VerifyPayV1: package = %v subscriptionID = %v, purchaseToken = %v", package1, subscriptionID, purchaseToken)
	switch types {
	case 0:
		info, err := p.c.VerifyProduct(context.Background(), package1, subscriptionID, purchaseToken)
		if err != nil {
			return gerror.Cause(err) // 验证产品失败，返回错误。
		}
		if info.PurchaseState == 0 {
			if err := cb(subscriptionID, info.OrderId); err != nil {
				return gerror.Cause(err) // 调用回调函数失败，返回错误。
			}
		}
	case 1:
		info, err := p.c.VerifySubscription(context.Background(), package1, subscriptionID, purchaseToken)
		if err != nil {
			return gerror.Cause(err) // 验证订阅失败，返回错误。
		}
		if len(info.OrderId) != 0 {
			if err := cb(subscriptionID, info.OrderId); err != nil {
				return gerror.Cause(err) // 调用回调函数失败，返回错误。
			}
		}
	}

	return nil
}

//func (p *GooglePay) VerifyPayTest(package1, subscriptionID, purchaseToken string) (*androidpublisher.ProductPurchase, error) {
//	return p.c.VerifyProduct(context.Background(), package1, subscriptionID, purchaseToken)
//}

func (p *GooglePay) VerifySubscriptionTest(package1, subscriptionID, purchaseToken string) (interface{}, error) {
	return p.c.VerifySubscription(context.Background(), package1, subscriptionID, purchaseToken)
}

// VerifySubSciption google 检查订阅是否有效
func (p *GooglePay) VerifySubSciption(package1, subscriptionID, purchaseToken string) (string, error) {
	info, err := p.c.VerifySubscription(context.Background(), package1, subscriptionID, purchaseToken)
	if err != nil {
		return "", gerror.Cause(err)
	}
	if len(info.OrderId) != 0 {
		return info.OrderId, nil
	}
	return "", nil
}
