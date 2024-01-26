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


create table todos (
                           id bigint not null auto_increment ,
                           user_id varchar(100) not null ,
    description text null,
                           title varchar(100) not null ,
                           created_at timestamp not null default current_timestamp,
                           updated_at timestamp not null default current_timestamp on update current_timestamp,
    deleted_at timestamp null,
                           primary key (id)
)engine = InnoDB;

select * from todos where id = 1;


select * from todos;


create table wallets (
    id varchar(100) not null ,
    user_id varchar(100) not null ,
    balance bigint not null ,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (user_id) references users(id)

) engine = InnoDB;


select * from wallets;

desc wallets;


create table addresses (
    id bigint not null auto_increment,
    user_id varchar(100) not null ,
    address varchar(100) not null ,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (user_id) references users(id)
) engine = InnoDB;

select * from addresses;


create table products (
    id varchar(100) not null,
    name varchar(100) not null ,
    price bigint not null ,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine = InnoDB;



create table user_like_product (
    user_id varchar(100) not null,
    product_id varchar(100) not null,
    primary key (user_id, product_id),
    foreign key (user_id) references users(id),
    foreign key (product_id) references products(id)
) engine = InnoDB;

desc user_like_product;

desc products;

select * from user_like_product;
select * from users;
select * from products;
select * from wallets;

drop table user_like_product;

show tables ;

desc wallets;







