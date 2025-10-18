ARG ENVFILE=".env.local"

FROM golang:1.25 as builder
WORKDIR /app

# Install dependencies and libvips
RUN apt-get update && apt-get install -y \
    --no-install-recommends \
    libvips-dev pkg-config build-essential \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o api .

FROM debian:bookworm-slim
WORKDIR /root

ARG ENVFILE
ENV ENVFILE=$ENVFILE

RUN apt-get update && apt-get install -y \
    --no-install-recommends \
    curl ca-certificates libvips \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sfS https://dotenvx.sh/install.sh | sh
COPY --from=builder /app/.env* ./

COPY --from=builder /app/api .
EXPOSE 8080

CMD ["/bin/sh", "-c", "dotenvx run -f $ENVFILE -- ./api"]
