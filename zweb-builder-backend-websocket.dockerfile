# ---------------------
# build zweb-builder-backend-websocket
FROM golang:1.19-bullseye as builder-for-backend

## set env
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

## build
WORKDIR /opt/zweb/zweb-builder-backend
RUN cd  /opt/zweb/zweb-builder-backend
RUN ls -alh

COPY ./ ./

RUN cat ./Makefile

RUN make build-websocket-server

RUN ls -alh ./bin/zweb-builder-backend-websocket


# -------------------
# build runner images
FROM alpine:latest as runner

WORKDIR /opt/zweb/zweb-builder-backend/bin/

## copy backend bin
COPY --from=builder-for-backend /opt/zweb/zweb-builder-backend/bin/zweb-builder-backend-websocket /opt/zweb/zweb-builder-backend/bin/


RUN ls -alh /opt/zweb/zweb-builder-backend/bin/



# run
EXPOSE 8002
CMD ["/bin/sh", "-c", "/opt/zweb/zweb-builder-backend/bin/zweb-builder-backend-websocket"]
