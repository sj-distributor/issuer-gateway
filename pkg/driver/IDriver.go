package driver

var (
	GRPC  = "GRPC"
	REDIS = "REDIS"
	ETCD  = "ETCD"
	AMQP  = "AMQP"
)

type OnMessageReceived = func(msg string)
type OnErrReceiving = func(err error)

type IPubSubDriver interface {
	Subscribe(ip string, onMegReceived OnMessageReceived, onErrReceiving ...OnErrReceiving) error
	Publish(msg string) error
}
