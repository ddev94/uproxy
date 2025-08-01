FROM golang:1.23-alpine AS dev

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o uproxy main.go

FROM golang:1.23-alpine AS prod

WORKDIR /app

LABEL traefik.enable=true
LABEL traefik.http.routers.api.rule=PathPrefix(`/apis`)
LABEL traefik.http.services.api.loadbalancer.server.port=8080

COPY --from=dev /app/uproxy /app

EXPOSE 8080

ENV PORT 8080

ENTRYPOINT ["/app/uproxy", "--target-url", "http://ad1b13f3cd5b04c39a31f5a26197bc6f-493533987.ap-northeast-3.elb.amazonaws.com:7199"]