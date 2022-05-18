docker build . -t audio-api:$1
docker save audio-api:$1 -o audio-api.$1
scp audio-api.$1 root@45.79.201.32:/root/docker-images