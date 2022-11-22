FROM alpine:latest

RUN apk add --update-cache tzdata
COPY be-wallet /be-wallet

ENTRYPOINT ["/be-wallet"]


