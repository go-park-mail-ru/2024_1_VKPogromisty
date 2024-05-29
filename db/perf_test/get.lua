-- Set the headers
wrk.headers["Content-Type"] = "application/json"

-- Set the method and path
wrk.method = "GET"
wrk.path = "/api/v1/posts/new"

-- Clear the body
wrk.body = nil

-- Handle the response
function response(status, headers, body)
end
