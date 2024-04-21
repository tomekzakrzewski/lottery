# Lotto Lottery Game

Number based lottery game. User selects 6 unique numbers between 1-99 and receives a ticket ID. Random generated numbers are drawn every Saturday at 12:00. Users can check the lottery results using their unique ticket ID.

## Technologies used

- GO
- Chi
- MongoDB
- Redis
- Docker
- gRPC

## Features
- Microservices architecture
- Containerized using Docker
- Caching with Redis
- gRPC and HTTP communication between services
- Logrus for logging
- Implemented clients for communication

## How to run
Requirements:
- Docker

docker-compose up

## Endpoints
/inputTicket POST
/result/{hash} GET


