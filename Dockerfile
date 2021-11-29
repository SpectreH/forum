FROM golang:1.16-alpine
LABEL maintainer="SpectreH"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.description="Forum in docker"

RUN apk update && apk add --no-cache gcc
RUN apk add --update gcc musl-dev

RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/
COPY . /usr/src/app/

RUN go build -o /forum
EXPOSE 8000

CMD [ "/forum" ]
