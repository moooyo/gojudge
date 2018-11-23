FROM scratch

# copy our static linked library

COPY judgeServer /home/gojudge/judgeServer/
COPY config.json /home/gojudge/judgeServer/
EXPOSE 12346/tcp

WORKDIR /home/gojudge/judgeServer
# tell we are exposing our service on port 8080

# run it!
CMD ["./judgeServer"]
