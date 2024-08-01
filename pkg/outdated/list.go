package outdated

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/pkg/wildcard"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type RunningImage struct {
	Namespace     string
	Pod           string
	InitContainer *string
	Container     *string
	Image         string
	PullableImage string
}

func (o Outdated) ListImages(ctx context.Context, configFlags *genericclioptions.ConfigFlags, imageNameCh chan string, ignoreNs []string) ([]RunningImage, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list namespaces")
	}

	runningImages := []RunningImage{}
	for _, namespace := range namespaces.Items {
		if isNamespaceExcluded(namespace.Name, ignoreNs) {
			continue
		}

		imageNameCh <- fmt.Sprintf("%s/", namespace.Name)

		pods, err := clientset.CoreV1().Pods(namespace.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list pods")
		}

		for _, pod := range pods.Items {
			for _, initContainerStatus := range pod.Status.InitContainerStatuses {
				pullable := initContainerStatus.ImageID
				if strings.HasPrefix(pullable, "docker-pullable://") {
					pullable = strings.TrimPrefix(pullable, "docker-pullable://")
				}
				runningImage := RunningImage{
					Pod:           pod.Name,
					Namespace:     pod.Namespace,
					InitContainer: &initContainerStatus.Name,
					Image:         initContainerStatus.Image,
					PullableImage: pullable,
				}

				imageNameCh <- fmt.Sprintf("%s/%s", namespace.Name, runningImage.Image)
				runningImages = append(runningImages, runningImage)
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				pullable := containerStatus.ImageID
				if strings.HasPrefix(pullable, "docker-pullable://") {
					pullable = strings.TrimPrefix(pullable, "docker-pullable://")
				}
				runningImage := RunningImage{
					Pod:           pod.Name,
					Namespace:     pod.Namespace,
					Container:     &containerStatus.Name,
					Image:         containerStatus.Image,
					PullableImage: pullable,
				}

				imageNameCh <- fmt.Sprintf("%s/%s", namespace.Name, runningImage.Image)
				runningImages = append(runningImages, runningImage)
			}
		}
	}

	// Remove exact duplicates
	cleanedImages := []RunningImage{}
	for _, runningImage := range runningImages {
		for _, cleanedImage := range cleanedImages {
			if cleanedImage.PullableImage == runningImage.PullableImage {
				goto NextImage
			}
		}

		cleanedImages = append(cleanedImages, runningImage)
	NextImage:
	}

	return cleanedImages, nil
}

func isNamespaceExcluded(namespace string, excluded []string) bool {
	for _, ex := range excluded {
		if wildcard.Match(ex, namespace) {
			return true
		}
	}

	return false
}
