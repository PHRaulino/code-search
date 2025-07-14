package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/aws"
)

func main() {
	// Configurar cliente do Vault
	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR") // ex: "https://vault.example.com:8200"
	
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Erro ao criar cliente Vault: %v", err)
	}

	// Método 1: Usando a biblioteca auth/aws do Vault (Recomendado)
	err = loginWithVaultAwsAuth(client)
	if err != nil {
		log.Fatalf("Erro no login com Vault AWS auth: %v", err)
	}

	// Método 2: Login manual (alternativo)
	// err = loginManualAws(client)
	// if err != nil {
	//     log.Fatalf("Erro no login manual: %v", err)
	// }

	// Testar o token obtido
	secret, err := client.Logical().Read("secret/data/test")
	if err != nil {
		log.Fatalf("Erro ao ler secret: %v", err)
	}
	
	fmt.Printf("Token válido! Secret lido: %+v\n", secret)
}

// Método 1: Usando a biblioteca oficial do Vault
func loginWithVaultAwsAuth(client *api.Client) error {
	// Configurar AWS SDK
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("erro ao carregar config AWS: %w", err)
	}

	// Criar auth method
	awsAuth, err := aws.NewAWSAuth(
		aws.WithRegion("us-east-1"), // ou sua região
		aws.WithIAMAuth(),
	)
	if err != nil {
		return fmt.Errorf("erro ao criar AWS auth: %w", err)
	}

	// Fazer login
	authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
	if err != nil {
		return fmt.Errorf("erro no login: %w", err)
	}

	fmt.Printf("Login bem-sucedido! Token: %s\n", authInfo.Auth.ClientToken)
	return nil
}

// Método 2: Login manual (para maior controle)
func loginManualAws(client *api.Client) error {
	// Configurar AWS SDK
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("erro ao carregar config AWS: %w", err)
	}

	// Criar cliente STS
	stsClient := sts.NewFromConfig(cfg)

	// Obter identity do caller
	callerIdentity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("erro ao obter caller identity: %w", err)
	}

	// Preparar requisição STS
	stsRequest, _ := stsClient.GetCallerIdentityRequest(context.TODO(), &sts.GetCallerIdentityInput{})
	
	// Assinar a requisição
	err = stsRequest.Sign()
	if err != nil {
		return fmt.Errorf("erro ao assinar requisição STS: %w", err)
	}

	// Extrair headers necessários
	headers := stsRequest.HTTPRequest.Header
	
	// Montar payload para o Vault
	data := map[string]interface{}{
		"iam_http_request_method":  "POST",
		"iam_request_url":         "https://sts.amazonaws.com/",
		"iam_request_body":        "Action=GetCallerIdentity&Version=2011-06-15",
		"iam_request_headers":     encodeHeaders(headers),
		"role":                    "my-vault-role", // Nome do role configurado no Vault
	}

	// Fazer login no Vault
	secret, err := client.Logical().Write("auth/aws/login", data)
	if err != nil {
		return fmt.Errorf("erro no login Vault: %w", err)
	}

	if secret.Auth == nil {
		return fmt.Errorf("resposta de auth vazia")
	}

	// Definir token no cliente
	client.SetToken(secret.Auth.ClientToken)
	
	fmt.Printf("Login manual bem-sucedido! Token: %s\n", secret.Auth.ClientToken)
	fmt.Printf("ARN: %s\n", *callerIdentity.Arn)
	
	return nil
}

// Função auxiliar para codificar headers
func encodeHeaders(headers map[string][]string) string {
	var result string
	for key, values := range headers {
		for _, value := range values {
			if result != "" {
				result += ","
			}
			result += fmt.Sprintf("%s:%s", key, value)
		}
	}
	return result
}

// Exemplo de como configurar variáveis de ambiente
func setupEnvironment() {
	// Definir endereço do Vault
	os.Setenv("VAULT_ADDR", "https://vault.example.com:8200")
	
	// Se usar perfil AWS específico
	os.Setenv("AWS_PROFILE", "my-sso-profile")
	
	// Ou definir credenciais diretamente (se aplicável)
	// os.Setenv("AWS_ACCESS_KEY_ID", "...")
	// os.Setenv("AWS_SECRET_ACCESS_KEY", "...")
	// os.Setenv("AWS_SESSION_TOKEN", "...")
}
