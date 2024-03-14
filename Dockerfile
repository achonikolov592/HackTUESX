FROM golang:1.22.1-alpine

WORKDIR /app
COPY ./ ./

RUN mkdir ./files

RUN apk add poppler-utils wv unrtf 

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o converter .

CMD ["./converter"]
EXPOSE 3334 