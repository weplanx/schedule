FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD excel/main /app/

EXPOSE 9000

CMD [ "./main" ]
