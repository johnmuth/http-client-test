FROM golang:1.8

RUN go get -u github.com/alecthomas/gometalinter && gometalinter --install

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]
