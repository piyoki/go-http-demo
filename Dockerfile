FROM golang:1.15-alpine as dev

WORKDIR /work

FROM golang:1.15-alpine as build

WORKDIR /app
COPY ./app/* /app/
RUN go build -o app

FROM alpine as runtime 
COPY --from=build /app/app /
COPY ./app/videos.json /
ENV SERVER_HOST=localhost:8080
CMD ./app

