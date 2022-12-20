package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
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

func (k *k8s) ExecPod(podName, containerName, namespace string, cmd []string) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	execOpt := &corev1.PodExecOptions{
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Container: containerName,
		Command:   cmd,
	}

	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		// 注意这里一定要import "k8s.io/client-go/kubernetes/scheme"，而不是"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
		VersionedParams(execOpt, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(K8s.Config, "POST", req.URL())
	if err != nil {
		fmt.Println("NewSPDYExecutor: " + err.Error())
		return stdout, stderr, errors.New("NewSPDYExecutor: " + err.Error())
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		fmt.Println("exec.Stream:" + err.Error())
		return stdout, stderr, errors.New("exec.Stream:" + err.Error())
	}

	//fmt.Println(stdout.String(), stderr.String())
	return
}

func test() {
	K8s.Init()
	stdout, stderr, err := K8s.ExecPod(
		"busybox-deployment-546c77d55b-p9nhh",
		"busybox",
		"default",
		[]string{"hostname"})
	if err != nil {
		return
	}

	fmt.Println(stdout.String(), stderr.String())
}
