
keys
```bash
> keys *
bluebell:post:time
bluebell:post:score
bluebell:post:voted:post1
bluebell:post:voted:post2
bluebell:community:1

> type bluebell:post:voted:post1
zset

> zrange bluebell:post:voted:post1 0 -1 withscores
user1
1
user2
1

> ZCOUNT bluebell:post:voted:post1 1 1
2
```