docker stop audio-api
docker rm audio-api
docker images | grep audio-api | awk '{print $3}' | xargs -L1 docker rmi
docker load -i audio-api.$1
docker run -d -v /root/voice-data:/data -p 8580:8580 audio-api:$1
