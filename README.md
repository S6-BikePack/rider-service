<div align="center">
  <a href="https://github.com/S6-BikePack">
    <img src="assets/logo.png" alt="logo" width="200" height="auto" />
  </a>
  <h1>BikePack - Rider-Service</h1>

  <p>
    Part of the S6 BikePack project.
  </p>


<!-- Badges -->
<p>

</p>

<h4>
    <a href="https://github.com/S6-BikePack">Home</a>
  <span> · </span>
    <a href="https://github.com/S6-BikePack/rider-service#-about-the-project">Documentation</a>
  </h4>
</div>

<br />

<!-- Table of Contents -->
# 📓 Table of Contents

- [About the Project](#-about-the-project)
    * [Architecture](#-architecture)
    * [Tech Stack](#%EF%B8%8F-tech-stack)
    * [Environment Variables](#-environment-variables)
    * [Messages](#-messages)
    * [Data](#-data)
- [Getting Started](%EF%B8%8F-getting-started)
    * [Prerequisites](%EF%B8%8F-prerequisites)
    * [Running Tests](#-running-tests)
    * [Run Locally](#-run-locally)
    * [Deployment](#-deployment)
- [Usage](#-usage)



<!-- About the Project -->
## ⭐ About the Project

The Rider-Service is the service for the BikePack project that handles all riders in the system. 
A rider it keeps track of which riders are online in the different service-areas and their last known location.
Using the system riders can register to the system.

<!-- Architecture -->
### 🏠 Architecture
For this service I have chosen a Hexagonal architecture. This keeps the service loosely coupled and thus flexible when having to change parts of the system.

<!-- TechStack -->
### 🛰️ Tech Stack
#### Language
  <ul>
    <li><a href="https://go.dev/">GoLang</a></li>
</ul>

#### Dependencies
  <ul>
    <li><a href="https://github.com/gin-gonic/gin">Gin</a><span> - Web framework</span></li>
    <li><a href="https://github.com/gin-gonic/gin">Amqp091-go</a><span> - Go AMQP 0.9.1 client</span></li>
    <li><a href="https://github.com/swaggo/swag">Swag</a><span> - Swagger documentation</span></li>
    <li><a href="https://gorm.io/index.html">GORM</a><span> - ORM library</span></li>
  </ul>

#### Database
  <ul>
    <li><a href="https://www.postgresql.org/">PostgreSQL</a></li>
</ul>

<!-- Env Variables -->
### 🔑 Environment Variables

This service has the following environment variables that can be set:

`PORT` - Port the service runs on

`RABBITMQ` - RabbitMQ connection string

`DATABASE` - Database connection string

<!-- Messages -->
## 📨 Messages

### Publishing
The service publishes the following messages to the RabbitMQ server:

---
**rider.create**

Published when a new rider is created in the system.
Sends the newly created rider in the  body.

```json
{
  "userid": "string", 
  "status": "int",
  "serviceAreaId": "int",
  "capacity": {
    "width": "int",
    "height": "int",
    "depth": "int"
  }
}
```



---
**rider.update**

Published when a delivery is updated in the system.
Sends the updated delivery in the  body.

```json
{
  "userid": "string", 
  "status": "int",
  "serviceAreaId": "int",
  "capacity": {
    "width": "int",
    "height": "int",
    "depth": "int"
  }
}
```



---
**rider.{service-area}.update.location**

Published when a rider updates their location.
Is pushed to a service-area specific topic.

```json
{
  "id": "string", 
  "location": {
    "latitude": "float",
    "longitude": "float
  }
}
```

<!-- Data -->

##  🗃️ Data

This service stores the following data:

```json
{
  "id": "string", 
  "user": {
    "id": "string",
    "name": "string",
    "lastname": "string",
  }
  "status": "int",
  "serviceArea": {
    "id": "int",
    "identifier": "string"
  },
  "capacity": {
    "width": "int",
    "height": "int",
    "depth": "int"
  },
  "location": {
    "latitide": "float",
    "longitude": "float"
  }
}
```

<!-- Getting Started -->
## 	🛠️ Getting Started

<!-- Prerequisites -->
### ‼️ Prerequisites

Building the project requires Go 1.18.

This project requires a PostgreSQL compatible database with a database named `rider` and a RabbitMQ server.
The easiest way to setup the project is to use the Docker-Compose file from the infrastructure repository.

<!-- Running Tests -->
### 🧪 Running Tests

-

<!-- Run Locally -->
### 🏃 Run Locally

Clone the project

```bash
  git clone https://github.com/S6-BikePack/rider-service
```

Go to the project directory

```bash
  cd rider-service
```

Run the project (Rest)

```bash
  go run cmd/rest/main.go
```


<!-- Deployment -->
### 🚀 Deployment

To build this project run (Rest)

```bash
  go build cmd/rest/main.go
```


<!-- Usage -->
## 👀 Usage

### REST
Once the service is running you can find its swagger documentation with all the endpoints at `/swagger`