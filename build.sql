drop database blog;
create database blog;
use blog;
create table entry (id INT auto_increment primary key, title VARCHAR(200) not null, date DATE not null, body TEXT not null);

load data local infile '/home/tn/git/vps/blog/entry.csv' into table entry fields terminated by ',';

select * from entry;
