package common

// ExtractBrowser 从User-Agent字符串中提取浏览器标识。
func ExtractBrowser(userAgent string) string {
	if userAgent == "" {
		return "未知"
	}
	uaLower := ToLower(userAgent)
	switch {
	case ContainsStr(uaLower, "edge"):
		return "Edge"
	case ContainsStr(uaLower, "chrome"):
		return "Chrome"
	case ContainsStr(uaLower, "firefox"):
		return "Firefox"
	case ContainsStr(uaLower, "safari"):
		return "Safari"
	case ContainsStr(uaLower, "msie") || ContainsStr(uaLower, "trident"):
		return "IE"
	case ContainsStr(uaLower, "postman"):
		return "Postman"
	default:
		return "其他"
	}
}

// ExtractOS 从User-Agent字符串中提取操作系统标识。
func ExtractOS(userAgent string) string {
	if userAgent == "" {
		return "未知"
	}
	uaLower := ToLower(userAgent)
	switch {
	case ContainsStr(uaLower, "windows"):
		return "Windows"
	case ContainsStr(uaLower, "mac os"):
		return "MacOS"
	case ContainsStr(uaLower, "android"):
		return "Android"
	case ContainsStr(uaLower, "iphone") || ContainsStr(uaLower, "ipad"):
		return "iOS"
	case ContainsStr(uaLower, "linux"):
		return "Linux"
	default:
		return "其他"
	}
}

// todo: 获取ip所属地区
func GetRegion(ip string) string {
	return ""
}
