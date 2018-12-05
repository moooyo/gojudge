#! /bin/bash

docker volume rm gojudgetest

docker rm gojudgetest_helper

docker volume create gojudgetest


docker run --name=gojudgetest_helper --mount type=volume,source=gojudgetest,target=/gojudgetest  busybox

docker cp . gojudgetest_helper:/gojudgetest/

docker rm gojudgetest_helper
