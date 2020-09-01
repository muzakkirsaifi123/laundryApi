FROM golang:1.14
WORKDIR /go/src/github.com/jchenriquez
EXPOSE 9001
RUN go get github.com/go-chi/chi
RUN go get github.com/jackc/pgx
RUN go get github.com/spf13/viper
RUN go get github.com/spf13/cobra
RUN go get github.com/go-chi/cors
RUN go get github.com/go-chi/jwtauth
RUN mkdir laundromat
WORKDIR laundromat
COPY . .
RUN go install -v ./...
CMD ["laundromat", "start"]