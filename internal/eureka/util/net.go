package util

import (
	"net/url"
	"path"
)

func JoinPathStr(host, p string) (string, error) {
	if u, err := url.Parse(host); err != nil {
		return "", err
	} else {
		u.Path = path.Join(u.Path, p)
		return u.String(), nil
	}
}
