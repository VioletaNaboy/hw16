FROM golang

WORKDIR /bin

COPY . .

RUN go build -o api .

CMD [ "/bin/api" ]