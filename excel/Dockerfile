FROM alpine:edge

RUN apk add tzdata

COPY ../dist/excel /app/main
WORKDIR /app

EXPOSE 9000

CMD [ "./main" ]
