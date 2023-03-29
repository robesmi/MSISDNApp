#!/usr/bin/env sh

set -ex

# Only executed on the very first initalization of a vault and saves the keys in a file on the container
init () {
vault operator init > /vault/file/keys
}

unseal () {
vault operator unseal $(grep 'Key 1:' /vault/file/keys | awk '{print $NF}')
vault operator unseal $(grep 'Key 2:' /vault/file/keys | awk '{print $NF}')
vault operator unseal $(grep 'Key 3:' /vault/file/keys | awk '{print $NF}')
}

log_in () {
   export ROOT_TOKEN=$(grep 'Initial Root Token:' /vault/file/keys | awk '{print $NF}')
   vault login $ROOT_TOKEN
   vault secrets enable -path=secret kv-v2
}

# Creates a token with the id that we give to docker compose up as an enviroment variable(or the default one)
create_token () {
   vault token create -id $MY_VAULT_TOKEN
}

set_app_secrets(){
   vault kv put secret/appvars PORT=$PORT \
   MYSQL_DRIVER=$MYSQL_DRIVER \
   MYSQL_SOURCE=$MYSQL_SOURCE \
   Secret=$Secret \
   AccessTokenPublicKey=$AccessTokenPublicKey \
   AccessTokenPrivateKey=$AccessTokenPrivateKey \
   RefreshTokenPublicKey=$RefreshTokenPublicKey \
   RefreshTokenPrivateKey=$RefreshTokenPrivateKey \
   GoogleClientId=$GoogleClientId \
   GoogleClientSecret=$GoogleClientSecret \
   GoogleRedirect=$GoogleRedirect \
   GithubClientId=$GithubClientId \
   GithubClientSecret=$GithubClientSecret \
   GithubRedirect=$GithubRedirect \
   EncryptKey=$EncryptKey

   vault kv put secret/superuser AdminUsername=$AdminUsername \
   AdminPassword=$AdminPassword
}

if [ -s /vault/file/keys ]; then
   unseal
else
   init
   unseal
   log_in
   create_token
   set_app_secrets
fi

vault status > /vault/file/status