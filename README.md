# Motivation
The goal of this project was for me to familiarized the technologies used, including:
- Golang  
- MySQL  
- Redis  
- Customized TCP
- Prometheus latency monitoring
- Connection pooling
  
# Design
Consists of user-facing web server and backend server. The backend server has access to database and cache, and interacts with the web server through connection pooling.
## Specifications
- Data size EstimationNumber of active users: 10,000,000
- Number of concurrent requests: 1000
- Size per user:
  - Username: 30 characters
  - Nickname: 30 characters
  - Password hash: 16 bytes
  - Token: 16 bytes
  - Token Datetime: 8 bytes
  - Image-url: 80 characters (to save space, we can store only the extension, rather than storing the whole url)
  - Total per user: < (220 bytes)
- Total size for metadata: 2.2GB

- Image file: 100kb
- Total size for images: 1TB

## MySQL table design
- Username: VARCHAR(30)
- Nickname: VARCHAR(30)
- Password hash: BINARY(16)
- Token:  BINARY(16)
- Token Datetime: DATETIME
- Image-url: VARCHAR(80)

## Packet format
Request packet from web server to backend:
{
"Format":
"User":{
"Username":
"Password":
"Nickname":
“Token”:
“URL”:
}
}

The packet format is specified in a file common to the backend and web server. With the use of json marshalling, this enables easy extension to the capabilities of a packet. Packets are of variable length, with end denoted by selected special character.

# What I would do if there were more time
- Implement the user-facing web page
- Add image storage server with AWS S3 integration
- Write unit tests
- Implmenent token staling based on login time
- Better error handling
- Limit number of incoming connections at backend server
- Gracefully handle packet transmission failure
  - Eg when client side timeouts, server might still try to respond on the connection
  