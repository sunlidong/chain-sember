sudo docker rm -f $(sudo docker ps -aq)
sudo docker network prune
sudo docker volume prune


