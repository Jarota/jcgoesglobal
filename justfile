default:
    just -l

build:
    podman build -t jcapi:latest --build-arg=ENVFILE=.env .

local:
    caddy start

    podman rm api

    podman run -it --name api \
    -v ./static/assets/pics:/root/static/assets/pics \
    -v ./app.db:/root/app.db \
    --publish 8080:8080 \
    jcapi

deploy-fe:
    rsync -azP --delete --exclude=assets/pics ./static/ \
    ubuntu@83.228.199.74:/jcgoesglobal/static

deploy-api USER HOST:
    podman save -o api.tar jcapi:latest
    scp api.tar {{ USER }}@{{ HOST }}:/jcgoesglobal
    ssh {{ USER }}@{{ HOST }} 'bash -s' < ./remote-deploy.sh
