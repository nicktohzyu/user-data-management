# Motivation
The goal of this project was for me to familiarize with the technologies used, including:
- Golang  
- MySQL  
- Redis  
- Customized TCP
- Prometheus latency monitoring
- Connection pooling
  
# Design
The project consists of user-facing web server and backend server. The backend server has access to database and cache, and interacts with the web server through connection pooling.
## Specifications
Data size Estimation
-  Number of active users: 10,000,000
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
&nbsp;  "Username":  
&nbsp;  "Password":  
&nbsp;  "Nickname":  
&nbsp;  “Token”:  
&nbsp;  “URL”:  
&nbsp;  }  
}

The packet format is specified in a file common to the backend and web server. Unused user fields are optional. With json marshalling, this enables easy extension to the capabilities of a packet. Packets are of variable length, with end denoted by selected special character.

# Testing
Tested with:
- 10,000,000 users
- up to 1000 workers
- 1000 requests per 0.2 seconds sent to rate limiter
- Effective 1000 requests per second
- 2048 maximum connections between backend server and database
Ran for 5+ minutes

Varying the number of connections in the pool, with 50 workers running the stress test, we get the following results (p99, ms):
| #pool connections | Round trip from client | Send message | Get connection | Free connection | DB get user metadata | DB update token |
|--|--|--|--|--|--|--|
| 50 | 9 | 0 | 0 | 0 | 5 | - |
| 100 | 9 | 0 | 0 | 0 | 6 | - |
| 200 | 9 | 0 | 0 | 0 | 6 | - |
| 400 | 8 | 0 | 0 | 0 | 4 | 5 |

Varying the number of connections in the pool, with 1000 workers running the stress test, we get the following results (p99, ms):
| #pool connections | Round trip from client | Send message | Get connection | Free connection | DB get user metadata | DB update token |
|--|--|--|--|--|--|--|
| 100 | 9 | 0 | 0 | 0 | 6 | 7 |
| 250 | 9 | 0 | 0 | 0 | 6 | 6 |
| 500 | 9 | 0 | 0 | 0 | 6 | 4 |
| 1000 | 9 | 0 | 0 | 0 | 6 | 3 |

Without the cache layer, varying from 100 to 1000 connections in the pool at 1000 QPS, p99 value for a round trip request remained between 8-9ms. This is an excellent latency level and exceeded the original goal of 200ms.

## Testing with Redis cache
With 1000 workers, 1000 connections, 1000 unique users, 60 second cache expiry  
| #pool connections | Round trip from client | Write to cache | Cache miss time* | Cache miss count | Cache hit time | Cache hit count | Cache miss rate|
|--|--|--|--|--|--|--|--|
| 1000 | 4 | 3 | 3 | 9.0e3 | 7 | 5.0e5 | 0.018 |
| 2000 | 4 | 6 | 10 | 1e4 | 2 | 1.2e6 | 0.0083 |
| 3000 | 3 | 6 | 9 | 9.0e3 | 2 | 1.5e6 | 0.006 |
| 4000 | 2 | 2 | 4 | 1.3e4 | 2 | 3.0e6 | 0.0043 |

\*Latency measurement for cache miss time includes a cache query and a subsequent DB query  
Number of cache misses reflects the expected value of #users*(test duration)/expiry

With Redis cache, testing from 1000 to 4000 QPS, round trip time p99 was reduced to under 4 ms. This demonstrates the effective implementation of the cache. However, it should be noted that the low cache miss rates are not representative of a real world scenario, as the queries per unique user in these tests are unnaturally large.

# Learning points
One issue I encountered early on was that round trip requests were taking in excess of 10 seconds. I added latency tracking (prometheus) to each section of the flow, and identified that the DB was the primary cause of the latency. I confirmed that the issue was caused by improperly setting up the table, such that no primary key had been set.

Another issue was that the backend server would stop responding after serving a specific number of requests. I investigated by first measuring this number, and looked in the code for matching values. I then tweaked those values to identify the variable that was correlated with this. It turned out that the matching variable was the number of connections in the DB pool. The cause of the error was that I did not close the DB rows after retrieving them, hence the DB connections remained waiting indefinitely, and the max number of connections was quickly reached.

# What I would do if there were more time
- Implement the user-facing web page
- Password hashing  
- Add image storage server with AWS S3 integration
- Write unit tests
- Implmenent token staling based on login time
- Better error handling
- Limit number of incoming connections at backend server
- Gracefully handle packet transmission failure
  - Eg when client side timeouts, server might still try to respond on the connection
- Testing with replicas (to estimate deviation)
  - Scripted testing
  
