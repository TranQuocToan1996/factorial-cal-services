How to start LOCAL:
- rename and fix the .env or the env in the docker compose file
- Storage: Create S3 bucket or use local disk. Or you can implement memory storage layer. Add base case in the /infrastructure/base_case
- docker compose up -d
- Refer the docs/swagger.json for API info. There are 3 APIs:
    - For request the calculate number only.
    POST /api/v1/factorial
    {
        "number": 4
    }

    - GET the result of the number
    GET /api/v1/factorial/{number}

    - GET metadata: Get the key, bucket, checksum, status of the request calculate. May use for client call S3 get the factorial of big numbers result
    GET /api/v1/factorial/metadata/{number}
- Call API request number. After that call GET the result of the number.


How to deploy to AWS ECS:
- AWS configure ID, Key, region. 
- Refer context some info setup.
- Storage: Create S3 bucket or use local disk. Or you can implement memory storage layer. Add base case in the /infrastructure/base_case
- Use step 1.1 or 1.2
1.1 Refer to the /infrastructure/terraform if u you Terraform
1.2 refer to /infrastructure/ecs and add to chatbot for generate aws commands to create
- Base on the .env file. Set the env files so the service can get it.
- Call API request number. After that call GET the result of the number.


For the deployed version, for simple no Auth require:
Swagger: https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/swagger/index.html

- POST https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial
{
    "number": 4
}

- GET https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial/4

- GET https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial/metadata/4