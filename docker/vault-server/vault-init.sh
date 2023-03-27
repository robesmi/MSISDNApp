#!/usr/bin/env sh

set -ex

# Does NOT work because 'vault server' command BLOCKS and further commands cannot be used without a second terminal
# Automating this to be seamless in a docker enviroment is going to take time and alcohol
init () {
vault server -config ./vault/config/config.hcl
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
}

create_token () {
   vault token create -id $MY_VAULT_TOKEN
}

if [ -s /vault/file/keys ]; then
   unseal
else
   init
   unseal
   log_in
   create_token
fi

vault status > /vault/file/status