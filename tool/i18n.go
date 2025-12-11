package tool

var (
	i18n_en map[string]string = map[string]string{
		"time":      "Time: ",
		"prev":      "Previous",
		"back_menu": "Back to menu",
		"next":      "Next",
	}
	i18n_es map[string]string = map[string]string{
		"time":      "Tiempo: ",
		"prev":      "Anterior",
		"back_menu": "Volver al menú",
		"next":      "Siguiente",
	}
	i18n_zh map[string]string = map[string]string{
		"time":      "时间：",
		"prev":      "上一页",
		"back_menu": "返回菜单",
		"next":      "下一页",
	}
	i18n_zh_tw map[string]string = map[string]string{
		"time":      "時間：",
		"prev":      "上一頁",
		"back_menu": "返回菜單",
		"next":      "下一頁",
	}
)

func GetI18n(languageID int, key string) string {
	switch languageID {
	case 1:
		return i18n_en[key]
	case 2:
		return i18n_zh[key]
	case 3:
		return i18n_zh_tw[key]
	case 4:
		return i18n_es[key]
	}
	return key
}

func GetLanguageCode(languageID int) string {
	switch languageID {
	case 1:
		return "en"
	case 2:
		return "zh"
	case 3:
		return "zh-tw"
	case 4:
		return "es"
	}
	return "en"
}
