create table if not exists "users"(
    "id" integer not null primary key autoincrement,
    "username" varchar(128) not null,
    "pwd" varchar(128) not null
);

create table if not exists "accounts"(
    "id" integer not null primary key autoincrement,
    "account_tag" varchar(16) not null,
    "user_id" integer not null,
    "ip" varchar(15) not null,
    "port" integer not null,
    foreign key("user_id") references "users"("id")
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

create table if not exists "friend_messages"(
    "id" integer not null primary key autoincrement,
    "friend_id" integer not null,
    "message_id" integer not null,
    "message_ts" integer not null,
    "self_message" boolean not null,
    "read_mark" boolean not null default 0,
    "hide" boolean not null default 0,
    "revoke" boolean not null default 0,
    foreign key("friend_id") references "friends"("id")
);

create table if not exists "friend_sub_messages"(
    "id" integer not null primary key autoincrement,
    "friend_message_id" integer not null,
    "is_text" integer not null,
    "message" integer not null,
    foreign key("friend_message_id") references "friend_messages"("id")
);
