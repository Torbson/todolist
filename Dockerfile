FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/git/todolist/
WORKDIR /go/src/git/todolist/
RUN go mod download
COPY . /go/src/git/todolist/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/todolist

FROM scratch
COPY --from=builder /go/src/git/todolist/bin/todolist /usr/bin/todolist
EXPOSE 8000 8000
ENTRYPOINT ["/usr/bin/todolist"]