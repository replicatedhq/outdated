# `kubectl outdated`

A `kubectl` plugin to show out-of-date images running in a cluster.

## Quick Start

```
curl https://krew.sh/outdated | bash
kubectl outdated
```

The plugin will scan for all pods in all namespaces that you have at least read access to. It will then connect to the registry that hosts the image, and (if there's permission), it will analyze your tag to the list of current tags.

The output is a list of all images, with the most out-of-date images in red, slightly outdated in yellow, and up-to-date in green.

### Example

[![kubectl outdated](https://asciinema.org/a/ESRC5ubIylWMSQgyi015j04oa.svg)](https://asciinema.org/a/ESRC5ubIylWMSQgyi015j04oa)
