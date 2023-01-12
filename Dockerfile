FROM alpine

RUN apk --update-cache upgrade
RUN apk add ffmpeg
