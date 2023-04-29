package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"kubea-demo/config" // 引入配置

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

// 类似于: k8s = {"client":["Test1","Test2"],"config":["config1","config2"]}
type k8s struct {
	// 提供多集群服务client
	ClientMap map[string]*kubernetes.Clientset
	// 提供集群config列表
	KubeConfMap map[string]string
}

// 根据集群名获取client
func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {

	client, ok := k.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群 %v 不存在，无法获取集群信息", client))

	}
	return client, nil
}

// 初始化client
func (k *k8s) Init() {
	// 初始化map,申请内存空间
	mp := make(map[string]string, 0)
	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)

	// 反序列化配置
	if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
		panic(fmt.Sprintf("kubeconfig配置文件反序列化失败, err => %v", err))
	}

	k.KubeConfMap = mp

	// 初始化client
	for key, value := range mp {
		// 构造访问config的配置，从文件中加载
		conf, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			panic(fmt.Sprintf("集群 %s : 初始化配置失败 %v\n", key, err))
		}

		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群 %s : client创建失败 %v\n", key, err))
		}

		k.ClientMap[key] = clientSet
		logger.Info(fmt.Sprintf("集群 %s : client创建成功\n", key))

	}

}
