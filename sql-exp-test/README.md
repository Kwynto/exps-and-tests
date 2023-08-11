# sql-exp-test

This is an experimental repository.


## MySQL

`docker pull mysql/mysql-server:latest` 

`docker run -p 3306:3306 -d --name=mysql mysql/mysql-server:latest` 
 
`docker logs mysql` 

You need to find the generated password and copy it. 

`docker exec -it mysql bash` 

`mysql -uroot -p` 

`ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';` 

`update mysql.user set host='%' where user='root' and host='localhost';` 

`flush privileges;` 

And restart the docker container!
