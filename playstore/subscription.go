package playstore

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
)

// AcknowledgeSubscription acknowledges a subscription purchase.
// 功能：确认订阅购买。
// 参数：packageName（应用包名），subscriptionID（订阅ID），token（购买令牌），req（确认请求对象）。
// 实现：使用PurchasesSubscriptionsService服务的Acknowledge方法来确认指定订阅。
func (c *Client) AcknowledgeSubscription(ctx context.Context, packageName string, subscriptionID string, token string,
	req *androidpublisher.SubscriptionPurchasesAcknowledgeRequest) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Acknowledge(packageName, subscriptionID, token, req).Context(ctx).Do()
	return err
}

// VerifySubscription verifies subscription status
// 功能：验证订阅状态。
// 参数：packageName（应用包名），subscriptionID（订阅ID），token（购买令牌）。
// 实现：使用PurchasesSubscriptionsService的Get方法来获取订阅的当前状态。
// 返回值：SubscriptionPurchase对象，包含订阅详情。
func (c *Client) VerifySubscription(ctx context.Context, packageName string, subscriptionID string, token string) (*androidpublisher.SubscriptionPurchase, error) {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	result, err := ps.Get(packageName, subscriptionID, token).Context(ctx).Do()
	return result, err
}

// CancelSubscription cancels a user's subscription purchase.
// 功能：取消用户的订阅购买。
// 参数：packageName（应用包名），subscriptionID（订阅ID），token（购买令牌）。
// 实现：使用PurchasesSubscriptionsService的Cancel方法来取消订阅。
func (c *Client) CancelSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Cancel(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RefundSubscription refunds a user's subscription purchase, but the subscription remains valid
// until its expiration time and it will continue to recur.
// 功能：退款用户的订阅购买，但订阅在到期前仍有效，并且会继续递延。
// 参数：packageName（应用包名），subscriptionID（订阅ID），token（购买令牌）。
// 实现：使用PurchasesSubscriptionsService的Refund方法来退款，但不取消订阅。
func (c *Client) RefundSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Refund(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RevokeSubscription refunds and immediately revokes a user's subscription purchase.
// Access to the subscription will be terminated immediately and it will stop recurring.
// 功能：退款并立即撤销用户的订阅购买。订阅将立即终止，并停止递延。
// 参数：packageName（应用包名），subscriptionID（订阅ID），token（购买令牌）。
// 实现：使用PurchasesSubscriptionsService的Revoke方法来退款并撤销订阅。
func (c *Client) RevokeSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Revoke(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}
