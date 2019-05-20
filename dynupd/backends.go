package main

import (
	"strings"
)

func evaluateBackendRules(rules [][]string) (m map[string]string) {
	m = map[string]string{}
	for _, rule := range rules {
		if strings.ToLower(rule[0]) == "manual" {
			g := rule[3]
			o := m[g]
			vs := strings.Split(rule[1], ",")
			for _, v := range vs {
				if len(o) > 0 {
					o = o + ","
				}
				o = o + v + ":" + rule[2]
			}
			m[g] = o
		}
	}
	return
}
