# backend-balance

This service sends the customer's balance and a summary of transactions, which are loaded from a CSV file, via email.

## Solution Explanation

To solve this challenge, AWS services were utilized, working with AWS Lambda, S3 buckets, and Amazon SES. The Lambda logic was implemented using the Golang programming language. It is triggered when a file is inserted or updated in the S3 bucket. The logic inside the Lambda reads the metadata of the S3 bucket, downloads the file, processes it, calculates the balance and the summary, and then sends an email with a summary, as per the requirements of the test.

### Features and Restrictions

* My personal AWS account was used to create the project's architecture.
* If the file size is greater than 1MB, an error is thrown, and the file is not processed. This is done to avoid the high processing of the lambda.
* The format of the .csv is the same of the statement of the test.
* The subject of the email is 'Transaction Summary'
* The company logo is located in this [S3 bucket](https://stori-resources.s3.amazonaws.com/stori_logo.png)

### Code Structure

The code is organized into four main layers:

* Controller: Here, the logic for reading the input, which is the file with the transactions, is implemented.
* Services: In this layer, the business logic for creating the summary is developed.
* Model: This layer contains the necessary DTOs (Data Transfer Objects).
* Notifications: The logic for sending emails using Amazon SES is placed in this layer.

`dependency_container` File:
This file contains the code for injecting each required instance between the different layers. Reviewing this file helps in understanding the flow.

`environment_variables` File:
This file provides a structure with all the attributes that should be populated with values from environment variables. For practical purposes, default values are also included.

## Installation

To compile the project, follow these instructions:

### AWS

* Create an S3 bucket.
* Create a Lambda function and configure it to be triggered when an object is inserted or updated in the S3 bucket.
* Configure the source email in AWS SES.
* In IAM, add the necessary policies to the Lambda role to read and download from the S3 bucket created in the first step.
* Also, add the policy to grant access to SES resources, enabling the use of the source email with the Lambda.

### Lambda Project

* Compile the code: `GOARCH=amd64 GOOS=linux go build -o main`
* Zip the compiled code: `zip main.zip main`
* Define environment variables. These values are set as environment variables:

- `SES_REGION`
- `DESTINATION_EMAIL`
- `SOURCE_EMAIL`

* Add the logic to the lambda, uploading the file `main.zip`.


## Missing or incomplete

* Some unit tests were added to demonstrate my knowledge of unit test coverage.
* Saving the information in a database. I apologize for not including this feature; however, I had a plan on how to implement it. I intended to create a MySQL database with two tables. The first table would be used to insert the number and type of documents, and the second one to insert the transactions. The document type and number would be extracted from the file name, which should have the format 'CC10432345.csv'. There is an example in this project.
* Refactoring of logic to get the HTML template from an S3 Bucket
