FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD excel /app/

EXPOSE 9000

CMD [ "./excel" ]
