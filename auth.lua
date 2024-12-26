wrk.method = "GET"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGF0X2lkIjoiMTE0MjUiLCJleHAiOjE3MzUwOTg5NTAsImlzX2hvdGVsaWVyIjp0cnVlLCJ1c2VyX2lkIjoxNCwidXNlcm5hbWUiOiJob3RlbGllciJ9.0VJfeZTa9y3JgbTihb9D4d8wWm549emvTu72otzHch4"

request = function()
    return wrk.format(nil, "/bookings/hotels?hotel_id=2")
end
