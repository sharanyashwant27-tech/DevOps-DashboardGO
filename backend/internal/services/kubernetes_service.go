package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/devops-command-center/backend/config"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesService struct {
	client  *kubernetes.Clientset
	enabled bool
	log     *zap.Logger
}

func NewKubernetesService(cfg config.KubernetesConfig, log *zap.Logger) *KubernetesService {
	svc := &KubernetesService{enabled: cfg.Enabled, log: log}
	if !cfg.Enabled {
		return svc
	}
	var restCfg *rest.Config
	var err error
	if cfg.InCluster {
		restCfg, err = rest.InClusterConfig()
	} else {
		kubeconfig := cfg.Kubeconfig
		if kubeconfig == "" {
			home, _ := os.UserHomeDir()
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		restCfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		log.Warn("kubernetes config unavailable", zap.Error(err))
		svc.enabled = false
		return svc
	}
	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		log.Warn("kubernetes client unavailable", zap.Error(err))
		svc.enabled = false
		return svc
	}
	svc.client = clientset
	return svc
}

func (s *KubernetesService) Enabled() bool { return s.enabled && s.client != nil }

func (s *KubernetesService) ensure() error {
	if !s.Enabled() {
		return fmt.Errorf("kubernetes is not available")
	}
	return nil
}

func (s *KubernetesService) ListNamespaces(ctx context.Context) ([]corev1.Namespace, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListReplicaSets(ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListDaemonSets(ctx context.Context, namespace string) ([]appsv1.DaemonSet, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListServices(ctx context.Context, namespace string) ([]corev1.Service, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListIngresses(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListNodes(ctx context.Context) ([]corev1.Node, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListPVs(ctx context.Context) ([]corev1.PersistentVolume, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListPVCs(ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) ListEvents(ctx context.Context, namespace string) ([]corev1.Event, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	list, err := s.client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *KubernetesService) PodLogs(ctx context.Context, namespace, pod, container string, tail int64) (string, error) {
	if err := s.ensure(); err != nil {
		return "", err
	}
	opts := &corev1.PodLogOptions{TailLines: &tail}
	if container != "" {
		opts.Container = container
	}
	req := s.client.CoreV1().Pods(namespace).GetLogs(pod, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	b, err := io.ReadAll(stream)
	return string(b), err
}

func (s *KubernetesService) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	if err := s.ensure(); err != nil {
		return err
	}
	scale, err := s.client.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = replicas
	_, err = s.client.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	return err
}

func (s *KubernetesService) RestartDeployment(ctx context.Context, namespace, name string) error {
	if err := s.ensure(); err != nil {
		return err
	}
	patch := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, metav1.Now().UTC().Format("2006-01-02T15:04:05Z"))
	_, err := s.client.AppsV1().Deployments(namespace).Patch(
		ctx, name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{},
	)
	return err
}

func (s *KubernetesService) DeletePod(ctx context.Context, namespace, name string) error {
	if err := s.ensure(); err != nil {
		return err
	}
	return s.client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
