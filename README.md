# MSISDNApp



Takes a full MSISDN number as input and returns:
- MNO Identifier
- Subscriber Number
- Country Identifier according to ISO 3166-1-alpha-2
- Country Code

Features authentication via JWT tokens. Registration is available with a native form or Oauth2 Social Login via Google/Github.  
Can also authenticate via POST calls to ```/api/register``` or ```/api/login``` with a JSON body.

Responds with JSON due to its wide compatibility and readibility by many languages and APIs.  
Can be called either via a POST call to its endpoint ```/service/api/lookup``` or the html page, both restricted to users.

Uses a Hashicorp Vault for storing and fetching the application secrets.

Has a administrator page for viewing and managing the database.

The app functionality does not account for mobile number portability, and uses a small initialized test set of values in the database as a proof of concept.

# ⚙️Usage

Populate the .env files in ```config/``` with your parameters.  
Set the enviroment variable ```MY_VAULT_TOKEN``` in your terminal with a value of the root token to be created.  
&nbsp;&nbsp;Linux: ```export MY_VAULT_TOKEN=token123```  
&nbsp;&nbsp;Windows: CMD:```set MY_VAULT_TOKEN=token123``` or Powershell: ```$env:MY_VAULT_TOKEN='token123'```  
Start the app with ```docker compose up```  


Normally the vault and the secrets would be manually set and managed. For the sake of seamless startup with docker that's done with a docker container that automatically unseals the vault and sets the app secrets as well as a root permission token with enviromental variables and files.
