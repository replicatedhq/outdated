kubectl-outdated finds and reports on any outdated images running in a Kubernetes cluster.

## Usage
The following assumes you have kubectl-outdated installed via

```shell
kubectl krew install outdated
```

### Scan images in your current kubecontext

```shell
kubectl outdated
```

### Scan images in another kubecontext

```shell
kubectl outdated --kubecconfig=/path/to/kubeconfig
```

## How it works
The plugin will iterate through readable namespaces, and look for pods. For every pod it can read, the plugin will read the podspec for the container images, and any init container images. Additionally, it collects the content sha of the image, so that it can be used to disambiguate between different versions pushed with the same tag.

After collecting a list of images running and deduplicating this list, the plugin will anonymously connect to all required image repositories and request a list of tags. For tags and images that follow strict semver naming, the list is simply sorted and the plugin reports how out of date the running image is.

For images that aren't semver named, the plugin starts to collect tags dates from the manifest and sorts to find any tag that was pushed after the tag that is running.

