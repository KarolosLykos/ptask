<h1 align="center">Periodic Task</h1>
<p align="center">A JSON/HTTP service, in golang, that returns the matching timestamps of a periodic task.</p>

<p align="center">

<a style="text-decoration: none" href="https://github.com/KarolosLykos/ptask/actions?query=workflow%3AGo+branch%3Amain">
<img src="https://img.shields.io/github/actions/workflow/status/KarolosLykos/ptask/build.yml?style=flat-square" alt="Build Status">
</a>

<a style="text-decoration: none" href="go.mod">
<img src="https://img.shields.io/badge/Go-v1.19-blue?style=flat-square" alt="Go version">
</a>

<a href="https://codecov.io/gh/KarolosLykos/ptask" style="text-decoration: none">
<img src="https://img.shields.io/codecov/c/github/KarolosLykos/ptask?color=magenta&style=flat-square" alt="Downloads">
</a>


---

## Application Structure

---

The application is built following clean architecture design.
- The `cmd` folder contains the starting point of the application.
- The `internal` folder contains `interfaces` and `implementations` for interacting with the application.
    - The `logger` folder contains the `Logger` interface and the `logrus.Logger` implementation.
    - The `api` folder contains the REST API server using `gorilla` router.
    - The `ptask` folder contains all the interfaces, implementations and logic specific to the domain layer.
---

## Dependencies

- [github.com/gorilla/mux](https://github.com/gorilla/mux) HTTP router
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) Structured logger
- [github.com/golang/mock](https://github.com/golang/mock) Mocking framework
- [github.com/stretchr/testify](https://github.com/stretchr/testify) Testing Library
- [github.com/swaggo/swag](https://github.com/swaggo/swag) Swagger

## Run Instructions

- ### Run

  The service accepts as optional command-line arguments the listen `host`/`port` and `debug` flag.
  ```go
  go run cmd/main.go -host 0.0.0.0 -p 8080 -d
  ```

- ### Run with docker compose
  - `Dockerfile` is a multistage file that builds the application.

  - `docker-compose.yml` which runs the application deployment.
  ```
  docker compose up -d
  ```

Then point your `curl` commands at `http://localhost:8080` or open [Swagger](http://localhost:8080/swagger/index.html#/default/get_ptlist) on your browser.

## Run tests
```go
go test ./...
```

## Endpoints

---

### Ptlist

<details>

### Periodic task list

Example request:
```bash
curl -X GET http://localhost:8080/ptlist?period=1h&tz=America/Los_Angeles&t1=20210714T204603Z&t2=20210715T123456Z
```

Example Response:

200 Status OK
```
{ 
  "status":"success",
  "data":["20210714T210000Z","20210714T220000Z","20210714T230000Z","20210715T000000Z","20210715T010000Z","20210715T020000Z","20210715T030000Z","20210715T040000Z","20210715T050000Z","20210715T060000Z","20210715T070000Z","20210715T080000Z","20210715T090000Z","20210715T100000Z","20210715T110000Z","20210715T120000Z"]
}
```


400 Bad Request
```
{
  "status": "error",
  "error": "bad request"
}
```

500 Internal Server error
```
{
    "status": "error",
    "error": "something went wrong"
}
```

