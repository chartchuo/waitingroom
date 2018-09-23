FROM ubuntu:18.04
WORKDIR /app
RUN mkdir /app/tmpl
ADD worker /app
ADD tmpl/wait.minify.tmpl /app/tmpl/wait.tmpl
ENTRYPOINT ["/app/worker"]
