FROM scratch

# copy our static linked library

COPY judgeServer /home/gojudge/judgeServer/
COPY config.json /home/gojudge/judgeServer/

EXPOSE 8080/tcp
EXPOSE 8081/tcp

WORKDIR /home/gojudge/judgeServer
# tell we are exposing our service on port 8080

# run it!
CMD ["./judgeServer"]
