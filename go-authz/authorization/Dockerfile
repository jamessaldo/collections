FROM golang:alpine AS build

WORKDIR /go/src/app

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /build/authorization .

FROM scratch

COPY --from=build /build/authorization /authorization
COPY --from=build /go/src/app/env/app.env.production /env/app.env
COPY --from=build /go/src/app/env/storage.env.production /env/storage.env

# Expose port 8888 for the service to listen on
EXPOSE 8888


ENTRYPOINT ["/authorization"]