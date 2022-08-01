# youtube-search-api

# Project Goal

To make an API to fetch latest videos sorted in reverse chronological order of their publishing date-time from YouTube for a given tag/search query in a paginated response.

<details>
  <summary>Click to expand!</summary>

## Requirements:

- [x] Server should call the YouTube API continuously in background (async) with some interval (say 10 seconds) for fetching the latest videos for a predefined search query and should store the data of videos (specifically these fields - Video title, description, publishing datetime, thumbnails URLs and any other fields you require) in a database with proper indexes.
- [x] A GET API which returns the stored video data in a paginated response sorted in descending order of published datetime.
- [x] A basic search API to search the stored videos using their title and description.
- [x] Add support for supplying multiple API keys so that if quota is exhausted on one, it automatically uses the next available key.
- [x] Optimise search api, so that it's able to search videos containing partial match for the search query in either video title or description.
    - Ex 1: A video with title *`How to make tea?`* should match for the search query `tea how`
    
### Reference:

- YouTube data v3 API: [https://developers.google.com/youtube/v3/getting-started](https://developers.google.com/youtube/v3/getting-started)
- Search API reference: [https://developers.google.com/youtube/v3/docs/search/list](https://developers.google.com/youtube/v3/docs/search/list)

</details>

# Tech Stack
- GO 1.18
- Postgres
- RabbitMQ

## How it works?

- When a client request data using REST API, Server will fetch for the latest videos based on the video published date and send back a paginated response.
- In the backgroud, Cron job will run continuously at scheduled interval and fetch the latest videos from YouTube and send that videos to AMQP.
- AMQP consumes the data and insert videos to Postgres Database.

## Getting Started

### Using Docker(Recommended):

1. Clone the repository using git clone:
```
$ git clone https://github.com/Gohelraj/youtube-search-api
$ cd youtube-search-api
```
2. Copy the `.env.example` file to new `.env` file:
```
$ cp .env.example .env
```
3. Update the `.env` file, Add one or multiple(comma separated) [YouTube data v3 API Keys](https://developers.google.com/youtube/v3/getting-started) in `GOOGLE_API_KEYS` variable.
4. Spin up the docker container:
```
$ docker-compose up
```
If permission error occurs, run command as root:
```
$ sudo docker-compose up
```
- The server will start listening on port `8087`
- Incase you have problems running due to ports or stuff already in use, try running the script below commands:
```
$ chmod +x ./docker_reset.sh 
$ sudo ./docker_reset.sh`
```
Note: Be careful while using it as it will kill and remove all other containers as well and thus might lead to loss of your work.

### Using Source Code:

#### Prerequisites you need to set up on your local computer:
1. [Golang](https://go.dev/doc/install)
2. [Postgres](https://www.postgresql.org/download/linux/ubuntu/)
3. [RabbitMQ](https://www.rabbitmq.com/download.html)
4. [Dbmate](https://github.com/amacneil/dbmate#installation)

#### Getting Started:

1. Clone the repository using git clone:
```
1) git clone https://github.com/Gohelraj/youtube-search-api
2) cd youtube-search-api
```
2. Copy the `.env.example` file to new `.env` file:
```
cp .env.example .env
```
3. Update the `.env` file and update below configurations:
   1. Add one or multiple(comma separated) [YouTube data v3 API Keys](https://developers.google.com/youtube/v3/getting-started) in `GOOGLE_API_KEYS` variable.
   2. Update AMQP and Postgres credentials with your local configurations.
   3. Add Postgres database URL in `DATABASE_URL` variable.
4. Run `dbmate migrate` to migrate database schema.
5. Run `go mod vendor` to install all the dependencies.
6. Run `go run cmd/main.go` to run the programme.

## API Endpoints

### 1. Get Videos With Pagination
`GET /videos` - Returns the latest videos sorted by descending order of published datetime in a paginated response.
#### Request Query Parameters:
| Param | Type | Default | Description| Sample |
| --- | --- | --- | --- | --- |
| limit | integer, optional | 50 | Number of records to return, Must be =< 200 | limit=100 |
| offset | integer, optional | 0 | Used to identify the starting point to return rows from | offset=100 |

### 2. Search Videos By Keyword (In Title/Description)
`POST /videos/search` - Returns the videos matching with the search keyword.
#### Request Body Parameters:
| Param | Type | Default | Description| Sample |
| --- | --- | --- | --- | --- |
| searchString | string, required |  | Search string to match in video's title and description  | {"searchString":"how to make tea"} |

_The exact API usage can be inspected via the [`api.postman_collection.json`](./api.postman_collection.json) postman collection._
