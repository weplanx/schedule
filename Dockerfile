FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

ADD schedule /app/

CMD [ "./schedule" ]
