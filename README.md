# Microsservices
<p align="left">	
  <img src="https://img.shields.io/github/last-commit/bentodvictor/microsservices?color=c1e7e3">
  <img src="https://img.shields.io/github/issues/bentodvictor/microsservices?color=c1e7e3">
  <img src="https://img.shields.io/github/issues-pr/bentodvictor/microsservices?color=c1e7e3">
  <img src="https://img.shields.io/github/languages/count/bentodvictor/microsservices?color=c1e7e3">
  <img src="https://img.shields.io/github/downloads/bentodvictor/microsservices/total?color=c1e7e3">
  <img src="https://img.shields.io/github/repo-size/bentodvictor/microsservices?color=c1e7e3">
  <img src="https://img.shields.io/badge/license-MIT-c1e7e3">
  <img alt="Stargazers" src="https://img.shields.io/github/stars/bentodvictor/microsservices?color=c1e7e3&logo=github">
</p>

Project developed at the event **Avança DEV** taught by [FullCycle](https://fullcycle.com.br/).

The objective of the course is to learn about microservices, which will be developed in GO, in addition to increasing the range of known programming languages.

---

## Microsservices Details
The application is using RabbitMQ (Docker image).

- *Checkup*: 
  - Load the base template,
  - Gets the data informed: coupon and credit card number,
  - Sends information to the **orders (orders_ex)** queue;
- *Payment*:
  - Reading and processing the information in the queue,
  - Invokes microservice to validate coupon;
- *Validate Coupon*:
  - Coupon validation,
  - Invokes microservice to validate credcard number;
- *Validate Credcard*:
  - Validate credcard number;
  
*If any microservice breaks, orders are sent to the **orders_dql (dlx)** queue and every 3 seconds an attempt is made to resubmit to the **orders (orders_ex)** queue.*

## Execute
All microservices are Docker images and a complete docker-compose image was created to execute and fully integrate the application.

To start the application, execute:

```docker-compose up -d```

To end the application, execute:

```docker-compose down```

If you wanna remove the images of all microsservice, run `docker image rm $image_id`.
