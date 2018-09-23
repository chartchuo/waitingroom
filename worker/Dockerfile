FROM ubuntu:18.04
WORKDIR /app
RUN mkdir /app/tmpl
ADD worker /app
ADD tmpl/wait.tmpl /app/tmpl
ENTRYPOINT ["/app/worker"]
