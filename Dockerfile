FROM ubuntu
ENV DEBIAN_FRONTEND noninteractive
WORKDIR ShortenedUrls
COPY . .
RUN apt -y update && \
    apt search golang-go && \
    apt search gccgo-go && \
    apt -y install golang-go \