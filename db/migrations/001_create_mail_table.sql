create type mail_status as enum ('QUEUED', 'SENT', 'ERROR');

create table mail
(
    id            int            not null generated always as identity,
    from_address  varchar(128)   not null,
    recipients_to varchar(128)[] not null,
    recipients_cc varchar(128)[],
    subject       varchar(64)    not null,
    body          text           not null,
    html          bool           not null,
    status        mail_status    not null,
    created_at    timestamp      not null,
    updated_at    timestamp      not null,
    try           smallint       not null,
    retry_at      timestamp      not null,
    constraint mail_pk primary key (id)
);

create index mail_status_ix on mail (status, updated_at) where status != 'SENT';
