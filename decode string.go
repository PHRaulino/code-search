package main

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

// DecodeHTMLEntities decodifica entidades HTML em uma string
func DecodeHTMLEntities(s string) string {
	// Primeiro, usa o decodificador padrão do Go para entidades nomeadas
	s = html.UnescapeString(s)
	
	// Em seguida, trata entidades numéricas decimais (&#xxx;)
	decimalPattern := regexp.MustCompile(`&#(\d+);`)
	s = decimalPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extrai o número
		num := match[2 : len(match)-1]
		if code, err := strconv.Atoi(num); err == nil {
			return string(rune(code))
		}
		return match
	})
	
	// Trata entidades numéricas hexadecimais (&#xXXX;)
	hexPattern := regexp.MustCompile(`&#[xX]([0-9a-fA-F]+);`)
	s = hexPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extrai o valor hexadecimal
		hex := match[3 : len(match)-1]
		if code, err := strconv.ParseInt(hex, 16, 32); err == nil {
			return string(rune(code))
		}
		return match
	})
	
	// Trata casos sem ponto e vírgula (&#xxx ou &#xXXX)
	decimalPatternNoSemi := regexp.MustCompile(`&#(\d+)`)
	s = decimalPatternNoSemi.ReplaceAllStringFunc(s, func(match string) string {
		num := match[2:]
		if code, err := strconv.Atoi(num); err == nil {
			return string(rune(code))
		}
		return match
	})
	
	hexPatternNoSemi := regexp.MustCompile(`&#[xX]([0-9a-fA-F]+)`)
	s = hexPatternNoSemi.ReplaceAllStringFunc(s, func(match string) string {
		hex := match[3:]
		if code, err := strconv.ParseInt(hex, 16, 32); err == nil {
			return string(rune(code))
		}
		return match
	})
	
	// Trata casos com barra invertida (\xxx ou \yyy)
	backslashPattern := regexp.MustCompile(`\\(\d{1,4})`)
	s = backslashPattern.ReplaceAllStringFunc(s, func(match string) string {
		num := match[1:]
		if code, err := strconv.Atoi(num); err == nil {
			// Verifica se é um código ASCII válido
			if code >= 0 && code <= 1114111 {
				return string(rune(code))
			}
		}
		return match
	})
	
	return s
}

// DecodeCommonHTMLEntities decodifica apenas as entidades HTML mais comuns
func DecodeCommonHTMLEntities(s string) string {
	replacements := map[string]string{
		"&lt;":   "<",
		"&gt;":   ">",
		"&amp;":  "&",
		"&quot;": "\"",
		"&apos;": "'",
		"&#39;":  "'",
		"&#34;":  "\"",
		"&#38;":  "&",
		"&#60;":  "<",
		"&#62;":  ">",
		// Adiciona suporte para notação com barra invertida
		"\\60":   "<",
		"\\062":  ">",
		"\\38":   "&",
		"\\034":  "\"",
		"\\039":  "'",
	}
	
	result := s
	for entity, char := range replacements {
		result = strings.ReplaceAll(result, entity, char)
	}
	
	return result
}

// Exemplo de uso
func main() {
	// Teste com diferentes formatos
	tests := []string{
		"&lt;div&gt;Hello World&lt;/div&gt;",
		"&#60;script&#62;alert('test')&#60;/script&#62;",
		"\\060div\\062Hello\\060/div\\062",
		"Mixed: &lt;p&gt; and \\060span\\062",
		"&#x3C;h1&#x3E;Title&#x3C;/h1&#x3E;",
	}
	
	println("Usando DecodeHTMLEntities (completo):")
	for _, test := range tests {
		println("Original:", test)
		println("Decodificado:", DecodeHTMLEntities(test))
		println()
	}
	
	println("\nUsando DecodeCommonHTMLEntities (simples):")
	for _, test := range tests {
		println("Original:", test)
		println("Decodificado:", DecodeCommonHTMLEntities(test))
		println()
	}
}
