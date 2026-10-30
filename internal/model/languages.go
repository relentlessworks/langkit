package model

// Language represents a supported language for detection.
type Language struct {
	Code       string   // ISO 639-1 code (e.g. "en", "fr", "de")
	Name       string   // English name
	NativeName string   // Native name
	Scripts    []string // Writing scripts (Latin, Cyrillic, Arabic, etc.)
}

// SupportedLanguages is the full list of languages langkit can detect.
var SupportedLanguages = []Language{
	{"en", "English", "English", []string{"Latin"}},
	{"es", "Spanish", "Español", []string{"Latin"}},
	{"fr", "French", "Français", []string{"Latin"}},
	{"de", "German", "Deutsch", []string{"Latin"}},
	{"it", "Italian", "Italiano", []string{"Latin"}},
	{"pt", "Portuguese", "Português", []string{"Latin"}},
	{"nl", "Dutch", "Nederlands", []string{"Latin"}},
	{"ru", "Russian", "Русский", []string{"Cyrillic"}},
	{"uk", "Ukrainian", "Українська", []string{"Cyrillic"}},
	{"pl", "Polish", "Polski", []string{"Latin"}},
	{"tr", "Turkish", "Türkçe", []string{"Latin"}},
	{"sv", "Swedish", "Svenska", []string{"Latin"}},
	{"no", "Norwegian", "Norsk", []string{"Latin"}},
	{"da", "Danish", "Dansk", []string{"Latin"}},
	{"fi", "Finnish", "Suomi", []string{"Latin"}},
	{"cs", "Czech", "Čeština", []string{"Latin"}},
	{"sk", "Slovak", "Slovenčina", []string{"Latin"}},
	{"hu", "Hungarian", "Magyar", []string{"Latin"}},
	{"ro", "Romanian", "Română", []string{"Latin"}},
	{"bg", "Bulgarian", "Български", []string{"Cyrillic"}},
	{"el", "Greek", "Ελληνικά", []string{"Greek"}},
	{"ar", "Arabic", "العربية", []string{"Arabic"}},
	{"he", "Hebrew", "עברית", []string{"Hebrew"}},
	{"fa", "Persian", "فارسی", []string{"Arabic"}},
	{"ur", "Urdu", "اردو", []string{"Arabic"}},
	{"hi", "Hindi", "हिन्दी", []string{"Devanagari"}},
	{"bn", "Bengali", "বাংলা", []string{"Bengali"}},
	{"ta", "Tamil", "தமிழ்", []string{"Tamil"}},
	{"te", "Telugu", "తెలుగు", []string{"Telugu"}},
	{"ml", "Malayalam", "മലയാളം", []string{"Malayalam"}},
	{"kn", "Kannada", "ಕನ್ನಡ", []string{"Kannada"}},
	{"mr", "Marathi", "मराठी", []string{"Devanagari"}},
	{"ne", "Nepali", "नेपाली", []string{"Devanagari"}},
	{"gu", "Gujarati", "ગુજરાતી", []string{"Gujarati"}},
	{"pa", "Punjabi", "ਪੰਜਾਬੀ", []string{"Gurmukhi"}},
	{"si", "Sinhala", "සිංහල", []string{"Sinhala"}},
	{"th", "Thai", "ไทย", []string{"Thai"}},
	{"lo", "Lao", "ລາວ", []string{"Lao"}},
	{"km", "Khmer", "ខ្មែរ", []string{"Khmer"}},
	{"my", "Burmese", "မြန်မာ", []string{"Myanmar"}},
	{"ka", "Georgian", "ქართული", []string{"Georgian"}},
	{"am", "Amharic", "አማርኛ", []string{"Ethiopic"}},
	{"ja", "Japanese", "日本語", []string{"Hiragana", "Katakana", "Han"}},
	{"ko", "Korean", "한국어", []string{"Hangul"}},
	{"zh", "Chinese", "中文", []string{"Han"}},
	{"vi", "Vietnamese", "Tiếng Việt", []string{"Latin"}},
	{"id", "Indonesian", "Bahasa Indonesia", []string{"Latin"}},
	{"ms", "Malay", "Bahasa Melayu", []string{"Latin"}},
	{"tl", "Filipino", "Filipino", []string{"Latin"}},
	{"sw", "Swahili", "Kiswahili", []string{"Latin"}},
	{"af", "Afrikaans", "Afrikaans", []string{"Latin"}},
	{"ca", "Catalan", "Català", []string{"Latin"}},
	{"eu", "Basque", "Euskara", []string{"Latin"}},
	{"gl", "Galician", "Galego", []string{"Latin"}},
	{"hr", "Croatian", "Hrvatski", []string{"Latin"}},
	{"sr", "Serbian", "Српски", []string{"Cyrillic", "Latin"}},
	{"sl", "Slovenian", "Slovenščina", []string{"Latin"}},
	{"et", "Estonian", "Eesti", []string{"Latin"}},
	{"lv", "Latvian", "Latviešu", []string{"Latin"}},
	{"lt", "Lithuanian", "Lietuvių", []string{"Latin"}},
	{"is", "Icelandic", "Íslenska", []string{"Latin"}},
	{"ga", "Irish", "Gaeilge", []string{"Latin"}},
	{"cy", "Welsh", "Cymraeg", []string{"Latin"}},
	{"az", "Azerbaijani", "Azərbaycan", []string{"Latin"}},
	{"kk", "Kazakh", "Қазақ", []string{"Cyrillic"}},
	{"uz", "Uzbek", "Oʻzbek", []string{"Latin"}},
	{"mn", "Mongolian", "Монгол", []string{"Cyrillic"}},
}

// LanguageByCode returns language info by ISO code.
func LanguageByCode(code string) (Language, bool) {
	for _, lang := range SupportedLanguages {
		if lang.Code == code {
			return lang, true
		}
	}
	return Language{}, false
}

// ListLanguageCodes returns all supported language codes.
func ListLanguageCodes() []string {
	codes := make([]string, len(SupportedLanguages))
	for i, lang := range SupportedLanguages {
		codes[i] = lang.Code
	}
	return codes
}
