# download qwen model
FROM gerke74/ollama-model-loader as downloader
RUN /ollama-pull mistral:latest


# build go app
FROM docker.io/library/golang:1.22 as builder

COPY service/go.mod src/go.mod
COPY service/go.sum src/go.sum
COPY service/pkg src/pkg/
COPY service/main.go src/main.go

RUN --mount=type=cache,target=/root/.cache/go-build cd src && \ 
        CGO_ENABLED=false go build -tags osusergo,netgo -o /service *.go;

# final app
FROM ubuntu:latest

RUN apt-get update && apt-get install -y supervisor curl
RUN curl -L https://ollama.com/download/ollama-linux-amd64 -o /usr/bin/ollama
RUN chmod +x /usr/bin/ollama

COPY --from=downloader /root/.ollama /root/.ollama
COPY --from=builder  /service /service

COPY supervisord.conf /etc/supervisor/conf.d/

CMD ["/usr/bin/supervisord"]
