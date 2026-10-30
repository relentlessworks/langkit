package model

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// DetectionResult holds the result of language detection.
type DetectionResult struct {
	Language    string  `json:"language"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	NativeName  string  `json:"native_name"`
	Confidence  float64 `json:"confidence"`
	Script      string  `json:"script"`
	Alternatives []AltResult `json:"alternatives"`
}

// AltResult is an alternative language candidate.
type AltResult struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// BatchResult holds results for batch detection.
type BatchResult struct {
	Index    int             `json:"index"`
	Text     string          `json:"text"`
	Result   DetectionResult `json:"result"`
}

// Detect identifies the language of the given text.
// It uses a combination of script detection and common word matching.
func Detect(text string) DetectionResult {
	text = strings.TrimSpace(text)
	if text == "" {
		return DetectionResult{
			Language:   "unknown",
			Code:       "und",
			Name:       "Unknown",
			Confidence: 0,
			Script:     "Unknown",
		}
	}

	// Step 1: Detect script
	script, scriptPct := DetectScript(text)

	// Step 2: Get candidate languages based on script
	candidates := candidatesByScript(script)

	// Step 3: Score each candidate using common word matching
	scores := scoreByCommonWords(text, candidates)

	// Step 4: If no common word matches, use script-only detection
	if len(scores) == 0 {
		return scriptOnlyResult(script, scriptPct)
	}

	// Step 5: Sort by score and build result
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	totalScore := 0.0
	for _, s := range scores {
		totalScore += s.score
	}

	// Build alternatives
	var alts []AltResult
	for i, s := range scores {
		confidence := s.score / totalScore * 100
		if i == 0 {
			continue // skip top, it's the main result
		}
		if confidence < 1.0 {
			break
		}
		lang, _ := LanguageByCode(s.code)
		alts = append(alts, AltResult{
			Code:       s.code,
			Name:       lang.Name,
			Confidence: round(confidence, 2),
		})
		if len(alts) >= 5 {
			break
		}
	}

	top := scores[0]
	topConfidence := top.score / totalScore * 100
	lang, _ := LanguageByCode(top.code)

	return DetectionResult{
		Language:     lang.Code,
		Code:         lang.Code,
		Name:         lang.Name,
		NativeName:   lang.NativeName,
		Confidence:   round(topConfidence, 2),
		Script:       script,
		Alternatives: alts,
	}
}

type langScore struct {
	code  string
	score float64
}

// candidatesByScript returns language codes that use the given script.
func candidatesByScript(script string) []string {
	if script == "Unknown" || script == "" {
		// Return all Latin-script languages as fallback
		var codes []string
		for _, lang := range SupportedLanguages {
			for _, s := range lang.Scripts {
				if s == "Latin" {
					codes = append(codes, lang.Code)
					break
				}
			}
		}
		return codes
	}

	var codes []string
	for _, lang := range SupportedLanguages {
		for _, s := range lang.Scripts {
			if s == script {
				codes = append(codes, lang.Code)
				break
			}
		}
	}
	return codes
}

// scoreByCommonWords scores candidate languages by matching common words.
func scoreByCommonWords(text string, candidates []string) []langScore {
	// Tokenize text into lowercase words
	words := tokenize(text)
	if len(words) == 0 {
		return nil
	}

	wordSet := make(map[string]bool)
	for _, w := range words {
		wordSet[w] = true
	}

	var scores []langScore
	for _, code := range candidates {
		common, ok := commonWords[code]
		if !ok {
			continue
		}
		score := 0.0
		matched := 0
		for w := range common {
			if wordSet[w] {
				score += 1.0
				matched++
			}
		}
		if score > 0 {
			// Normalize by text length to avoid bias toward longer texts
			scores = append(scores, langScore{
				code:  code,
				score: score / float64(len(words)),
			})
		}
	}
	return scores
}

// scriptOnlyResult builds a result when only script detection is available.
func scriptOnlyResult(script string, pct float64) DetectionResult {
	// Map script to primary language
	scriptToLang := map[string]string{
		"Cyrillic":   "ru",
		"Greek":      "el",
		"Arabic":     "ar",
		"Hebrew":     "he",
		"Devanagari": "hi",
		"Bengali":    "bn",
		"Tamil":      "ta",
		"Telugu":     "te",
		"Malayalam":  "ml",
		"Kannada":    "kn",
		"Gujarati":   "gu",
		"Gurmukhi":   "pa",
		"Sinhala":    "si",
		"Thai":       "th",
		"Lao":        "lo",
		"Khmer":      "km",
		"Myanmar":    "my",
		"Georgian":   "ka",
		"Ethiopic":   "am",
		"Hangul":     "ko",
		"Hiragana":   "ja",
		"Katakana":   "ja",
		"Han":        "zh",
	}

	code, ok := scriptToLang[script]
	if !ok {
		return DetectionResult{
			Language:   "unknown",
			Code:       "und",
			Name:       "Unknown",
			Confidence: 0,
			Script:     script,
		}
	}

	lang, _ := LanguageByCode(code)
	return DetectionResult{
		Language:   lang.Code,
		Code:       lang.Code,
		Name:       lang.Name,
		NativeName: lang.NativeName,
		Confidence: round(pct, 2),
		Script:     script,
	}
}

// tokenize splits text into lowercase word tokens.
func tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || r == '\'' || r == '-' {
			current.WriteRune(unicode.ToLower(r))
		} else {
			if current.Len() > 0 {
				w := current.String()
				if len(w) >= 2 {
					words = append(words, w)
				}
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		w := current.String()
		if len(w) >= 2 {
			words = append(words, w)
		}
	}
	return words
}

// round rounds to n decimal places.
func round(v float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(v*pow) / pow
}

// DetectBatch detects languages for multiple texts.
func DetectBatch(texts []string) []BatchResult {
	results := make([]BatchResult, len(texts))
	for i, text := range texts {
		results[i] = BatchResult{
			Index:  i,
			Text:   text,
			Result: Detect(text),
		}
	}
	return results
}

// SupportedScripts returns all script names langkit can detect.
func SupportedScripts() []string {
	var names []string
	for _, s := range scripts {
		names = append(names, s.Name)
	}
	return names
}
