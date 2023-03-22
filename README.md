# MSISDNApp

Starts with ```docker compose up```

Takes a full MSISDN number as input and returns:
- MNO Identifier
- Subscriber Number
- Country Identifier according to ISO 3166-1-alpha-2
- Country Code

Features authentication via JWT tokens. Registration is available with a native form or Oauth2 Social Login via Google/Github.  
Can also authenticate via POST calls to ```/api/register``` or ```/api/login``` with a JSON body.

Responds with JSON due to its wide compatibility and readibility by many languages and APIs.  
Can be called either via a POST call to its endpoint ```/service/api/lookup``` or the html page, both restricted to users.

Has a administrator page for viewing and managing the database.
