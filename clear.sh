#!/usr/bin/env bash

### 删除容器  重要提示：记得用管理员权限

sudo docker rm -f $(sudo docker ps -aq)
sudo docker network prune
sudo docker volume prune


