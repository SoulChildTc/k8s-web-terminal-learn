package main

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
}

func (k *k8s) Init() {
	var err error
	k.Config, err = clientcmd.BuildConfigFromFlags("", "../../kube-config")
	if err != nil {
		panic("kubeConfig解析失败: " + err.Error())
	}

	k.ClientSet, err = kubernetes.NewForConfig(k.Config)
	if err != nil {
		panic("kubernetes client创建失败: " + err.Error())
	} else {
		logger.Info("kubernetes client初始化成功")
	}
}
