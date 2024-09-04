
# ES

参考文档:
- es基础
  - https://www.cnblogs.com/jiujuan/p/17387600.html
- go-es
  - https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/examples.html

## 安装

- es
```bash
systemctl status firewalld.service
# 建立一些必要的文件夹，挂载到容器
mkdir -p /data/elasticsearch/config
mkdir -p /data/elasticsearch/data
mkdir -p /data/elasticsearch/plugins
# chmod 777 -R /data/elasticsearch
chmod 777 /data/elasticsearch/*
# 写入配置
echo "http.host: 0.0.0.0" >> /data/elasticsearch/config/elasticsearch.yml
# 创建容器，7.10.1即可，内存不需要给太大，服务太多虚拟机带不起来了
docker run -d --restart=always \
--name es --privileged \
-p 9200:9200 -p 9300:9300 \
-e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
-e "discovery.type=single-node" \
-e "xpack.security.enabled=false" \
elasticsearch:8.14.3



```
- kiabna
```bash

docker run -d --name kibana \
-e ELASTICSEARCH_HOSTS="http://192.168.109.128:9200" \
-p 5601:5601 kibana:7.10.1

```

