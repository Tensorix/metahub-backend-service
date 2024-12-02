create table if not exists "users"(
    "id" integer not null primary key autoincrement,
    "username" varchar(128) not null,
    "pwd" varchar(128) not null
);

create table if not exists "srvs"(
    "id" integer not null primary key autoincrement,
    "img_name" varchar(128) not null,
    "container" varchar(128) not null,
    "ip_addr" varchar(15) not null,
    "port" integer not null
);

create table if not exists "accounts"(
    "id" integer not null primary key autoincrement,
    "uid" integer not null,
    "account_tag" varchar(16) not null,
    "user_id" integer not null,
    "srv_id" integer not null,
    foreign key("user_id") references "users"("id"),
    foreign key("srv_id") references "srvs"("id")
);

create table if not exists "friends"(
    "id" integer not null primary key autoincrement,
    "account_id" integer not null,
    "nickname" varchar(128) not null,
    "uid" integer not null,
    "remark" varchar(128),
    "deleted" boolean not null default 0,
    foreign key("account_id") references accounts("id")
);

create table if not exists "friendmsgs"(
    "id" integer not null primary key autoincrement,
    "msg_id" integer not null,
    "friend_id" integer not null,
    "msg_ts" integer not null,
    "self_msg" boolean not null,
    "has_text" boolean not null,
    "has_img" boolean not null,
    "hide" boolean not null default 0,
    "revoke" boolean not null default 0,
    foreign key("friend_id") references "friends"("id")
);

create table if not exists "groups"(
    "id" integer not null primary key autoincrement,
    "account_id" integer not null,
    "gid" integer not null,
    "group_name" varchar(128) not null,
    "member_count" integer not null,
    "max_member_count" integer not null,
    "deleted" boolean not null default 0,
    foreign key("account_id") references "accounts"("id")
);

create table if not exists "gmembers"(
    "id" integer not null primary key autoincrement,
    "group_id" integer not null,
    "uid" integer not null,
    "nickname" varchar(128) not null,
    "isfriend" boolean not null
);

create table if not exists "groupmsgs"(
    "id" integer not null primary key autoincrement,
    "group_id" integer not null,
    "uid" integer not null,
    "msg_ts" integer not null,
    "has_text" boolean not null,
    "has_img" boolean not null,
    "hide" boolean not null default 0,
    "revoke" boolean not null default 0,
    foreign key("group_id") references "groups"("id")
);

create table if not exists "texts"(
    "id" integer not null primary key autoincrement,
    "pid" integer,
    "gid" integer,
    "msg_index" integer not null,
    "msg" text not null
);

create table if not exists "imgs"(
    "id" integer not null primary key autoincrement,
    "pid" integer,
    "gid" integer,
    "msg_index" integer not null,
    "file_name" varchar(32) not null
);
