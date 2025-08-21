package main

import (
	"sync"
)

// Integração com seu Service existente
type Service struct {
	// Seus campos existentes...
	
	// Novo matcher ultra rápido
	matcher     *UltraFastMatcher
	matcherOnce sync.Once
	patterns    []string // Cache dos patterns para lazy loading
}

// Substituição da função original
func (s *Service) shouldIncludeFile(relativePath string, patterns []string) (bool, PatternType) {
	// Lazy initialization do matcher
	s.matcherOnce.Do(func() {
		opts := &MatcherOptions{
			EnableCache:        true,
			CaseSensitive:     true,
			MatchBasenameOnly: true,
		}
		
		matcher, err := NewUltraFastMatcher(patterns, opts)
		if err != nil {
			// Fallback para implementação original se houver erro
			s.matcher = nil
			return
		}
		s.matcher = matcher
		s.patterns = patterns
	})
	
	// Se o matcher foi inicializado com sucesso, usa ele
	if s.matcher != nil {
		matched := s.matcher.Match(relativePath)
		if matched {
			return true, PatternTypeMatch // ou o tipo que você usar
		}
		return false, PatternTypeNone
	}
	
	// Fallback para implementação original
	return s.shouldIncludeFileFallback(relativePath, patterns)
}

// Versão otimizada da função matchPattern original
func (s *Service) matchPattern(path, pattern string) bool {
	// Se temos matcher inicializado, cria um temporário para este pattern
	if matcher, err := NewUltraFastMatcher([]string{pattern}, nil); err == nil {
		return matcher.Match(path)
	}
	
	// Fallback para implementação original
	return s.matchPatternFallback(path, pattern)
}

// Fallbacks para compatibilidade (suas funções originais)
func (s *Service) shouldIncludeFileFallback(relativePath string, patterns []string) (bool, PatternType) {
	// Sua implementação original aqui...
	return false, PatternTypeNone
}

func (s *Service) matchPatternFallback(path, pattern string) bool {
	// Sua implementação original aqui...
	return false
}

// Método utilitário para recarregar patterns
func (s *Service) UpdatePatterns(newPatterns []string) error {
	opts := &MatcherOptions{
		EnableCache:        true,
		CaseSensitive:     true,
		MatchBasenameOnly: true,
	}
	
	matcher, err := NewUltraFastMatcher(newPatterns, opts)
	if err != nil {
		return err
	}
	
	s.matcher = matcher
	s.patterns = newPatterns
	return nil
}

// Método para obter estatísticas
func (s *Service) GetMatcherStats() *MatcherStats {
	if s.matcher != nil {
		stats := s.matcher.Stats()
		return &stats
	}
	return nil
}

// Exemplo de uso específico para seu cenário
func (s *Service) ProcessFiles(filePaths []string, patterns []string) []string {
	// Inicializa matcher se necessário
	if s.matcher == nil {
		if err := s.UpdatePatterns(patterns); err != nil {
			// Handle error - talvez log e use fallback
			return s.processFilesFallback(filePaths, patterns)
		}
	}
	
	var matchedFiles []string
	
	// Processa em batch para máxima performance
	results := s.matcher.MatchBatch(filePaths)
	
	for i, matched := range results {
		if matched {
			matchedFiles = append(matchedFiles, filePaths[i])
		}
	}
	
	return matchedFiles
}

func (s *Service) processFilesFallback(filePaths []string, patterns []string) []string {
	// Implementação fallback usando seu código original
	var matchedFiles []string
	// ... sua lógica original
	return matchedFiles
}

// === EXEMPLO DE BENCHMARK ===
/*
func BenchmarkMatchPattern(b *testing.B) {
	patterns := []string{
		"*.go", "*.js", "*.test.go", 
		"node_modules/*", "vendor/*", 
		"**/*.min.js", "**/.git/**",
	}
	
	testPaths := []string{
		"main.go", "src/app.js", "test.test.go",
		"node_modules/react/index.js", "vendor/lib.go",
		"dist/app.min.js", ".git/config",
	}
	
	// Benchmark original
	service := &Service{}
	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range testPaths {
				for _, pattern := range patterns {
					service.matchPatternFallback(path, pattern)
				}
			}
		}
	})
	
	// Benchmark UltraFast
	matcher, _ := NewUltraFastMatcher(patterns, nil)
	b.Run("UltraFast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, path := range testPaths {
				matcher.Match(path)
			}
		}
	})
	
	// Benchmark UltraFast Batch
	b.Run("UltraFastBatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			matcher.MatchBatch(testPaths)
		}
	})
}
*/

// === UTILITÁRIOS ADICIONAIS ===

// Para debug - mostra como cada pattern foi categorizado
func AnalyzePatternOptimization(patterns []string) {
	matcher, err := NewUltraFastMatcher(patterns, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	stats := matcher.Stats()
	total := len(patterns)
	
	fmt.Printf("📊 Pattern Optimization Analysis:\n")
	fmt.Printf("Total patterns: %d\n\n", total)
	
	fmt.Printf("✅ Super Fast (O(1) HashMap):\n")
	fmt.Printf("  • Exact paths: %d\n", stats.ExactPaths)
	fmt.Printf("  • Exact basenames: %d\n", stats.ExactBasenames) 
	fmt.Printf("  • Extensions: %d\n", stats.Extensions)
	
	fmt.Printf("\n⚡ Very Fast (O(n) but n small):\n")
	fmt.Printf("  • Prefixes: %d\n", stats.Prefixes)
	fmt.Printf("  • Suffixes: %d\n", stats.Suffixes)
	
	fmt.Printf("\n🐌 Slower (Glob patterns):\n")
	fmt.Printf("  • Complex globs: %d\n", stats.ComplexGlobs)
	
	optimized := stats.ExactPaths + stats.ExactBasenames + stats.Extensions + stats.Prefixes + stats.Suffixes
	percentage := float64(optimized) / float64(total) * 100
	
	fmt.Printf("\n🎯 Optimization Rate: %.1f%% (%d/%d patterns optimized)\n", 
		percentage, optimized, total)
}
