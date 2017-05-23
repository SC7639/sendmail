FROM golang

# Install sendmail
RUN apt-get update && apt-get install sendmail -y

# Set FQDN
RUN line=$(head -n 1 /etc/hosts) && line2=$(echo $line | awk '{print $2}') && echo "$line $line2.localdomain" >> /etc/hosts

RUN mkdir -p /go/app/

WORKDIR /go/app
ADD main.go .
ADD sendmail_test.go .
ADD test.sh .

# Make bash file executable
RUN chmod +x test.sh

CMD ["./test.sh"]
