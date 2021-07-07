create table users
(
    username varchar(30) not null,
    password varchar(30) null,
    nickname varchar(30) null,
    token    char(16)    null,
    constraint users_username_uindex
        unique (username)
);

alter table users
    add primary key (username);