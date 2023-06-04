package scraper

import "strings"

func rootDomain(url string) string {
	if strings.HasPrefix(url, "http") {
		return strings.Join(strings.Split(url, "/")[:3], "/")
	}
	return ""
}

func urlDir(url string) string {
	splitURL := strings.Split(strings.TrimSuffix(url, "/"), "/")
	if len(splitURL) <= 3 {
		return url
	}
	dir := strings.Join(splitURL[:len(splitURL)-1], "/")
	return dir
}

func validLink(url string) bool {
	prefixes := []string{"mailto", "rsync", "ftp", "javascript", "#"}
	for _, p := range prefixes {
		if strings.HasPrefix(url, p) {
			return false
		}
	}
	return true
}

func trimAfterHash(v string) string {
	if i := strings.Index(v, "#"); i != -1 {
		v = v[:i]
	}
	return v
}

func prefixRoot(parent, v string) string {
	if strings.HasPrefix(v, "/") {
		v = rootDomain(parent) + v
	} else if !strings.HasPrefix(v, "http") {
		v = urlDir(parent) + "/" + strings.TrimLeft(v, "./")
	}
	return v
}

func buildLink(parent, link string) string {
	pipeline := []func(string) string{
		trimAfterHash,
		func(v string) string {
			return prefixRoot(parent, v)
		},
	}

	link = pipe(pipeline...)(link)
	return link
}

func pipe(ops ...func(string) string) func(string) string {
	return func(v string) string {
		for _, op := range ops {
			v = op(v)
		}
		return v
	}
}
