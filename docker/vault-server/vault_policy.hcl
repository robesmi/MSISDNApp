path "kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "kv/my-secret" {
  capabilities = ["read"]
}

path "secret/*" {
  capabilities = ["create". "read", "update", "list"]
}