FROM golang:latest as build

WORKDIR /go/src/app
COPY . .

RUN go build -a -tags netgo,osusergo -ldflags '-w -extldflags "-static"' -trimpath -o mailrg

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12

COPY --from=build /go/src/app/mailrg /

ENV SMTP_SERVER="smtp.gmail.com"
ENV SMTP_USER="example@gmail.com"
ENV SMTP_PASS="xyz"

ENV MAILRG_TOKEN="CorrectHorseBatteryStaple"
ENV DATA_DIR="/data"

EXPOSE 2163
VOLUME /data

ENTRYPOINT ["/mailrg"]