create table users
(
	id                  bigserial,
	public_id           uuid,
	first_name          text,
	last_name           text,
	email               text,
	password            text,
	status              text,
	created_at          timestamp,
	updated_at          timestamp,
	activation_code     text,
	reset_password_code text,
	otp_secret          text,
	otp_recovery_codes  int[],
	otp_enabled         bool default false
);

create unique index public_id_uidx
	on users (public_id);

create unique index users_email_uidx
	on users (email);


