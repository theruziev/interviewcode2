create table contents
(
	id           bigserial,
	source       text,
	preview_url  text,
	preview_type text,
	content_type text,
	description  text,
	media_url    text,
	status       text,
	tags         text[],
	created_at   timestamp,
	updated_at   timestamp
);
