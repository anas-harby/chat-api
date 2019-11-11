from golang:1.12
RUN go get github.com/go-redis/redis github.com/bsm/redislock github.com/gorilla/mux github.com/imroc/req
WORKDIR /
RUN mkdir /go/src/github.com/anas-harby/
RUN mkdir /go/src/github.com/anas-harby/go-chat-creation-api

COPY . /go/src/github.com/anas-harby/go-chat-creation-api
COPY entrypoint.sh /usr/bin/entrypoint-go.sh
RUN chmod +x /usr/bin/entrypoint-go.sh
COPY wait-for-it.sh /usr/bin
RUN chmod +x /usr/bin/wait-for-it.sh
ENTRYPOINT ["entrypoint-go.sh"]
EXPOSE 8080

CMD ["/go/src/github.com/anas-harby/go-chat-creation-api/cmd/go-chat-creation-api/go-chat-creation-api"]
