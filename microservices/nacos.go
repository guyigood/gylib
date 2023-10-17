package microservices

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosIni struct {
	RegisterInstance struct {
		ServerIP    string `json:"server_ip"`
		ServicePort uint64 `json:"service_port"`
		ServiceName string `json:"service_name"`
		GroupName   string `json:"group_name"`
		ClusterName string `json:"cluster_name"`
	} `json:"register_instance"`
	ServerConfigs struct {
		IpAddr string `json:"ip_addr"`
		Port   uint64 `json:"port"`
	} `json:"server_configs"`
	ClientConfig struct {
		NamespaceId string `json:"namespace_id"`
		UserName    string `json:"user_name"`
		PassWord    string `json:"pass_word"`
	} `json:"client_config"`
}

type NacosClass struct {
	Ini        NacosIni
	IsClient   bool
	NameClient naming_client.INamingClient
}

func NewNacosServer() *NacosClass {
	that := new(NacosClass)
	that.Ini.RegisterInstance.ClusterName = "default"
	that.Ini.RegisterInstance.GroupName = "DEFAULT_GROUP"
	return that
}

func (that *NacosClass) CreatNacosClient() bool {
	// 创建 Nacos 客户端
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: that.Ini.ServerConfigs.IpAddr,
			Port:   that.Ini.ServerConfigs.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         that.Ini.ClientConfig.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		Username:            that.Ini.ClientConfig.UserName,
		Password:            that.Ini.ClientConfig.PassWord,
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}
	var err error
	that.NameClient, err = clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	that.IsClient = true
	return true
}
func (that *NacosClass) RegisterServer() bool {
	// 注册服务实例
	flag, err := that.NameClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          that.Ini.RegisterInstance.ServerIP,
		Port:        that.Ini.RegisterInstance.ServicePort,
		ServiceName: that.Ini.RegisterInstance.ServiceName,
		GroupName:   that.Ini.RegisterInstance.GroupName,
		Weight:      10,
		ClusterName: that.Ini.RegisterInstance.ClusterName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return flag
}

func (that *NacosClass) GetServerList(groupname string) (model.ServiceList, error) {
	instances, err := that.NameClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: that.Ini.ClientConfig.NamespaceId,
		GroupName: groupname,
	})
	if err != nil {
		//fmt.Println(err)
		return instances, err
	} else {
		return instances, nil
	}

}

func (that *NacosClass) GetIPList(servername, groupname string, cl_list []string) ([]model.Instance, error) {
	instances, err := that.NameClient.SelectAllInstances(vo.SelectAllInstancesParam{
		Clusters:    cl_list, // 默认值DEFAULT
		ServiceName: servername,
		GroupName:   groupname, // 默认值DEFAULT_GROUP
	})
	if err != nil {
		//fmt.Println(err)
		return instances, err
	} else {
		return instances, nil
	}

}

// cl_list=[]string{"DEFAULT"}
func (that *NacosClass) GetOneHealth(servername, groupname string, cl_list []string) *model.Instance {
	oneserver, _ := that.NameClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: servername,
		GroupName:   groupname, // 默认值DEFAULT_GROUP
		Clusters:    cl_list,   // 默认值DEFAULT
	})
	return oneserver
}

// cl_list=[]string{"DEFAULT"}
func (that *NacosClass) GetOneHealthDefault(servername string) *model.Instance {
	oneserver, _ := that.NameClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: servername,
	})
	return oneserver
}

func (that *NacosClass) GetIPListDefault(servername string) ([]model.Instance, error) {
	instances, err := that.NameClient.SelectAllInstances(vo.SelectAllInstancesParam{

		ServiceName: servername,
	})
	if err != nil {
		//fmt.Println(err)
		return instances, err
	} else {
		return instances, nil
	}

}
