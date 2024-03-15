package playstore

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

// New 创建并返回一个包含访问androidpublisher API所需凭证的http客户端。
//
// @Description: 通过提供的JSON密钥创建一个配置好的Client实例，可用于与Google Play Store API交互。
// @param jsonKey 用于构建JWT配置的JSON密钥字节切片。
// @return *Client 返回初始化好的Client实例。
// @return error 如果在创建过程中遇到任何错误，则返回非nil的error。
func New(jsonKey []byte) (*Client, error) {
	// 设置http客户端超时时间为10秒
	c := &http.Client{Timeout: 10 * time.Second}
	// 为context设置HTTP客户端，以便在OAuth2流程中使用
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	// 使用JSON密钥和所需范围配置JWT
	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	// 验证JWT配置是否正确，并获取访问令牌
	val := conf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, err
	}

	// 使用配置的HTTP客户端初始化androidpublisher服务
	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	// 返回初始化好的Client实例
	return &Client{service}, err
}

// NewWithClient returns http client which includes the custom http client.
// 使用自定义的http客户端创建并返回一个包含访问androidpublisher API所需凭证的http客户端。
func NewWithClient(jsonKey []byte, cli *http.Client) (*Client, error) {
	if cli == nil {
		return nil, fmt.Errorf("client is nil")
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cli)

	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{service}, err
}

// VerifySignature 验证应用内购买的签名。
// 您需要为您的 Android 应用的内购准备公钥，可在 https://play.google.com/apps/publish/ 上完成。
// 参数：
//
//	base64EncodedPublicKey string - 经过 Base64 编码的公钥字符串。
//	receipt []byte - 购买收据的字节数据。
//	signature string - 购买收据的签名字符串。
//
// 返回值：
//
//	isValid bool - 标识签名是否验证成功。
//	err error - 验证过程中遇到的错误。
func VerifySignature(base64EncodedPublicKey string, receipt []byte, signature string) (isValid bool, err error) {
	// 准备公钥
	decodedPublicKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key")
	}
	publicKey, _ := publicKeyInterface.(*rsa.PublicKey)

	// 从收据生成哈希值
	hasher := sha1.New()
	hasher.Write(receipt)
	hashedReceipt := hasher.Sum(nil)

	// 解码签名
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature")
	}

	// 验证签名
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature); err != nil {
		return false, nil
	}

	return true, nil
}
