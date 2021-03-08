# `kubectl outdated`

`kubectl` `outdated` is a `kubectl` plugin that displays all out-of-date images running in a Kubernetes cluster.

## How it Works

The plugin will iterate through readable namespaces, and look for pods. For every pod it can read, the plugin will read the `podspec` for the container images, and any `init` container images. Additionally, it collects the content sha of the image, so that it can be used to disambiguate between different versions pushed with the same tag.

After collecting a list of images running and deduplicating this list, the plugin will anonymously connect to all required image repositories and request a list of tags. For tags and images that follow strict [semver](https://semver.org/) naming, the list is simply sorted and the plugin reports how out of date the running image is.

For images that aren't [semver](https://semver.org/) named, the plugin starts to collect tags dates from the manifest and sorts to find any tag that was pushed after the tag that is running.

## Quickstart

### Prerequisites

<mark>Note:</mark> You will need [git](https://git-scm.com/downloads) to install the `krew` plugin.

the `outdated` plugin is installed using the `krew` plugin manager for Kubernetes CLI. Installation instructions for `krew` can be found [here](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).

### Installation

After installing & configuring the k8s `krew` plugin, install `outdated` using the following command:

````
$ kubectl krew install outdated
````

### Usage

[![kubectl outdated](https://asciinema.org/a/ESRC5ubIylWMSQgyi015j04oa.svg)](https://asciinema.org/a/ESRC5ubIylWMSQgyi015j04oa)

The plugin will scan for all pods in all namespaces that you have at least read access to. It will then connect to the registry that hosts the image, and (if there's permission), it will analyze your tag to the list of current tags.

Scan all available images in your current `kubecontext` with the command:

````
kubectl outdated
````

The output is a list of all images, with the most out-of-date images in red, slightly outdated in yellow, and up-to-date in green.

### Contributing to `outdated`

Find a bug? Want to add a new feature? Want to write docs? Send a [pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests) & we'll review it! 