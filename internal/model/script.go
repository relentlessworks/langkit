package model

// Script represents a Unicode writing system.
type Script struct {
	Name  string
	Ranges []runeRange
}

type runeRange struct {
	lo, hi rune
}

func r(lo, hi rune) runeRange { return runeRange{lo, hi} }

// Unicode script ranges for detection.
var scripts = []Script{
	{"Latin", []runeRange{
		r(0x0041, 0x005A), r(0x0061, 0x007A), // Basic Latin
		r(0x00C0, 0x00FF), // Latin-1 Supplement
		r(0x0100, 0x017F), // Latin Extended-A
		r(0x0180, 0x024F), // Latin Extended-B
		r(0x0250, 0x02AF), // IPA Extensions
		r(0x1E00, 0x1EFF), // Latin Extended Additional
	}},
	{"Cyrillic", []runeRange{
		r(0x0400, 0x04FF), // Cyrillic
		r(0x0500, 0x052F), // Cyrillic Supplement
		r(0x2DE0, 0x2DFF), // Cyrillic Extended-A
		r(0xA640, 0xA69F), // Cyrillic Extended-B
	}},
	{"Greek", []runeRange{
		r(0x0370, 0x03FF), // Greek and Coptic
		r(0x1F00, 0x1FFF), // Greek Extended
	}},
	{"Arabic", []runeRange{
		r(0x0600, 0x06FF), // Arabic
		r(0x0750, 0x077F), // Arabic Supplement
		r(0x08A0, 0x08FF), // Arabic Extended-A
		r(0xFB50, 0xFDFF), // Arabic Presentation Forms-A
		r(0xFE70, 0xFEFF), // Arabic Presentation Forms-B
	}},
	{"Hebrew", []runeRange{
		r(0x0590, 0x05FF), // Hebrew
		r(0xFB1D, 0xFB4F), // Hebrew Presentation Forms
	}},
	{"Devanagari", []runeRange{
		r(0x0900, 0x097F), // Devanagari
		r(0xA8E0, 0xA8FF), // Devanagari Extended
	}},
	{"Bengali", []runeRange{r(0x0980, 0x09FF)}},
	{"Tamil", []runeRange{r(0x0B80, 0x0BFF)}},
	{"Telugu", []runeRange{r(0x0C00, 0x0C7F)}},
	{"Malayalam", []runeRange{r(0x0D00, 0x0D7F)}},
	{"Kannada", []runeRange{r(0x0C80, 0x0CFF)}},
	{"Gujarati", []runeRange{r(0x0A80, 0x0AFF)}},
	{"Gurmukhi", []runeRange{r(0x0A00, 0x0A7F)}},
	{"Sinhala", []runeRange{r(0x0D80, 0x0DFF)}},
	{"Thai", []runeRange{r(0x0E00, 0x0E7F)}},
	{"Lao", []runeRange{r(0x0E80, 0x0EFF)}},
	{"Khmer", []runeRange{r(0x1780, 0x17FF)}},
	{"Myanmar", []runeRange{r(0x1000, 0x109F)}},
	{"Georgian", []runeRange{r(0x10A0, 0x10FF)}},
	{"Ethiopic", []runeRange{
		r(0x1200, 0x137F),
		r(0x1380, 0x139F),
		r(0x2D80, 0x2DDF),
	}},
	{"Hangul", []runeRange{
		r(0xAC00, 0xD7AF), // Hangul Syllables
		r(0x1100, 0x11FF), // Hangul Jamo
		r(0x3130, 0x318F), // Hangul Compatibility Jamo
	}},
	{"Hiragana", []runeRange{r(0x3040, 0x309F)}},
	{"Katakana", []runeRange{
		r(0x30A0, 0x30FF),
		r(0xFF65, 0xFF9F), // Halfwidth Katakana
	}},
	{"Han", []runeRange{
		r(0x4E00, 0x9FFF), // CJK Unified Ideographs
		r(0x3400, 0x4DBF), // CJK Extension A
		r(0x20000, 0x2A6DF), // CJK Extension B
		r(0xF900, 0xFAFF), // CJK Compatibility Ideographs
	}},
}

// DetectScript identifies the dominant script in the given text.
// Returns the script name and the percentage of characters in that script.
func DetectScript(text string) (string, float64) {
	counts := make(map[string]int)
	total := 0
	for _, ch := range text {
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			continue
		}
		for _, s := range scripts {
			for _, rng := range s.Ranges {
				if ch >= rng.lo && ch <= rng.hi {
					counts[s.Name]++
					total++
					break
				}
			}
		}
	}
	if total == 0 {
		return "Unknown", 0
	}
	bestScript := ""
	bestCount := 0
	for name, count := range counts {
		if count > bestCount {
			bestCount = count
			bestScript = name
		}
	}
	return bestScript, float64(bestCount) / float64(total) * 100
}

// DetectAllScripts returns counts for all scripts found in the text.
func DetectAllScripts(text string) map[string]int {
	counts := make(map[string]int)
	for _, ch := range text {
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			continue
		}
		for _, s := range scripts {
			for _, rng := range s.Ranges {
				if ch >= rng.lo && ch <= rng.hi {
					counts[s.Name]++
					break
				}
			}
		}
	}
	return counts
}
