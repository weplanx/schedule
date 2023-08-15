FROM alpine:edge

RUN apk add tzdata

COPY main /app/main
WORKDIR /app

EXPOSE 9000

CMD [ "./main" ]
