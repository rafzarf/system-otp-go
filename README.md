# Implementation of OTP System using Go and React.js

Here is the Implementation of OTP system using Go as a Backend and React.js as Frontend, the project made for one KSJ Course Module.

### Brief Project Description

This project aims to create a simple OTP (One-Time Password) system for demonstrating network security using Go for the backend and React.js for the frontend.

### Faculty Lead(s)
Dr. Susetyo Bagas Bhaskoro, S.T., M.T.

## How to Run

### Prerequisites

Before you start this project, these are the prerequisites required to be installed:

-   Have the latest version of Golang installed on your system. If you don't have it installed, you can download it from [https://go.dev/doc/install](https://go.dev/doc/install).
-   Have a basic understanding of API designs and CRUD (Create, Read, Update, Delete) patterns.
-   Familiarity with Gin Gonic (a web framework for Golang) and GORM (a Golang Object-Relational Mapping library) will be beneficial.

## Running the Golang 2FA App Locally

To run the Golang-based Two-Factor Authentication (2FA) application locally, follow these steps:

- Clone this repository into your local pc.
- Open your terminal and navigate to the project directory.
- Run the following command to install all the necessary Go packages:
```
go mod tidy
```
- Start the Gin Gonic HTTP Server by Running:
```
go run main.go
```
- You should now have the Golang 2FA app running locally. You can test the 2FA verification system using any API testing software of your choice.

## Running the 2FA Frontend Built with React.js

To run the frontend of the Two-Factor Authentication (2FA) system, which is built with React.js, follow these steps:

1.  Clone or download the source code from this repository

2.  Open your terminal and navigate to the project directory.

3.  Run the following command to install all the required dependencies using `npm` or `yarn`:
```
npm install
```
4. Start the Vite deelopment server by running:
```
npm dev
```
5. Open your web browser and navigate to http://localhost:3000 to test the two-factor authentication system against the API


## API 
### Backend API Specification

We have implemented the backend API for the OTP system following the OpenAPI 3.0.3 specification. Below are the details of the API endpoints:

-   **Server URL:** [http://localhost:8080/v1](http://localhost:8080/v1) (Development server)

#### User Registration

-   **Endpoint:** `/auth/register`
-   **Description:** Register a new user
-   **Request Body:** User registration data
-   **Response:** Successful registration

#### User Login

-   **Endpoint:** `/auth/login`
-   **Description:** User login
-   **Request Body:** User login credentials
-   **Response:** Successful login

#### Generate OTP Token

-   **Endpoint:** `/auth/otp/generate`
-   **Description:** Generate OTP token
-   **Request Body:** User login credentials
-   **Response:** OTP token generated successfully

#### Verify OTP Token

-   **Endpoint:** `/auth/otp/verify`
-   **Description:** Verify OTP token
-   **Request Body:** OTP verification data
-   **Response:** OTP token verified successfully

#### Validate OTP Token

-   **Endpoint:** `/auth/otp/validate`
-   **Description:** Validate OTP token
-   **Request Body:** OTP validation data
-   **Response:** OTP token validated successfully

#### Disable OTP for User

-   **Endpoint:** `/auth/otp/disable`
-   **Description:** Disable OTP for user
-   **Request Body:** OTP disable data
-   **Response:** OTP disabled successfully

#### Send OTP via Telegram

-   **Endpoint:** `/auth/send-otp-telegram`
-   **Description:** Send OTP via Telegram
-   **Request Body:** Send OTP to Telegram
-   **Response:** OTP sent via Telegram successfully
