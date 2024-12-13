package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mine/internal"
	"mine/internal/utils"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRedirect(ctx *fiber.Ctx) error {
	UrlRedirect := ctx.Query("urlRedirect", "")
	StatusStr := ctx.Query("status", "0")
	status, _ := strconv.Atoi(StatusStr)
	return ctx.Redirect(UrlRedirect, status)
}

func TestLog(ctx *fiber.Ctx) error {
	type Input struct {
		Loop int
		Text string
	}
	input := &Input{}
	err := ctx.BodyParser(input)
	if err != nil {
		return ctx.Status(200).JSON(err)
	}
	ListData := []string{}
	for i := 0; i < input.Loop; i++ {
		ListData = append(ListData, input.Text)
	}
	internal.Log.Info("TestLog", zap.Any("data", ListData))
	return nil
}

func TestLogKibana(ctx *fiber.Ctx) error {
	var temp interface{}
	type Input struct {
		TypeRequest string      `json:"type_request"`
		Topic       string      `json:"topic"`
		Async       bool        `json:"async"`
		Data        interface{} `json:"data"`
	}
	input := &Input{}
	err := ctx.BodyParser(input)
	if err != nil {
		return ctx.Status(200).JSON(err)
	}
	config := kafka.WriterConfig{
		Brokers: internal.Brokers, // Danh sách các Kafka broker
		Topic:   input.Topic,      // Tên topic muốn gửi message đến
		Async:   input.Async,
	}
	// Tạo một writer mới
	writer := kafka.NewWriter(config)
	byteInput, _ := json.Marshal(input.Data)
	if input.TypeRequest == "1" {
		// Gửi một message đến Kafka
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: byteInput,
			},
		)
		if err != nil {
			temp = fmt.Sprintf("Failed to write message: %v\n", err)
			return ctx.Status(200).JSON(temp)
		}
	} else {
		m := utils.KibanaMessageAll{}
		err := json.Unmarshal(byteInput, &m)
		if err != nil {
			return ctx.Status(200).JSON(err)
		}
		buffers := new(bytes.Buffer)
		json.NewEncoder(buffers).Encode(m)
		byteDate := buffers.Bytes()
		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: byteDate,
			},
		)
		if err != nil {
			temp = fmt.Sprintf("Failed to write message: %v\n", err)
			return ctx.Status(200).JSON(temp)
		}
	}

	// Đóng writer
	writer.Close()
	return ctx.Status(200).JSON(temp)
}

func TestAPI(ctx *fiber.Ctx) error {
	funcName := "TestAPI"
	type StructAPI struct {
		Url      string
		Header   map[string]string
		Body     map[string]interface{}
		IsGet    bool
		IsProxy  bool
		Param    map[string]string
		TypeCall string
	}
	body := &StructAPI{}
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.JSON(err)
	}
	internal.Log.Info("Call", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param))
	// if body.TypeCall == "NORMAL" {
	resp, err := utils.Request(body.Url, body.IsGet, body.Header, body.Param, body.Body, 30, body.IsProxy)
	if err != nil {
		internal.Log.Error("Response", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param), zap.Error(err))
		return ctx.JSON(err)
	}
	internal.Log.Info("Response", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param), zap.Any("resp", resp.String()), zap.Any("httpStatus", resp.Status()))
	res := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("Unmarshal", zap.Any("resp", resp.String()), zap.Error(err))
		return ctx.JSON(resp.String())
	}
	return ctx.JSON(res)
}

func testPod() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		return
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting kubernetes config: %v\n", err)
		return
	}
	// An empty string returns all namespaces
	namespace := "kube-system"
	pods, err := ListPods(namespace, clientset)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %v\n", pod.Name)
	}
	var message string
	if namespace == "" {
		message = "Total Pods in all namespaces"
	} else {
		message = fmt.Sprintf("Total Pods in namespace `%s`", namespace)
	}
	fmt.Printf("%s %d\n", message, len(pods.Items))

	// ListNamespaces function call returns a list of namespaces in the kubernetes cluster
	namespaces, err := ListNamespaces(clientset)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
	fmt.Printf("Total namespaces: %d\n", len(namespaces.Items))
}

func ListPods(namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	fmt.Println("Get Kubernetes Pods")
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return pods, nil
}

func ListNamespaces(client kubernetes.Interface) (*v1.NamespaceList, error) {
	fmt.Println("Get Kubernetes Namespaces")
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting namespaces: %v\n", err)
		return nil, err
	}
	return namespaces, nil
}
