FROM scratch

ADD demo-worker ./

CMD ["./demo-worker"]
