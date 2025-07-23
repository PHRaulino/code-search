package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

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

type CacheData struct {
	Repositories []string    `json:"repositories"`
	Results      []RequirementsResult `json:"results"`
	Timestamp    time.Time   `json:"timestamp"`
	TeamSlug     string      `json:"team_slug"`
}

const CACHE_FILE = "github_cache.json"
const CACHE_DURATION = 1 * time.Hour // Cache válido por 1 hora

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	orgName := "sua-org-aqui"     // Substitua pela sua organização
	teamSlug := "seu-team-aqui"   // Substitua pelo slug do seu team
	
	ctx := context.Background()
	
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Tentar carregar do cache primeiro
	cachedData := loadCache(teamSlug)
	
	var repos []string
	var results []RequirementsResult
	
	if cachedData != nil && time.Since(cachedData.Timestamp) < CACHE_DURATION {
		fmt.Println("Using cached data...")
		repos = cachedData.Repositories
		results = cachedData.Results
	} else {
		fmt.Println("Fetching fresh data from GitHub API...")
		
		// Buscar repositórios do team
		var err error
		repos, err = getTeamRepositories(ctx, client, orgName, teamSlug)
		if err != nil {
			log.Fatal("Error getting team repositories:", err)
		}

		fmt.Printf("Found %d repositories for team '%s'\n\n", len(repos), teamSlug)

		// Buscar requirements.txt em cada repositório
		for _, repo := range repos {
			repoResults, err := searchRequirementsInRepo(ctx, client, orgName, repo)
			if err != nil {
				fmt.Printf("Error searching in %s: %v\n", repo, err)
				continue
			}
			results = append(results, repoResults...)
		}

		// Salvar no cache
		saveCache(CacheData{
			Repositories: repos,
			Results:      results,
			Timestamp:    time.Now(),
			TeamSlug:     teamSlug,
		})
	}

	// Filtrar resultados
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

func getTeamRepositories(ctx context.Context, client *github.Client, orgName, teamSlug string) ([]string, error) {
	var allRepos []string
	
	opts := &github.ListOptions{PerPage: 100}

	for {
		repos, resp, err := client.Teams.ListTeamReposBySlug(ctx, orgName, teamSlug, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			// Filtrar apenas repositórios Python
			if isPythonRepository(repo) {
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

func isPythonRepository(repo *github.Repository) bool {
	// Verifica pela linguagem principal
	if repo.Language != nil && strings.ToLower(*repo.Language) == "python" {
		return true
	}
	
	// Verifica pelo nome do repositório (se contém indicadores Python)
	if repo.Name != nil {
		name := strings.ToLower(*repo.Name)
		if strings.Contains(name, "python") || strings.Contains(name, "py-") || strings.HasSuffix(name, "-py") {
			return true
		}
	}
	
	return false
}

func searchRequirementsInRepo(ctx context.Context, client *github.Client, orgName, repoName string) ([]RequirementsResult, error) {
	var results []RequirementsResult

	// Buscar arquivos requirements.txt diretamente no repositório
	requirementsPaths := []string{
		"app/requirements.txt",
		"requirements.txt",
		"app/requirements/requirements.txt",
		"app/requirements/base.txt",
		"app/requirements/production.txt",
	}

	for _, path := range requirementsPaths {
		fileContent, _, _, err := client.Repositories.GetContents(ctx, orgName, repoName, path, nil)
		if err != nil {
			// Arquivo não existe, continuar
			continue
		}

		content, err := fileContent.GetContent()
		if err != nil {
			fmt.Printf("Error decoding content for %s/%s: %v\n", repoName, path, err)
			continue
		}

		hasXRay := strings.Contains(strings.ToLower(content), "aws-sdk-xray")
		hasBoto3 := strings.Contains(strings.ToLower(content), "boto3")

		results = append(results, RequirementsResult{
			Repository: repoName,
			FilePath:   path,
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
		if result.HasXRay && !result.HasBoto3 {
			filtered = append(filtered, result)
		}
	}
	
	return filtered
}

func loadCache(teamSlug string) *CacheData {
	if _, err := os.Stat(CACHE_FILE); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(CACHE_FILE)
	if err != nil {
		fmt.Printf("Error reading cache file: %v\n", err)
		return nil
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Printf("Error parsing cache file: %v\n", err)
		return nil
	}

	// Verificar se o cache é para o mesmo team
	if cache.TeamSlug != teamSlug {
		fmt.Println("Cache is for different team, ignoring...")
		return nil
	}

	return &cache
}

func saveCache(data CacheData) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling cache data: %v\n", err)
		return
	}

	if err := ioutil.WriteFile(CACHE_FILE, jsonData, 0644); err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
		return
	}

	fmt.Printf("Cache saved to %s\n", CACHE_FILE)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
