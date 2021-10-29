package register

import (
	"fmt"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/event"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/middleware/logger"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

const (
	HTTPService = iota
	RPCService
)

var Client *api.Client

func Init(serviceType int, serviceName string, serviceIp string, servicePort int) {
	//创建配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", serviceIp, servicePort)
	//创建客户端
	Client, _ = api.NewClient(cfg)
	//创建服务ID
	id := fmt.Sprintf("%s", uuid.NewV4())
	//服务注册
	ServiceRegister(serviceType, id, serviceName, serviceIp, servicePort)
	//事件句柄
	event.RegisterHandler(event.ServiceTerm, func() {
		err := ServiceDeregister(id)
		if err != nil {
			logger.Global.Info(serviceName+"服务注销失败：", err)
		} else {
			logger.Global.Info(serviceName + "服务注销成功")
		}
	})
}

func ServiceRegister(serviceType int, id string, name string, ip string, port int) {

	var check *api.AgentServiceCheck

	switch serviceType {
	case RPCService:
		check = &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", ip, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	case HTTPService:
		fallthrough
	default:
		check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", ip, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	}

	reg := new(api.AgentServiceRegistration)
	reg.ID = id
	reg.Name = name
	reg.Address = ip
	reg.Port = port
	reg.Tags = []string{name}

	reg.Check = check

	err := Client.Agent().ServiceRegister(reg)

	if err != nil {
		panic(err)
	}
}

func ServiceDeregister(id string) error {
	return Client.Agent().ServiceDeregister(id)
}
