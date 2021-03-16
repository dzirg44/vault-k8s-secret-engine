package servian

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

// TODO: Finish help text
const backendHelp = `
The Vault dynamic service account backend provides on-demand, dynamic 
credentials for a short-lived k8s service account
`

// TODO: default value for a ttl to make the ttl parameter optional
// TODO: set maxTtlInSeconds via plugin configuration, remove hardcoded value

func K8sServiceAccountFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	k := KubernetesService{}
	b := Backend(&k)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	for k, v := range conf.Config {
		b.Logger().Info(fmt.Sprintf("Config: %s -> %s\n", k, v))
	}
	return b, nil
}

// TODO: implement a backend InitializeFunc to ensure we can connect to k8s
func Backend(k KubernetesInterface) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		Paths: []*framework.Path{
			configurePlugin(&b),
			readSecret(&b),
		},
		Secrets: []*framework.Secret{
			createSecret(&b),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				configPath,
			},
		},
		// TODO: Do we need to use `TypeCredential` instead?
		BackendType: logical.TypeLogical,
	}
	b.kubernetesService = k
	return &b
}

type backend struct {
	*framework.Backend
	kubernetesService KubernetesInterface
}
