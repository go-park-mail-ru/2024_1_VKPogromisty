#!/bin/bash

# Make the curl request and save the response headers to a file
curl -D headers.txt -X POST -H "Content-Type: application/json" \
     -H "Origin: http://localhost:3000" \
     -d '{"email":"petr09mitin@mail.ru","password":"admin1"}' \
     http://localhost:8080/api/v1/auth/login

# Parse the Set-Cookie header to extract the Session cookie
session_cookie=$(grep -i 'Set-Cookie: Session' headers.txt |  sed 's/Set-Cookie: //' | tr -d '\r')

# Make a new curl request to /api/v1/csrf/ with the Session cookie
curl -o csrf.txt -X GET -H "Content-Type: application/json" \
     -H "Cookie: ${session_cookie}" \
     -H "Origin: http://localhost:3000" \
     http://localhost:8080/api/v1/csrf/

# Parse the X-CSRF-Token header
csrf_token=$(grep -i '{"body":{"csrfToken":"' csrf.txt | sed -n 's|.*"csrfToken":"\([^"]*\)".*|\1|p')

echo "Session cookie: ${session_cookie}"
echo "CSRF token: ${csrf_token}"

wrk -t15 -c30 -d60s http://localhost:8080 \
    -H "X-CSRF-Token: ${csrf_token}" \
    -H "Cookie: ${session_cookie}" \
    -H "Origin: http://localhost:3000" \
    -s ./db/perf_test/get.lua
