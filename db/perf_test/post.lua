-- Generate a random boundary
local boundary = '----WebKitFormBoundary' .. math.random(1000, 9999)

-- Set the headers
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary

-- Set the method and path
wrk.method = "POST"
wrk.path = "/api/v1/posts/"

-- Construct the body
local body = '--' .. boundary .. '\r\n'
           .. 'Content-Disposition: form-data; name="content"\r\n\r\n'
           .. 'post content' .. '\r\n'
           .. '--' .. boundary .. '--\r\n'
wrk.body = body

-- Handle the response
function response(status, headers, body)
end
