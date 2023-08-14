# sql-exp-test

This is an experimental repository.


## MySQL

```bash
docker pull mysql/mysql-server:latest
``` 

```bash
docker run -p 3306:3306 -d --name=mysql mysql/mysql-server:latest`
``` 
 
```bash
docker logs mysql
``` 

You need to find the generated password and copy it. 

```bash
docker exec -it mysql bash
``` 

```bash
mysql -uroot -p
``` 

```mysql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';
update mysql.user set host='%' where user='root' and host='localhost';
flush privileges;
```

And restart the docker container! 

```mysql
create database entities;
```

and 

```mysql
show databases;
```

This is end setings. 

You should run 
```bash
MYSQL_PASSWORD=root go run ./cmd/mysql/main.go
```


## PostgreSQL

docker run --name postgresql -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:latest


