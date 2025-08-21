package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobwas/glob"
)

// TypedPattern representa um pattern com tipo e negação
type TypedPattern struct {
	Pattern   string
	Type      string // "Code", "Doc", etc.
	IsNegated bool   // true se o pattern começa com !
}

// MatchResult contém o resultado do match com tipo
type MatchResult struct {
	Matched bool
	Type    string
}

// UltraFastMatcher - Matcher super otimizado para TypedPatterns
type UltraFastMatcher struct {
	// HashMap lookups - O(1) - agora com tipos
	exactPaths     map[string]string  // path -> type
	exactBasenames map[string]string  // basename -> type
	extensions     map[string]string  // ext -> type
	
	// Slice lookups - O(n) mas rápido - agora com tipos
	prefixes       []typedPrefix
	suffixes       []typedSuffix
	
	// Para patterns complexos - com tipos
	compiledGlobs  []typedGlob
	
	// Patterns negados (!) - processados por último
	negatedPatterns []TypedPattern
	
	// Cache de resultados
	resultCache    sync.Map // map[string]MatchResult
	cacheEnabled   bool
}

type typedPrefix struct {
	prefix string
	ptype  string
}

type typedSuffix struct {
	suffix string
	ptype  string
}

type typedGlob struct {
	glob  glob.Glob
	ptype string
}

// MatcherOptions - opções de configuração
type MatcherOptions struct {
	EnableCache        bool
	CaseSensitive     bool
	MatchBasenameOnly bool
}

// NewUltraFastMatcher cria matcher com TypedPatterns
func NewUltraFastMatcher(patterns []TypedPattern, opts *MatcherOptions) (*UltraFastMatcher, error) {
	if opts == nil {
		opts = &MatcherOptions{
			EnableCache:        true,
			CaseSensitive:     true,
			MatchBasenameOnly: true,
		}
	}
	
	m := &UltraFastMatcher{
		exactPaths:     make(map[string]string),
		exactBasenames: make(map[string]string),
		extensions:     make(map[string]string),
		cacheEnabled:   opts.EnableCache,
	}
	
	if err := m.compilePatterns(patterns, opts); err != nil {
		return nil, fmt.Errorf("failed to compile patterns: %w", err)
	}
	
	return m, nil
}

// compilePatterns categoriza TypedPatterns por tipo
func (m *UltraFastMatcher) compilePatterns(patterns []TypedPattern, opts *MatcherOptions) error {
	for _, tp := range patterns {
		pattern := tp.Pattern
		if len(pattern) == 0 {
			continue
		}
		
		// Processa patterns negados separadamente
		if tp.IsNegated {
			m.negatedPatterns = append(m.negatedPatterns, tp)
			continue
		}
		
		// Case insensitive se necessário
		if !opts.CaseSensitive {
			pattern = strings.ToLower(pattern)
		}
		
		// Categoriza por tipo de pattern
		switch {
		// 1. Extensões: *.go, *.js, *.test.go
		case strings.HasPrefix(pattern, "*.") && !strings.Contains(pattern[2:], "*"):
			ext := pattern[1:] // Remove '*', mantém '.'
			// Se já existe, mantém o primeiro (precedência)
			if _, exists := m.extensions[ext]; !exists {
				m.extensions[ext] = tp.Type
			}
			
		// 2. Paths exatos: "main.go", "src/app/main.go"
		case !strings.ContainsAny(pattern, "*?[]{}"):
			if _, exists := m.exactPaths[pattern]; !exists {
				m.exactPaths[pattern] = tp.Type
			}
			if opts.MatchBasenameOnly {
				basename := filepath.Base(pattern)
				if _, exists := m.exactBasenames[basename]; !exists {
					m.exactBasenames[basename] = tp.Type
				}
			}
			
		// 3. Prefixos: "src/*", "vendor/*", "node_modules/*"
		case strings.HasSuffix(pattern, "/*") && !strings.Contains(pattern[:len(pattern)-2], "*"):
			prefix := pattern[:len(pattern)-1] // Remove '*'
			m.prefixes = append(m.prefixes, typedPrefix{prefix, tp.Type})
			
		// 4. Sufixos: "*/test", "*/tests"
		case strings.HasPrefix(pattern, "*/") && !strings.Contains(pattern[2:], "*"):
			suffix := pattern[1:] // Remove '*'
			m.suffixes = append(m.suffixes, typedSuffix{suffix, tp.Type})
			
		// 5. Patterns complexos (**, globs, etc.)
		default:
			g, err := glob.Compile(pattern)
			if err != nil {
				return fmt.Errorf("failed to compile pattern %s: %w", pattern, err)
			}
			m.compiledGlobs = append(m.compiledGlobs, typedGlob{g, tp.Type})
		}
	}
	
	return nil
}

// Match verifica se path corresponde a algum pattern e retorna tipo
func (m *UltraFastMatcher) Match(path string) MatchResult {
	if len(path) == 0 {
		return MatchResult{false, ""}
	}
	
	// 1. Verifica cache
	if m.cacheEnabled {
		if cached, ok := m.resultCache.Load(path); ok {
			return cached.(MatchResult)
		}
	}
	
	result := m.doMatch(path)
	
	// 2. Salva no cache
	if m.cacheEnabled {
		m.resultCache.Store(path, result)
	}
	
	return result
}

// doMatch executa lógica de matching otimizada
func (m *UltraFastMatcher) doMatch(path string) MatchResult {
	// 1. Exact path match - O(1)
	if ptype, exists := m.exactPaths[path]; exists {
		return m.checkNegated(path, MatchResult{true, ptype})
	}
	
	// 2. Exact basename match - O(1)
	basename := filepath.Base(path)
	if ptype, exists := m.exactBasenames[basename]; exists {
		return m.checkNegated(path, MatchResult{true, ptype})
	}
	
	// 3. Extension match - O(1)
	if ext := filepath.Ext(path); ext != "" {
		if ptype, exists := m.extensions[ext]; exists {
			return m.checkNegated(path, MatchResult{true, ptype})
		}
	}
	
	// 4. Extensões compostas (.test.go, .min.js)
	if dotIndex := strings.LastIndex(basename, "."); dotIndex > 0 {
		if secondDot := strings.LastIndex(basename[:dotIndex], "."); secondDot > 0 {
			compoundExt := basename[secondDot:]
			if ptype, exists := m.extensions[compoundExt]; exists {
				return m.checkNegated(path, MatchResult{true, ptype})
			}
		}
	}
	
	// 5. Prefix match - O(n)
	for _, tp := range m.prefixes {
		if strings.HasPrefix(path, tp.prefix) {
			return m.checkNegated(path, MatchResult{true, tp.ptype})
		}
	}
	
	// 6. Suffix match - O(n)
	for _, ts := range m.suffixes {
		if strings.HasSuffix(path, ts.suffix) {
			return m.checkNegated(path, MatchResult{true, ts.ptype})
		}
	}
	
	// 7. Complex glob patterns
	for _, tg := range m.compiledGlobs {
		if tg.glob.Match(path) || tg.glob.Match(basename) {
			return m.checkNegated(path, MatchResult{true, tg.ptype})
		}
	}
	
	return MatchResult{false, ""}
}

// checkNegated verifica se algum pattern negado (!pattern) cancela o match
func (m *UltraFastMatcher) checkNegated(path string, result MatchResult) MatchResult {
	if !result.Matched {
		return result
	}
	
	basename := filepath.Base(path)
	
	for _, negated := range m.negatedPatterns {
		pattern := negated.Pattern
		
		// Remove ! do início se existir
		if strings.HasPrefix(pattern, "!") {
			pattern = pattern[1:]
		}
		
		// Verifica se o pattern negado faz match
		matched := false
		switch {
		case strings.HasPrefix(pattern, "*.") && !strings.Contains(pattern[2:], "*"):
			ext := pattern[1:]
			matched = strings.HasSuffix(path, ext) || strings.HasSuffix(basename, ext)
			
		case !strings.ContainsAny(pattern, "*?[]{}"):
			matched = (path == pattern) || (basename == pattern)
			
		case strings.HasSuffix(pattern, "/*"):
			prefix := pattern[:len(pattern)-1]
			matched = strings.HasPrefix(path, prefix)
			
		case strings.HasPrefix(pattern, "*/"):
			suffix := pattern[1:]
			matched = strings.HasSuffix(path, suffix)
			
		default:
			if g, err := glob.Compile(pattern); err == nil {
				matched = g.Match(path) || g.Match(basename)
			}
		}
		
		// Se pattern negado faz match, cancela o resultado
		if matched {
			return MatchResult{false, ""}
		}
	}
	
	return result
}

// MatchBatch processa múltiplos paths
func (m *UltraFastMatcher) MatchBatch(paths []string) []MatchResult {
	results := make([]MatchResult, len(paths))
	for i, path := range paths {
		results[i] = m.Match(path)
	}
	return results
}

// MatchSimple retorna apenas bool (compatibilidade)
func (m *UltraFastMatcher) MatchSimple(path string) bool {
	return m.Match(path).Matched
}

// ClearCache limpa cache de resultados
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
		ExactPaths:      len(m.exactPaths),
		ExactBasenames:  len(m.exactBasenames),
		Extensions:      len(m.extensions),
		Prefixes:        len(m.prefixes),
		Suffixes:        len(m.suffixes),
		ComplexGlobs:    len(m.compiledGlobs),
		NegatedPatterns: len(m.negatedPatterns),
		CacheSize:       cacheSize,
	}
}

type MatcherStats struct {
	ExactPaths      int
	ExactBasenames  int
	Extensions      int
	Prefixes        int
	Suffixes        int
	ComplexGlobs    int
	NegatedPatterns int
	CacheSize       int
}

func (s MatcherStats) String() string {
	return fmt.Sprintf(
		"MatcherStats{ExactPaths: %d, ExactBasenames: %d, Extensions: %d, "+
		"Prefixes: %d, Suffixes: %d, ComplexGlobs: %d, NegatedPatterns: %d, CacheSize: %d}",
		s.ExactPaths, s.ExactBasenames, s.Extensions, 
		s.Prefixes, s.Suffixes, s.ComplexGlobs, s.NegatedPatterns, s.CacheSize,
	)
}

// === EXEMPLO DE USO ===
func main() {
	// Patterns com tipos
	patterns := []TypedPattern{
		// Código
		{Pattern: "*.go", Type: "Code", IsNegated: false},
		{Pattern: "*.js", Type: "Code", IsNegated: false},
		{Pattern: "*.test.go", Type: "Code", IsNegated: false},
		{Pattern: "src/*", Type: "Code", IsNegated: false},
		
		// Documentação
		{Pattern: "*.md", Type: "Doc", IsNegated: false},
		{Pattern: "*.txt", Type: "Doc", IsNegated: false},
		{Pattern: "docs/*", Type: "Doc", IsNegated: false},
		
		// Ignorados (negados)
		{Pattern: "!*.test.go", Type: "", IsNegated: true},
		{Pattern: "!node_modules/*", Type: "", IsNegated: true},
		{Pattern: "!.git/*", Type: "", IsNegated: true},
		
		// Outros
		{Pattern: "*.log", Type: "Log", IsNegated: false},
		{Pattern: "*.tmp", Type: "Temp", IsNegated: false},
	}
	
	// Cria matcher
	opts := &MatcherOptions{
		EnableCache:        true,
		CaseSensitive:     true,
		MatchBasenameOnly: true,
	}
	
	matcher, err := NewUltraFastMatcher(patterns, opts)
	if err != nil {
		panic(err)
	}
	
	// Testa paths
	testPaths := []string{
		"main.go",                 // Code
		"app.test.go",            // Code mas negado por !*.test.go
		"README.md",              // Doc
		"src/utils.js",           // Code (src/*)
		"docs/guide.txt",         // Doc (docs/*)
		"node_modules/lib.js",    // Negado por !node_modules/*
		"app.log",                // Log
		"temp.tmp",               // Temp
		".git/config",            // Negado por !.git/*
	}
	
	fmt.Println("🚀 UltraFastMatcher com Tipos:")
	fmt.Println(strings.Repeat("=", 50))
	
	for _, path := range testPaths {
		result := matcher.Match(path)
		status := "❌"
		typeInfo := ""
		
		if result.Matched {
			status = "✅"
			typeInfo = fmt.Sprintf(" [%s]", result.Type)
		}
		
		fmt.Printf("%s %s%s\n", status, path, typeInfo)
	}
	
	fmt.Println("\n📊 Estatísticas:")
	fmt.Println(matcher.Stats())
	
	// Batch processing
	fmt.Println("\n⚡ Processamento em Batch:")
	results := matcher.MatchBatch(testPaths)
	
	typeCount := make(map[string]int)
	for _, result := range results {
		if result.Matched {
			typeCount[result.Type]++
		}
	}
	
	fmt.Printf("Processados: %d paths\n", len(testPaths))
	for ptype, count := range typeCount {
		fmt.Printf("  %s: %d matches\n", ptype, count)
	}
}
