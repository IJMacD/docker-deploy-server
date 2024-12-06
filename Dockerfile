FROM golang:1.23.1 AS build

WORKDIR /go

COPY go.mod ./
RUN go mod download -x

COPY src ./

ARG TARGETARCH
ENV CGO_ENABLED=0

RUN go build -v -o ./build/$TARGETARCH/docker-deploy-server .
COPY static ./build/$TARGETARCH/static/
COPY tmpl ./build/$TARGETARCH/tmpl/

FROM scratch AS final
ARG TARGETARCH

COPY --from=build /go/build/$TARGETARCH/ /

EXPOSE 8080
VOLUME "/fleets"
CMD ["/docker-deploy-server"]