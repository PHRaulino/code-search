package main

import (
“context”
“fmt”
“log”
“net/http”
“net/url”
“os”
“strings”

```
"github.com/aws/aws-sdk-go-v2/config"
"github.com/aws/aws-sdk-go-v2/service/sts"
"github.com/hashicorp/vault/api"
```

)

func main() {
// Configurar cliente do Vault
vaultConfig := api.DefaultConfig()
vaultConfig.Address = os.Getenv(“VAULT_ADDR”)

```
client, err := api.NewClient(vaultConfig)
if err != nil {
	log.Fatalf("Erro ao criar cliente Vault: %v", err)
}

// Método 1: Login automático extraindo credenciais do SSO
err = loginVaultWithAutoCredentials(client)
if err != nil {
	log.Fatalf("Erro no login automático: %v", err)
}

// Testar o token
secret, err := client.Logical().Read("secret/data/test")
if err != nil {
	log.Printf("Erro ao ler secret (normal se não existir): %v", err)
} else {
	fmt.Printf("✅ Token válido! Secret: %+v\n", secret)
}
```

}

// Extrair credenciais automaticamente do SSO e fazer login manual
func loginVaultWithAutoCredentials(client *api.Client) error {
fmt.Println(“🔄 Extraindo credenciais AWS automaticamente…”)

```
// Carregar credenciais do SSO
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao carregar config AWS: %w", err)
}

// Obter credenciais do provider
creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao obter credenciais: %w", err)
}

fmt.Printf("✅ Credenciais obtidas automaticamente\n")
fmt.Printf("🔑 Access Key: %s\n", creds.AccessKeyID[:10]+"...")
fmt.Printf("🔑 Session Token: %s\n", creds.SessionToken[:20]+"...")

// Obter ARN atual
stsClient := sts.NewFromConfig(cfg)
identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
if err != nil {
	return fmt.Errorf("erro ao obter identity: %w", err)
}
fmt.Printf("👤 Identity: %s\n", *identity.Arn)

// Criar requisição STS assinada manualmente
loginData, err := createVaultLoginPayload(cfg, creds)
if err != nil {
	return fmt.Errorf("erro ao criar payload de login: %w", err)
}

// Fazer login no Vault
vaultRole := os.Getenv("VAULT_ROLE")
if vaultRole == "" {
	vaultRole = "SUA-ROLE-DO-VAULT" // substitua pelo nome da sua role
}

loginData["role"] = vaultRole

fmt.Printf("🚀 Fazendo login no Vault com role: %s\n", vaultRole)

secret, err := client.Logical().Write("auth/aws/login", loginData)
if err != nil {
	return fmt.Errorf("erro no login Vault: %w", err)
}

if secret.Auth == nil {
	return fmt.Errorf("resposta de auth vazia")
}

// Definir token no cliente
client.SetToken(secret.Auth.ClientToken)

fmt.Printf("🎉 Login bem-sucedido!\n")
fmt.Printf("🔑 Token: %s\n", secret.Auth.ClientToken)
fmt.Printf("⏰ TTL: %d segundos\n", secret.Auth.LeaseDuration)

return nil
```

}

// Criar payload de login para o Vault usando credenciais automáticas
func createVaultLoginPayload(cfg config.Config, creds aws.Credentials) (map[string]interface{}, error) {
// Criar cliente STS
stsClient := sts.NewFromConfig(cfg)

```
// Preparar requisição GetCallerIdentity
input := &sts.GetCallerIdentityInput{}
req, err := stsClient.GetCallerIdentityRequest(context.TODO(), input)
if err != nil {
	return nil, fmt.Errorf("erro ao criar request STS: %w", err)
}

// Assinar requisição
err = req.Sign()
if err != nil {
	return nil, fmt.Errorf("erro ao assinar request: %w", err)
}

// Extrair dados da requisição assinada
httpReq := req.HTTPRequest

// Ler body da requisição
body := "Action=GetCallerIdentity&Version=2011-06-15"

// Codificar headers
headers := encodeHeaders(httpReq.Header)

// Criar payload para o Vault
payload := map[string]interface{}{
	"iam_http_request_method":  httpReq.Method,
	"iam_request_url":         httpReq.URL.String(),
	"iam_request_body":        body,
	"iam_request_headers":     headers,
}

fmt.Printf("📋 Payload criado para o Vault:\n")
fmt.Printf("   Method: %s\n", payload["iam_http_request_method"])
fmt.Printf("   URL: %s\n", payload["iam_request_url"])
fmt.Printf("   Headers: %d caracteres\n", len(headers))

return payload, nil
```

}

// Codificar headers para o formato esperado pelo Vault
func encodeHeaders(headers http.Header) string {
var parts []string
for name, values := range headers {
for _, value := range values {
parts = append(parts, fmt.Sprintf(”%s:%s”, name, value))
}
}
return strings.Join(parts, “,”)
}

// Alternativa: Usar as credenciais como variáveis de ambiente temporárias
func loginWithTemporaryEnvVars(client *api.Client) error {
fmt.Println(“🔄 Método alternativo: Definindo env vars temporárias…”)

```
// Carregar credenciais
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao carregar config: %w", err)
}

creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	return fmt.Errorf("erro ao obter credenciais: %w", err)
}

// Salvar valores atuais (para restaurar depois)
oldAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
oldSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
oldSessionToken := os.Getenv("AWS_SESSION_TOKEN")

// Definir credenciais como env vars temporárias
os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)

defer func() {
	// Restaurar valores originais
	os.Setenv("AWS_ACCESS_KEY_ID", oldAccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", oldSecretKey)
	os.Setenv("AWS_SESSION_TOKEN", oldSessionToken)
}()

fmt.Println("✅ Credenciais definidas temporariamente como env vars")

// Agora usar a biblioteca oficial do Vault (que vai pegar as env vars)
awsAuth, err := aws.NewAWSAuth(
	aws.WithRegion("us-east-1"),
	aws.WithIAMAuth(),
	aws.WithRole(os.Getenv("VAULT_ROLE")),
)
if err != nil {
	return fmt.Errorf("erro ao criar AWS auth: %w", err)
}

authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
if err != nil {
	return fmt.Errorf("erro no login: %w", err)
}

fmt.Printf("🎉 Login bem-sucedido! Token: %s\n", authInfo.Auth.ClientToken)
return nil
```

}

// Função para extrair e mostrar credenciais (debug)
func showCurrentCredentials() {
fmt.Println(“🔍 Credenciais atuais no ambiente:”)

```
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	fmt.Printf("❌ Erro: %v\n", err)
	return
}

creds, err := cfg.Credentials.Retrieve(context.TODO())
if err != nil {
	fmt.Printf("❌ Erro: %v\n", err)
	return
}

fmt.Printf("🔑 Access Key ID: %s\n", creds.AccessKeyID)
fmt.Printf("🔑 Secret Access Key: %s***\n", creds.SecretAccessKey[:10])
fmt.Printf("🔑 Session Token: %s***\n", creds.SessionToken[:20])
fmt.Printf("📅 Expires: %v\n", creds.Expires)
fmt.Printf("🏷️  Source: %s\n", creds.Source)
```

}