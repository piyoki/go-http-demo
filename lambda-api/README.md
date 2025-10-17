# Go Lambda-like HTTP Server

## How to Run

```shell
go run main.go
```

Test with curl:

```shell
curl -X POST http://localhost:8080/events/new -v
```

You’ll see a 204 No Content response, and the server will terminate right after handling the request.
