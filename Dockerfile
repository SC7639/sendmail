FROM golang

RUN mkdir -p /go/app/

WORKDIR /go/app
ADD sendmail.go .
ADD sendmail_test.go .
ADD test.sh .

# Make bash file executable
RUN chmod +x test.sh

CMD ["./test.sh"]
