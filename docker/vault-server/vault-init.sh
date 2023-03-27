#!/usr/bin/env sh

# Start vault
vault server -config ./vault/config/config.hcl

# Export values
export VAULT_ADDR='https://0.0.0.0:8200'
export VAULT_SKIP_VERIFY='true'

# Parse unsealed keys
mapfile -t keyArray < <( grep "Unseal Key " < generated_keys.txt  | cut -c15- )

vault operator unseal ${keyArray[0]}
vault operator unseal ${keyArray[1]}
vault operator unseal ${keyArray[2]}

# Get root token
mapfile -t rootToken < <(grep "Initial Root Token: " < generated_keys.txt  | cut -c21- )
echo ${rootToken[0]} > root_token.txt

export VAULT_TOKEN=${rootToken[0]}

# Enable kv
vault secrets enable -path=kv-v2 kv-v2

# Enable userpass and add default user
vault auth enable userpass
vault policy write vault-policy vault-policy.hcl
vault write auth/userpass/users/admin password=${SECRET_PASS} policies=vault-policy

# Add test value to my-secret
vault kv put kv/my-secret my-value=verysecret