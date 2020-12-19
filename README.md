# forest

基于`elton`的脚手架，实现了数据校验、行为统计等功能。


## commit

feat：新功能（feature）

fix：修补bug

docs：文档（documentation）

style： 格式（不影响代码运行的变动）

refactor：重构（即不是新增功能，也不是修改bug的代码变动）

test：增加测试

chore：构建过程或辅助工具的变动

## 启动数据库

### postgres

```
docker pull postgres:alpine

docker run -d --restart=always \
  -v $PWD/forest:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=A123456 \
  -p 5432:5432 \
  --name=forest \
  postgres:alpine

docker exec -it forest sh

psql -c "CREATE DATABASE forest;" -U postgres
psql -c "CREATE USER vicanso WITH PASSWORD 'A123456';" -U postgres
psql -c "GRANT ALL PRIVILEGES ON DATABASE forest to vicanso;" -U postgres
```

## redis

```
docker pull redis:alpine

docker run -d --restart=always \
  -p 6379:6379 \
  --name=redis \
  redis:alpine
```

## 规范

- 所有自定义的error都必须为hes.Error
