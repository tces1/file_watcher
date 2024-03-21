FROM alpine
WORKDIR /root
COPY ./file_watcher /usr/local/bin/file_watcher
RUN chmod +x /usr/local/bin/file_watcher
COPY ./file_watcher.yaml /root/file_watcher.yaml
CMD ["/usr/local/bin/file_watcher"]
