package commands

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

// Connect sets up a connection to Consul and hands back a client.
func Connect() (*consul.Client, error) {
	config := consul.DefaultConfig()
	config.Token = Token
	config.Address = ConsulServer
	consul, err := consul.NewClient(config)
	if err != nil {
		Log("Consul connection is bad.", "info")
		return nil, err
	}
	return consul, nil
}

// Get grabs a key's value from Consul.
func Get(c *consul.Client, key string) string {
	var value string
	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
	}
	return value
}

// Set sets a key's value inside Consul.
func Set(c *consul.Client, key, value string) bool {
	p := &consul.KVPair{Key: key, Value: []byte(value)}
	kv := c.KV()
	Log(fmt.Sprintf("key='%s' value='%s'", key, value), "info")
	_, err := kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
	return true
}
