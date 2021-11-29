clear
docker image build -f Dockerfile -t web-image .
docker images
sleep 1
docker container run -p 8000:8000 --detach --name container web-image
docker ps -a
echo 'http://localhost:8000'