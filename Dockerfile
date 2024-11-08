FROM golang:1.23 as build

ENV CGO_ENABLED=0
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

RUN go build /app/cmd/main/tufin.go

RUN curl -sfL https://get.k3s.io | sh -s - --write-kubeconfig-mode 644 --disable-agent

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/tufin /tufin
COPY --from=build /usr/local/bin/k3s /usr/local/bin/k3s

ENTRYPOINT ["/tufin"]
