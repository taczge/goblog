# mysql> source build.sql;
drop database blog;
create database blog;
use blog;
# 17 = 2015-01-02-123859
create table entry (id CHAR(17) primary key, title VARCHAR(200) not null, date DATE not null, body TEXT not null);

#load data local infile 'entry.csv' into table entry fields terminated by ',';

select * from entry;
