module github.com/replicatedhq/outdated

go 1.15

require (
	github.com/docker/docker v23.0.0+incompatible
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/fatih/color v1.7.0
	github.com/genuinetools/reg v0.16.1
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.16.0 // indirect
	github.com/hashicorp/go-version v1.1.0
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/minio/minio v0.0.0-20190813204106-bf9b619d8656
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.4
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/tj/go-spin v1.1.0
	gotest.tools/v3 v3.4.0 // indirect
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab // indirect
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190815110238-8ff09bc626d6 // indirect
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.8.0
