FROM golang:latest
WORKDIR /app/
ADD . /app/
#CMD ["chmod +x skeleton"]
CMD ["./skeleton"]
