package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobwas/glob"
)

// UltraFastMatcher - Matcher super otimizado para diferentes tipos de patterns
type UltraFastMatcher struct {
	// HashMap lookups - O(1)
	exactPaths     map[string]bool  // "main.go", "src/app.go"
	exactBasenames map[string]bool  // Para match apenas do nome do arquivo
	extensions     map[string]bool  // ".go", ".js", ".test.go"
	
	// Slice lookups - O(n) mas n é pequeno e operações são muito rápidas
	prefixes       []string         // "src/", "vendor/", "node_modules/"
	suffixes       []string         // "/test", "/tests"
	
	// Para patterns complexos - último recurso
	compiledGlobs  []glob.Glob      // Para patterns que precisam de glob
	
	// Cache de resultados - para paths consultados frequentemente
	resultCache    sync.Map         // map[string]bool
	cacheEnabled   bool
}

// Opções de configuração
type MatcherOptions struct {
	EnableCache        bool
	CaseSensitive     bool
	MatchBasenameOnly bool  // Se deve testar apenas o nome do arquivo também
}

// NewUltraFastMatcher cria um novo matcher otimizado
func NewUltraFastMatcher(patterns []string, opts *MatcherOptions) (*UltraFastMatcher, error) {
	if opts == nil {
		opts = &MatcherOptions{
			EnableCache:        true,
			CaseSensitive:     true,
			MatchBasenameOnly: true,
		}
	}
	
	m := &UltraFastMatcher{
		exactPaths:     make(map[string]bool),
		exactBasenames: make(map[string]bool),
		extensions:     make(map[string]bool),
		cacheEnabled:   opts.EnableCache,
	}
	
	// Pré-processa e categoriza patterns
	if err := m.compilePatterns(patterns, opts); err != nil {
		return nil, fmt.Errorf("failed to compile patterns: %w", err)
	}
	
	return m, nil
}

// compilePatterns categoriza e otimiza patterns por tipo
func (m *UltraFastMatcher) compilePatterns(patterns []string, opts *MatcherOptions) error {
	for _, pattern := range patterns {
		if len(pattern) == 0 {
			continue
		}
		
		// Case insensitive se necessário
		if !opts.CaseSensitive {
			pattern = strings.ToLower(pattern)
		}
		
		// Categoriza o pattern pelo tipo para máxima otimização
		switch {
		
		// 1. Extensões simples: *.go, *.js, *.test.go
		case strings.HasPrefix(pattern, "*.") && !strings.Contains(pattern[2:], "*"):
			ext := pattern[1:] // Remove apenas o '*', mantém o '.'
			m.extensions[ext] = true
			
		// 2. Paths exatos: "main.go", "src/app/main.go"
		case !strings.ContainsAny(pattern, "*?[]{}"):
			m.exactPaths[pattern] = true
			if opts.MatchBasenameOnly {
				basename := filepath.Base(pattern)
				m.exactBasenames[basename] = true
			}
			
		// 3. Prefixos simples: "src/*", "vendor/*", "node_modules/*"
		case strings.HasSuffix(pattern, "/*") && !strings.Contains(pattern[:len(pattern)-2], "*"):
			prefix := pattern[:len(pattern)-1] // Remove apenas o '*'
			m.prefixes = append(m.prefixes, prefix)
			
		// 4. Sufixos simples: "*/test", "*/tests"  
		case strings.HasPrefix(pattern, "*/") && !strings.Contains(pattern[2:], "*"):
			suffix := pattern[1:] // Remove apenas o '*'
			m.suffixes = append(m.suffixes, suffix)
			
		// 5. Patterns com ** (doublestar) - converte para glob
		case strings.Contains(pattern, "**"):
			// Substitui ** por ** para glob
			globPattern := strings.ReplaceAll(pattern, "**", "**")
			g, err := glob.Compile(globPattern)
			if err != nil {
				return fmt.Errorf("failed to compile doublestar pattern %s: %w", pattern, err)
			}
			m.compiledGlobs = append(m.compiledGlobs, g)
			
		// 6. Outros patterns complexos
		default:
			g, err := glob.Compile(pattern)
			if err != nil {
				return fmt.Errorf("failed to compile glob pattern %s: %w", pattern, err)
			}
			m.compiledGlobs = append(m.compiledGlobs, g)
		}
	}
	
	return nil
}

// Match verifica se um path corresponde a algum pattern
func (m *UltraFastMatcher) Match(path string) bool {
	if len(path) == 0 {
		return false
	}
	
	// 1. Verifica cache primeiro
	if m.cacheEnabled {
		if cached, ok := m.resultCache.Load(path); ok {
			return cached.(bool)
		}
	}
	
	result := m.doMatch(path)
	
	// Salva no cache
	if m.cacheEnabled {
		m.resultCache.Store(path, result)
	}
	
	return result
}

// doMatch executa a lógica de matching otimizada
func (m *UltraFastMatcher) doMatch(path string) bool {
	// 1. Exact path match - O(1) HashMap lookup
	if m.exactPaths[path] {
		return true
	}
	
	// 2. Exact basename match - O(1) HashMap lookup  
	basename := filepath.Base(path)
	if m.exactBasenames[basename] {
		return true
	}
	
	// 3. Extension match - O(1) HashMap lookup
	if ext := filepath.Ext(path); ext != "" && m.extensions[ext] {
		return true
	}
	
	// Também verifica extensões compostas como .test.go
	if dotIndex := strings.LastIndex(basename, "."); dotIndex > 0 {
		// Para "app.test.go", verifica ".test.go" 
		if secondDot := strings.LastIndex(basename[:dotIndex], "."); secondDot > 0 {
			compoundExt := basename[secondDot:]
			if m.extensions[compoundExt] {
				return true
			}
		}
	}
	
	// 4. Prefix match - O(n) mas muito rápido
	for _, prefix := range m.prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	
	// 5. Suffix match - O(n) mas muito rápido
	for _, suffix := range m.suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	
	// 6. Complex glob patterns - último recurso
	for _, g := range m.compiledGlobs {
		if g.Match(path) || g.Match(basename) {
			return true
		}
	}
	
	return false
}

// MatchBatch verifica múltiplos paths de uma vez - ainda mais otimizado
func (m *UltraFastMatcher) MatchBatch(paths []string) []bool {
	results := make([]bool, len(paths))
	
	for i, path := range paths {
		results[i] = m.Match(path)
	}
	
	return results
}

// ClearCache limpa o cache de resultados
func (m *UltraFastMatcher) ClearCache() {
	if m.cacheEnabled {
		m.resultCache.Range(func(key, value interface{}) bool {
			m.resultCache.Delete(key)
			return true
		})
	}
}

// Stats retorna estatísticas do matcher
func (m *UltraFastMatcher) Stats() MatcherStats {
	cacheSize := 0
	if m.cacheEnabled {
		m.resultCache.Range(func(k, v interface{}) bool {
			cacheSize++
			return true
		})
	}
	
	return MatcherStats{
		ExactPaths:     len(m.exactPaths),
		ExactBasenames: len(m.exactBasenames),
		Extensions:     len(m.extensions),
		Prefixes:       len(m.prefixes),
		Suffixes:       len(m.suffixes),
		ComplexGlobs:   len(m.compiledGlobs),
		CacheSize:      cacheSize,
	}
}

type MatcherStats struct {
	ExactPaths     int
	ExactBasenames int
	Extensions     int
	Prefixes       int
	Suffixes       int
	ComplexGlobs   int
	CacheSize      int
}

func (s MatcherStats) String() string {
	return fmt.Sprintf(
		"MatcherStats{ExactPaths: %d, ExactBasenames: %d, Extensions: %d, "+
		"Prefixes: %d, Suffixes: %d, ComplexGlobs: %d, CacheSize: %d}",
		s.ExactPaths, s.ExactBasenames, s.Extensions, 
		s.Prefixes, s.Suffixes, s.ComplexGlobs, s.CacheSize,
	)
}

// === EXEMPLO DE USO ===
func main() {
	// Patterns típicos de .gitignore
	patterns := []string{
		// Extensões
		"*.log",
		"*.tmp", 
		"*.test.go",
		"*.min.js",
		
		// Paths exatos
		".DS_Store",
		"Thumbs.db",
		"main.go",
		
		// Prefixos  
		"node_modules/*",
		"vendor/*",
		".git/*",
		"dist/*",
		"build/*",
		
		// Sufixos
		"*/tmp",
		"*/cache", 
		
		// Patterns complexos
		"**/*.test.js",
		"**/node_modules/**",
		"**/.git/**",
	}
	
	// Cria o matcher
	opts := &MatcherOptions{
		EnableCache:        true,
		CaseSensitive:     true, 
		MatchBasenameOnly: true,
	}
	
	matcher, err := NewUltraFastMatcher(patterns, opts)
	if err != nil {
		panic(err)
	}
	
	// Testa alguns paths
	testPaths := []string{
		"app.log",                    // *.log
		"src/main.go",               // main.go (basename)
		"tests/unit.test.go",        // *.test.go
		"node_modules/react/index.js", // node_modules/*
		"vendor/pkg/lib.go",         // vendor/*
		".git/config",               // .git/*
		"dist/bundle.min.js",        // dist/* e *.min.js
		"src/cache/data",            // */cache
		"deep/path/to/file.test.js", // **/*.test.js
		"project/.git/hooks/pre-commit", // **/.git/**
	}
	
	fmt.Println("🚀 UltraFastMatcher Results:")
	fmt.Println("=" * 40)
	
	for _, path := range testPaths {
		matched := matcher.Match(path)
		status := "❌"
		if matched {
			status = "✅"
		}
		fmt.Printf("%s %s\n", status, path)
	}
	
	fmt.Println("\n📊 Matcher Statistics:")
	fmt.Println(matcher.Stats())
	
	// Teste de performance com batch
	fmt.Println("\n⚡ Batch Processing:")
	results := matcher.MatchBatch(testPaths)
	matchCount := 0
	for _, matched := range results {
		if matched {
			matchCount++
		}
	}
	fmt.Printf("Processed %d paths, %d matches\n", len(testPaths), matchCount)
}

// Para instalar a dependência:
// go get github.com/gobwas/glob
