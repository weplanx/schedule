FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD worker /app/

CMD [ "./worker" ]
