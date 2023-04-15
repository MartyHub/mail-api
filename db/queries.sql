-- name: CreateMail :one
insert into mail (from_address, recipients_to, recipients_cc, subject, body, html, status, created_at, updated_at, try, retry_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning *
;

-- name: GetMail :one
select *
from mail
where id = $1
;

-- name: GetMailsToSend :many
select *
from mail
where status = 'QUEUED'
  and retry_at <= $1
order by created_at
limit $2 for update skip locked
;

-- name: UpdateMail :exec
update mail
set status     = $1,
    updated_at = $2,
    try        = try + 1,
    retry_at   = $3
where id = $4
;
