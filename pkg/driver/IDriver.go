package driver

import "cert-gateway/bus/pb"

var (
	GRPC  = "GRPC"
	REDIS = "REDIS"
	ETCD  = "ETCD"
	AMQP  = "AMQP"
)

type OnMessageReceived = func(certs []*pb.Cert)
type OnErrReceiving = func(err error)

type IProvider interface {
	//GatewaySubscribe Gateway订阅
	GatewaySubscribe(localIp string, onMegReceived OnMessageReceived, onErrReceiving ...OnErrReceiving) error

	//SendCertificateToGateway 发送证书同步给某个 Gateway
	SendCertificateToGateway(localIP string) error

	// SyncCertificateToProvider Issuer发送证书给Provider
	SyncCertificateToProvider(certificateList *pb.CertificateList) error
}
