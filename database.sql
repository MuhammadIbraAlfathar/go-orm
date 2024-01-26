create table users (
    id varchar(100) not null ,
    password varchar(100) not null ,
    name varchar(100) not null ,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
)engine = InnoDB;

select * from users;

alter table users change name first_name varchar(100);

alter table users add column middle_name varchar(100) null after first_name;
alter table users add column last_name varchar(100) null after middle_name;


create table user_logs (
                       id int auto_increment ,
                       user_id varchar(100) not null ,
                       action varchar(100) not null ,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp on update current_timestamp,
                       primary key (id)
)engine = InnoDB;

select * from user_logs;

delete from user_logs;

alter table user_logs modify created_at bigint not null;
alter table user_logs modify updated_at bigint not null;

desc user_logs;