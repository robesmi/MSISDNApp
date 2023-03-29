package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type Vault struct {
	Vault api.Client
}
//go:generate mockgen -destination=../mocks/vault/mockVault.go -package=vault github.com/robesmi/MSISDNApp/vault VaultInterface

type VaultInterface interface{
	Insert(string, map[string]interface{}) error
	Fetch(string, ...string) (map[string]string, error)
}

// New takes a hashicorp vault config and a registered vault token as input
// and returns a client to interact with the vault
func New(conf *api.Config, token string) (VaultInterface, error){
	newClient, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	newVault := Vault{*newClient}
	newVault.Vault.SetToken(token)
	return &newVault, nil
}


// Insert inserts a new key/value pair at the provided path with the provided values
func (v Vault) Insert(path string, kv map[string]interface{}) (error){

	ctx := context.Background()
	
	_, err := v.Vault.KVv2("secret").Put(ctx, path, kv)
	if err != nil {
		return err
	}
	return nil
}

// Fetch takes the path and an arbitrary amount of keys that should
// be present in the vault at that path and returns a map[string]string
// with the ones that match. If no keys are provided, returns all values
func (v Vault) Fetch(path string, key ...string) (map[string]string, error){

	ctx := context.Background()

	res, err := v.Vault.KVv2("secret").Get(ctx, path)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)

	if len(key) == 0{
		for k,v := range res.Data{
			strKey := fmt.Sprint(k)
			strVal := fmt.Sprint(v)
			result[strKey] = strVal
		}
		return result, nil
	}else{
		for _, v := range key{
			if val, ok := res.Data[v].(string); ok {
				result[v] = val
			}
		}
		return result, nil
	}
}