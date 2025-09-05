default:
    just -l

local:
    caddy start
    dotenvx run -f .env.local -- go run main.go

deploy-static:
    rsync -azP --delete --exclude=assets/pics ./static/ \
    ubuntu@83.228.199.74:/jcgoesglobal/static

deploy-api:
    rsync -azP \
    --exclude=* \
    --include=*.go,go.mod,go.sum \
    . ubuntu@83.228.199.74:/jcgoesglobal


    ssh ubuntu@83.228.199.74 'cd /jcgoesglobal && go build -o app main.go'
