FROM golang:1.14.6-alpine3.12 AS builder
ENV ENV=there
ENV POSTGRES_USER=gweshgnedfhbja
ENV POSTGRES_PASSWORD=c68ccf60efdc5f5e8d1bb9cfe3635d49df9979cd7127c032b8921f7bf543744c
ENV POSTGRES_DB=d45omspb4mfdkf
ENV POSTGRES_HOST=ec2-54-78-36-245.eu-west-1.compute.amazonaws.com
ENV POSTGRES_PORT=5432
COPY go.mod go.sum /go/src/git/todolist/
WORKDIR /go/src/git/todolist/
RUN go mod download
COPY . /go/src/git/todolist/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/todolist

FROM scratch AS production
ENV ENV=there
ENV PORT=8000
ENV POSTGRES_USER=gweshgnedfhbja
ENV POSTGRES_PASSWORD=c68ccf60efdc5f5e8d1bb9cfe3635d49df9979cd7127c032b8921f7bf543744c
ENV POSTGRES_DB=d45omspb4mfdkf
ENV POSTGRES_HOST=ec2-54-78-36-245.eu-west-1.compute.amazonaws.com
ENV POSTGRES_PORT=5432
COPY --from=builder /go/src/git/todolist/bin/todolist /usr/bin/todolist
EXPOSE 8000 8000
ENTRYPOINT ["/usr/bin/todolist"]