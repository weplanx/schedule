FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD worker/main /app/

CMD [ "./main" ]
