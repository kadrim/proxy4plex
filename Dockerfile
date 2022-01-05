FROM golang:1.17-alpine

LABEL maintainer="Kadrim <kadrim@users.noreply.github.com>"

WORKDIR /proxy4plex
ADD . .

RUN go build

EXPOSE 80
EXPOSE 3000

CMD [ "./proxy4plex" ]