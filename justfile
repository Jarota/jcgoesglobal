default:
    just -l

build:
    podman build -t travelsapi --build-arg=ENVFILE=.env.local .

local:
    caddy start

    podman rm api

    podman run -it --name api \
    -v ./static/assets/pics:/root/static/assets/pics \
    -v ./app.db:/root/app.db \
    --publish 8080:8080 \
    travelsapi

deploy-fe:
    rsync -azP --delete --exclude=assets/pics ./static/ \
    ubuntu@83.228.199.74:/jcgoesglobal/static

deploy-api:
    rsync -azP \
    --include="main.go" --include="go.mod" \
    --include="go.sum" --include="internal/***" \
    --exclude="*" . ubuntu@83.228.199.74:/jcgoesglobal

    ssh ubuntu@83.228.199.74 'cd /jcgoesglobal && go build -o app main.go'

deploy: deploy-fe deploy-api
