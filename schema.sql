create table tbl_users
(
	id int auto_increment primary key,
	business_name varchar(55) not null,
	full_name varchar(155) not null,
	business_email varchar(55) not null,
	business_phone varchar(55) not null,
	password varchar(55) not null,
	is_verified tinyint(1) default 0 not null,
	email_code varchar(55) not null,
	password_code varchar(55) not null
);