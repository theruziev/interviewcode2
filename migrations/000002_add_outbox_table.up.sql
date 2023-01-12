create table outbox
(
	id         bigserial,
	topic      text,
	msg        jsonb,
	status     text,
	created_at timestamp,
	updated_at timestamp
);

create index outbox_created_at_status_topic_idx
	on outbox (created_at, status, topic);

