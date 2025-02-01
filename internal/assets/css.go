package assets

import _ "embed"

//go:embed style.css
var css string

func GetCss() string {
	return css
}
