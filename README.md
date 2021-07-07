# Motivation
The goal of this project was for me to familiarized the technologies used, including:
- Golang  
- MySQL  
- Redis  
- Customized TCP
- Prometheus latency monitoring
- Connection pooling

# What I would do if there were more time
- Implement the user-facing web page
- Add image storage server with AWS S3 integration
- Write unit tests
- Implmenent token staling based on login time
- Better error handling
- Limit number of incoming connections at backend server
- Gracefully handle packet transmission failure
  - Eg when client side timeouts, server might still try to respond on the connection
    
