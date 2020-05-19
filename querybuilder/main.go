package querybuilder

import (
	"net/url"
	"strings"
)

type Query struct {
	query url.Values
}

func (query Query) URLSafe() string {
	return strings.Replace(query.query.Encode(), "+", "%20", -1)
}

type leftQuery struct {
	leftPart string
}

func (query Query) Parameter(parameter string) {
	query.query[parameter] = []string{}
}

func (query Query) Is(parameter string, value ...string) {
	query.query[parameter] = value
}
