package main

import log "github.com/sirupsen/logrus"

const p95cap = 1000

func getP95(p95 []int) int {
	n := len(p95)
	if n == 0 {
		return 0
	}
	return p95[n-1]
}

func getP95Max(p95 []int) int {
	n := len(p95)
	if n == 0 {
		return 0
	}
	return p95[0]
}

func calP95(p95 []int, count int, rt int) []int {
	l := len(p95)
	if l == 0 {
		p95 = append(p95, rt)
		return p95
	}
	n := count*5/100 + 1 // number of element to store
	if n == 0 {
		n = 1
	}
	if n > p95cap {
		n = p95cap
	}
	if n > l {
		if n > l+1 {
			log.Errorf(" n>l+1 somthing worng in p95  l:%v,n:%v,count:%v\n", l, n, count)
		}
		p95 = append(p95, p95[l-1])
	}

	i := n - 1
	for ; i >= 0; i-- {
		if rt < p95[i] {
			break
		}
	}
	i++ //rt should replease next index
	if i == n {
		return p95
	}
	for j := n - 1; j > i; j-- { // shift value in array
		p95[j] = p95[j-1]
	}
	p95[i] = rt

	return p95
}
