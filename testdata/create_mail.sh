#!/usr/bin/env sh

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

curl -XPOST http://localhost:8125/api/v1/mail \
  -H 'Content-Type: application/json' \
  -d @"${script_dir}/create_mail_input.json"
