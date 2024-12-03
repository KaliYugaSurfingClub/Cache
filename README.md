# NoSQL stateful database

# Install
```cmd
git clone https://github.com/KaliYugaSurfingClub/Cache.git cache/src
cd cache/src
go build -o ../cache main.go
```

# The use
```cmd
cache -port=YOUR_PORT -logs_path=YOUR_FILE_FOR_STATE -time_for_shutdown=YOUR_TIME
```

# TCP API 
- 

# Go library
- 

# Rest light-wight API
## Get
- URL: `/v1/{key}`
- Method: `GET`
- Response variants: 
    - Body: `your requesting value`, StatusCode: `200`
    - Body: `no such key`, StatusCode `404`
    - StatusCode `500`

## Put (idempotent)
- URL: `/v1/{key}`
- Method: `PUT`
- Request Body: `your value to save` (simple text)
- Response variants:
  - StatusCode `201`
  - StatusCode `500`

## Delete (idempotent)
- URL: `/v1/{key}`
- Method: `DELETE`
- Response variants:
    - StatusCode `200`

## Clear (idempotent) delete all data
- URL: `/v1/operation/clear`
- Method: `DELETE`
- Response variants:
    - StatusCode `200`

