# AWS Serverless Microservices with GO challenge

Follow this steps to deploy and use

## 1) Create an IAM User
Enter in AWS Web

Create a new user in IAM.
Assign the following permissions:
- AdministratorAccess
- AmazonAPIGatewayInvokeFullAccess
- AmazonS3FullAccess
- AWSLambdaFullAccess

These steps will ensure that the IAM user has the necessary permissions.



## 2) Follow the guide for AWS SAM Deployments

- To run in local or deploy Follow this [guide](./src/microservices/README.md) from AWS SAM proyects

## 3) Use the API

If the deployment is successful, you will see in the logs that an API Gateway endpoint has been provisioned with a structure similar to this:

Key: `microservicesAPI`
Description: API Gateway endpoint URL for Prod environment
Value: `https://{awslink}.{region}.amazonaws.com/Prod/`

The available services are:

**/orders**

Here you will make a POST request with the header `Content-Type: application/json` and the following JSON structure:

```json
{
  "user_id": "string",
  "item": "string",
  "quantity": number,
  "total_price": number
}
```

**/payments**

Here you will use the `order_id` generated from the POST request to `/orders`, and make a new POST request to process the payment with the following structure:

```json
{
  "order_id": "string",
  "status": "string"
}
```

These endpoints allow you to interact with the microservices via API Gateway in your production environment.

exampleURI: `https://{awslink}.{region}.amazonaws.com/Prod/orders`

## 4) 
- Enter in AWS web console and explore DynamoDB , SQS and Lambda to check the interactions
