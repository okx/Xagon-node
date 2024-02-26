package nacos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// GetOneInstance is get the instance
func GetOneInstance(urls string, nameSpace string, param vo.SelectOneHealthInstanceParam) (instance *model.Instance, err error) {
	serverConfigs, err := getServerConfigs(urls)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve nacos server url %s: %s", urls, err.Error())
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig": constant.ClientConfig{
			NamespaceId:         nameSpace,
			TimeoutMs:           defaultTimeoutMs,
			NotLoadCacheAtStart: true,
			LogDir:              "/dev/null",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos client when getting one service. error: %s", err.Error())
	}

	instance, err = namingClient.SelectOneHealthyInstance(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s service in [%s, %s]. error: %s", param, urls, nameSpace, err.Error())
	}
	return instance, nil
}

func getServerConfigs(urls string) ([]constant.ServerConfig, error) {
	// nolint
	var configs []constant.ServerConfig
	for _, url := range strings.Split(urls, ",") {
		laddr := strings.Split(url, ":")
		serverPort, err := strconv.Atoi(laddr[1])
		if err != nil {
			return nil, err
		}
		configs = append(configs, constant.ServerConfig{
			IpAddr: laddr[0],
			Port:   uint64(serverPort),
		})
	}
	return configs, nil
}
