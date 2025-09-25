#!/bin/sh
set -e

cd /jcgoesglobal

podman load -i api.tar
podman stop api && podman rm api
podman run -d --name api \
  -v ./static/assets/pics:/root/static/assets/pics \
  -v ./app.db:/root/app.db \
  --publish 8080:8080 \
  jcapi
