# DIF - Docker Image Forwarder
## TLDR; Usage

```sh
dif forward --src-host https://registry-1.docker.io/ --dst-host https://private-registry.my/ --dst-user registry-username --dst-pass very-secret-password alpine my-alpine latest 3.7
```
^-- will copy *alpine* image tags *latest* and *3.7* from Docker Hub to *my-alpine* repository in registry *https://private-registry.my/*

## Usage

### forward
```
NAME:
   dif forward - Forward docker image from source registry to destination registry

USAGE:
   dif forward [command options] [source repository] [destination repository] [tag a] [tag b] ... [tag n] - if no tags given, all tags will be forwarded

OPTIONS:
   --src-host value  Source registry hostname (default: "https://registry-1.docker.io/") [$DIF_SRC_HOST]
   --dst-host value  Destination registry hostname (default: "https://registry-1.docker.io/") [$DIF_DST_HOST]
   --src-user value  Source registry username [$DIF_SRC_USER]
   --dst-user value  Destination registry username [$DIF_DST_USER]
   --src-pass value  Source registry password [$DIF_SRC_PASS]
   --dst-pass value  Destination registry password [$DIF_DST_PASS]
   --help, -h        show help (default: false)
```

**You can use env vars instead of options - variable names defined in square brackets at the end of option description.**
