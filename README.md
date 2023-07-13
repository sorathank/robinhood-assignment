

# Robinhood Assignment

  

Robinhood Assignment is a simple backend system for storing and fetching interview data. The system handles descriptions, creation times, usernames, a list of comments (which includes the creation time and commenter's name), and status.

  

## Tech Stack

  

The project utilizes:

- Golang

- Gin Framework

- MongoDB

- Redis

- Docker (Dockerfile, Docker-compose)

  

## Prerequisites

  

Ensure you have Docker and Docker-compose installed on your machine.

  

## Installation

  

To run this application:

  

1. Clone the repository

```bash

git  clone  https://github.com/your-username/robinhood-assignment.git

```

  

2. Navigate to the project folder

```bash

cd  robinhood-assignment

```

  

3. Run the Docker-compose command

```bash

docker-compose  up

```

  

This will launch the application and its associated services. The application will be accessible at localhost:8080 by default.

  

Please replace PORT with the correct port number which is mentioned in the .env file or Docker-compose configuration.


# API Endpoints

## Create Interview
**URL**: `/interview`
**Method**: `POST`
**Body**:
```json
{
    "Description": "Test"
}
```
**Response**:
```json
{
    "Create Interview": "Success"
}
```

## Get Interview by ID

**URL**: `/interview/id/:interviewId`
 **Example**: `localhost:8080/interview/id/64afcc9dafab7c977749add6` 
 **Method**: `GET` 
 **Response**:
 ```json
{
    "comments": [
        //...Comments Data...
    ],
    "interview": {
        "Id": "64afcc9dafab7c977749add6",
        "Description": "Test2",
        "User": "creator1",
        "Status": "Todo",
        "CreatedTime": "2023-07-13T10:06:21.6Z"
    }
}

```

## Get Interviews by Page

**URL**: `/interview/page/:page` 
**Example**: `localhost:8080/interview/page/1` 
**Method**: `GET` 
**Response**:
```json
[
    {
        "Id": "64afcc9dafab7c977749add6",
        "Description": "Test2",
        "User": "creator1",
        "Status": "Todo",
        "CreatedTime": "2023-07-13T10:06:21.6Z"
    }
]

```

## Update Interview Status

**URL**: `/interview/status` 
**Method**: `PUT` 
**Body**:
```json
{
    "InterviewId": "64afc1feafab7c977749add5",
    "Status": "In Progress"
}

```
**Response**:
```json
{
    "Update Status": "Success"
}
```
OR
```json
{
    "Update Status": "invalid status"
}
```

## Create Comment

**URL**: `/comment` 
**Method**: `POST` 
**Body**:
```json
{
    "InterviewId": "64afcc9dafab7c977749add6",
    "Content": "test1"
}
```
**Response**:
```json
{
    "Create Comment": "Success"
}
```

## Create User
**URL**: `/user`
**Method**: `POST`
**Body**:
```json
{
    "username": "creator3",
    "password": "creator3",
    "email": "creator3@mail.com"
}
```
**Response**:
```json
{
    "Create User": "Success"
}
```
OR
```json
{
    "Create User": "Duplicated Username"
}
```
OR
```json
{
    "Create User": "Duplicated Email"
}
```

## User Login
**URL**: `/login`
**Method**: `POST`
**Body**:
```json
{
    "username": "creator3",
    "password": "creator3"
}
```
**Response**:
```json
{
    "result": "Login Success"
}
```
OR
```json
{
    "error": "Username or Password is incorrect"
}
```