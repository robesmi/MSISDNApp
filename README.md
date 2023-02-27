# MSISDNApp

Starts with ```docker compose up```

Takes a full MSISDN number as input and returns:
- MNO Identifier
- Subscriber Number
- Country Identifier according to ISO 3166-1-alpha-2
- Country Code

Responds via JSON due to its wide compatability and readibility by many languages and APIs.  
Can be called either via a POST call to its endpoint ```/lookup``` or via the html page.
