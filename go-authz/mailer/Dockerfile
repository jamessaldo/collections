FROM golang:alpine AS build

WORKDIR /go/src/app

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /build/mailer .

FROM scratch

COPY --from=build /build/mailer /mailer
COPY --from=build /go/src/app/env/app.env.production /env/app.env
COPY --from=build /go/src/app/env/storage.env.production /env/storage.env
COPY --from=build /go/src/app/env/mailer.env.production /env/mailer.env

ENTRYPOINT ["/mailer"]