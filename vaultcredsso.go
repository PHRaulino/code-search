package main

import (
“context”
“fmt”
“log”
“os”
“runtime”

```
"github.com/aws/aws-sdk-go-v2/config"
"github.com/aws/aws-sdk-go-v2/service/sts"
"github.com/hashicorp/vault/api"
"github.com/hashicorp/vault/api/auth/aws"
```

)

func main() {
// Configurar Vault
vaultConfig := api.DefaultConfig()
vaultConfig.Address = os.Getenv(“VAULT_ADDR”)

```
client, err := api.NewClient(vaultConfig)
if err != nil {
	log.Fatalf("Erro ao criar cliente Vault: %v", err)
}

// Login híbrido seguro
err = loginVaultHybridSecure(client)
if err != nil {
	log.Fatalf("Erro no login: %v", err)
}

fmt.Println("🎉 Login realizado com sucesso!")
```

}

// Solução híbrida: env vars temporárias com limpeza agressiva
func loginVaultHybridSecure(client *api.Client) error {
fmt.Println(“🔒 Login híbrido seguro…”)

```
// 1. Verificar credenciais
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao carregar config: %w", err)
}

stsClient := sts.NewFromConfig(cfg)
identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
if err != nil {
	return fmt.Errorf("erro ao obter identity: %w", err)
}
fmt.Printf("✅ Identity: %s\n", *identity.Arn)

// 2. Extrair credenciais
creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao obter credenciais: %w", err)
}

// 3. Contexto isolado para login
err = doIsolatedVaultLogin(client, creds)
if err != nil {
	return err
}

// 4. Limpeza agressiva da memória
runtime.GC()
runtime.GC() // Forçar garbage collection duplo

fmt.Println("🧹 Limpeza de memória realizada")

return nil
```

}

// Fazer login em contexto isolado
func doIsolatedVaultLogin(client *api.Client, creds aws.Credentials) error {
// Escopo limitado - variáveis locais são mais seguras
func() {
// Backup env vars apenas neste escopo
origAccessKey := os.Getenv(“AWS_ACCESS_KEY_ID”)
origSecretKey := os.Getenv(“AWS_SECRET_ACCESS_KEY”)
origToken := os.Getenv(“AWS_SESSION_TOKEN”)

```
	// Defer para garantir limpeza mesmo se der panic
	defer func() {
		// Limpeza agressiva
		os.Setenv("AWS_ACCESS_KEY_ID", "")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "")
		os.Setenv("AWS_SESSION_TOKEN", "")
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_SESSION_TOKEN")
		
		// Restaurar valores originais
		if origAccessKey != "" {
			os.Setenv("AWS_ACCESS_KEY_ID", origAccessKey)
		}
		if origSecretKey != "" {
			os.Setenv("AWS_SECRET_ACCESS_KEY", origSecretKey)
		}
		if origToken != "" {
			os.Setenv("AWS_SESSION_TOKEN", origToken)
		}
	}()
	
	// Definir credenciais temporariamente
	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	
	// Login rápido
	vaultRole := os.Getenv("VAULT_ROLE")
	if vaultRole == "" {
		panic("VAULT_ROLE não definida")
	}

	awsAuth, err := aws.NewAWSAuth(
		aws.WithRegion("us-east-1"),
		aws.WithIAMAuth(),
		aws.WithRole(vaultRole),
	)
	if err != nil {
		panic(fmt.Sprintf("erro ao criar AWS auth: %v", err))
	}

	authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
	if err != nil {
		panic(fmt.Sprintf("erro no login: %v", err))
	}

	fmt.Printf("🎉 Token obtido: %s\n", authInfo.Auth.ClientToken)
	
	// As env vars são limpas automaticamente no defer acima
}()

return nil
```

}

// Versão ainda mais paranóica (usando processo filho)
func loginVaultExtremeSecure(client *api.Client) error {
// Para casos extremamente sensíveis, você poderia:
// 1. Criar um processo filho
// 2. Passar credenciais via pipe/socket
// 3. O filho faz login e retorna só o token
// 4. Matar o processo filho

```
// Implementação complexa - só se realmente necessário
return fmt.Errorf("implementar apenas se necessário")
```

}