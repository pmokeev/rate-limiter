# Rate Limiter

An HTTP service capable of limiting the number of requests (rate limit) from one IPv4 subnet. If there are no restrictions, then service returns a Chuck Norris quote.

It is allowed to make 100 requests from one subnet per 2 minutes. Subnet: `/24 (mask 255.255.255.0)`

## Schema
<p align="center" width="100%">
    <img src="https://i.ibb.co/Vv8PGmn/Untitled-Diagram-drawio.png"> 
</p>

## Installation

The application is packaged in [docker](https://www.docker.com/) containers. You must also have docker-compose installed in order to run the application. Command to run the application:

```bash
make
```

## Used

- Golang 1.18
- Redis
- Docker

## License
[MIT](https://choosealicense.com/licenses/mit/)