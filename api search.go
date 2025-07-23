package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type RequirementsResult struct {
	Repository string
	FilePath   string
	Content    string
	HasXRay    bool
	HasBoto3   bool
}

func main() {
	// Token do GitHub (configure como variável de ambiente)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Nome da organização
	orgName := "sua-org-aqui" // Substitua pelo nome da sua organização
	
	ctx := context.Background()
	
	// Configurar cliente autenticado
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Buscar repositórios da organização que contenham "sigapp" no nome
	repos, err := searchRepositories(ctx, client, orgName)
	if err != nil {
		log.Fatal("Error searching repositories:", err)
	}

	fmt.Printf("Found %d repositories with 'sigapp' in name\n\n", len(repos))

	// Buscar requirements.txt em cada repositório
	var results []RequirementsResult
	for _, repo := range repos {
		repoResults, err := searchRequirementsInRepo(ctx, client, orgName, repo)
		if err != nil {
			fmt.Printf("Error searching in %s: %v\n", repo, err)
			continue
		}
		results = append(results, repoResults...)
	}

	// Filtrar resultados que atendem aos critérios
	filteredResults := filterResults(results)

	// Exibir resultados
	fmt.Printf("=== RESULTS ===\n")
	fmt.Printf("Found %d requirements.txt files matching criteria:\n\n", len(filteredResults))

	for _, result := range filteredResults {
		fmt.Printf("Repository: %s\n", result.Repository)
		fmt.Printf("File: %s\n", result.FilePath)
		fmt.Printf("Has aws-sdk-xray: %t\n", result.HasXRay)
		fmt.Printf("Has boto3: %t\n", result.HasBoto3)
		fmt.Println("Content preview:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Println(result.Content[:min(len(result.Content), 500)])
		if len(result.Content) > 500 {
			fmt.Println("... (truncated)")
		}
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println()
	}
}

func searchRepositories(ctx context.Context, client *github.Client, orgName string) ([]string, error) {
	var allRepos []string
	
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			if strings.Contains(strings.ToLower(*repo.Name), "sigapp") {
				allRepos = append(allRepos, *repo.Name)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func searchRequirementsInRepo(ctx context.Context, client *github.Client, orgName, repoName string) ([]RequirementsResult, error) {
	var results []RequirementsResult

	// Buscar arquivos requirements.txt no repositório
	query := fmt.Sprintf("filename:requirements.txt repo:%s/%s", orgName, repoName)
	
	searchOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	searchResult, _, err := client.Search.Code(ctx, query, searchOpts)
	if err != nil {
		return nil, err
	}

	for _, codeResult := range searchResult.CodeResults {
		// Obter conteúdo do arquivo
		fileContent, _, _, err := client.Repositories.GetContents(
			ctx, 
			orgName, 
			repoName, 
			*codeResult.Path, 
			nil,
		)
		if err != nil {
			fmt.Printf("Error getting content for %s/%s: %v\n", repoName, *codeResult.Path, err)
			continue
		}

		content, err := fileContent.GetContent()
		if err != nil {
			fmt.Printf("Error decoding content for %s/%s: %v\n", repoName, *codeResult.Path, err)
			continue
		}

		// Verificar se contém aws-sdk-xray e se NÃO contém boto3
		hasXRay := strings.Contains(strings.ToLower(content), "aws-sdk-xray")
		hasBoto3 := strings.Contains(strings.ToLower(content), "boto3")

		results = append(results, RequirementsResult{
			Repository: repoName,
			FilePath:   *codeResult.Path,
			Content:    content,
			HasXRay:    hasXRay,
			HasBoto3:   hasBoto3,
		})
	}

	return results, nil
}

func filterResults(results []RequirementsResult) []RequirementsResult {
	var filtered []RequirementsResult
	
	for _, result := range results {
		// Critério: deve ter aws-sdk-xray MAS NÃO deve ter boto3
		if result.HasXRay && !result.HasBoto3 {
			filtered = append(filtered, result)
		}
	}
	
	return filtered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
