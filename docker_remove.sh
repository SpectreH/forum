docker stop $(docker ps -l -a -q)
docker rm container
docker rmi web-image