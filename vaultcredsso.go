package main

import (
“context”
“fmt”
“log”
“os”

```
"github.com/aws/aws-sdk-go-v2/config"
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

// Login automático com SSO
err = loginVaultAutoSSO(client)
if err != nil {
	log.Fatalf("Erro no login: %v", err)
}

fmt.Println("🎉 Login realizado com sucesso!")
```

}

// Função principal - extrai credenciais SSO e faz login automático
func loginVaultAutoSSO(client *api.Client) error {
fmt.Println(“🔄 Extraindo credenciais SSO automaticamente…”)

```
// 1. Carregar credenciais do SSO
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao carregar config AWS: %w", err)
}

// 2. Obter credenciais do provider (extrai do SSO automaticamente)
creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao obter credenciais SSO: %w", err)
}

fmt.Println("✅ Credenciais SSO extraídas com sucesso")

// 3. Definir as credenciais como variáveis de ambiente temporárias
// (para que a biblioteca do Vault possa usá-las)
originalEnv := preserveOriginalEnv()
defer restoreOriginalEnv(originalEnv)

os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)

fmt.Println("🔧 Credenciais definidas temporariamente")

// 4. Usar a biblioteca oficial do Vault
vaultRole := os.Getenv("VAULT_ROLE")
if vaultRole == "" {
	return fmt.Errorf("defina a variável VAULT_ROLE com o nome da sua role no Vault")
}

awsAuth, err := aws.NewAWSAuth(
	aws.WithRegion("us-east-1"), // ajuste para sua região
	aws.WithIAMAuth(),
	aws.WithRole(vaultRole),
)
if err != nil {
	return fmt.Errorf("erro ao criar AWS auth: %w", err)
}

// 5. Fazer login no Vault
fmt.Printf("🚀 Fazendo login no Vault com role: %s\n", vaultRole)

authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
if err != nil {
	return fmt.Errorf("erro no login Vault: %w", err)
}

fmt.Printf("🎉 Login bem-sucedido!\n")
fmt.Printf("🔑 Token: %s\n", authInfo.Auth.ClientToken)
fmt.Printf("⏰ TTL: %d segundos\n", authInfo.Auth.LeaseDuration)
fmt.Printf("📋 Policies: %v\n", authInfo.Auth.Policies)

return nil
```

}

// Preservar variáveis de ambiente originais
func preserveOriginalEnv() map[string]string {
return map[string]string{
“AWS_ACCESS_KEY_ID”:     os.Getenv(“AWS_ACCESS_KEY_ID”),
“AWS_SECRET_ACCESS_KEY”: os.Getenv(“AWS_SECRET_ACCESS_KEY”),
“AWS_SESSION_TOKEN”:     os.Getenv(“AWS_SESSION_TOKEN”),
}
}

// Restaurar variáveis de ambiente originais
func restoreOriginalEnv(original map[string]string) {
for key, value := range original {
if value == “” {
os.Unsetenv(key)
} else {
os.Setenv(key, value)
}
}
}

// Função auxiliar para debug - mostrar credenciais atuais
func debugCredentials() {
fmt.Println(“🔍 Debug: Verificando credenciais…”)

```
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	fmt.Printf("❌ Erro ao carregar config: %v\n", err)
	return
}

creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	fmt.Printf("❌ Erro ao obter credenciais: %v\n", err)
	return
}

fmt.Printf("✅ Access Key: %s\n", creds.AccessKeyID[:10]+"...")
fmt.Printf("✅ Secret Key: %s\n", creds.SecretAccessKey[:10]+"...")
fmt.Printf("✅ Session Token: %s\n", creds.SessionToken[:20]+"...")
fmt.Printf("✅ Source: %s\n", creds.Source)
```

}