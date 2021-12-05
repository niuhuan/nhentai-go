package nhentai


const TagLanguageJapanese = 6346
const TagLanguageEnglish = 12227
const TagLanguageChinese = 29963

// lang **使用标签ids确定语言
func lang(ids []int) string {
	if contains(ids, TagLanguageJapanese) {
		return "JP"
	} else if contains(ids, TagLanguageEnglish) {
		return "EN"
	} else if contains(ids, TagLanguageChinese) {
		return "CH"
	}
	return ""
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
