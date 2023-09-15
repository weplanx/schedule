FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD main /app/

EXPOSE 9000

CMD [ "./main" ]
