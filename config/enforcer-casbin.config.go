package config

import (
	"log"

	"github.com/casbin/casbin/v2"
)

var EnforcerCasbin *casbin.Enforcer

func ConfigCasbinEnforcer() {
	var err error

	EnforcerCasbin, err = casbin.NewEnforcer("casbin/rbac_model.conf", "casbin/rbac_policy.csv")
	if err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}
	err = EnforcerCasbin.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}
}
