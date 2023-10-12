FROM golang:latest

RUN apt-get update --allow-unauthenticated --allow-insecure-repositories
RUN apt-get install libasound-dev libportaudio2 libportaudiocpp0 portaudio19-dev -y

RUN mkdir /app
WORKDIR /app

ADD . /app

CMD go run .

koch             74585
