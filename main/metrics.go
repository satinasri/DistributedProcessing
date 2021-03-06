package main

import (
	"time"
	"strconv"
	"math/big"
)

type metric struct {
	Perf int
	IsPrime bool
	hPerf time.Duration
	Val string
}

//Maps node name/ID to Reputation
type RepMetrics struct {
	CurrentMetrics map[string]Reputation
}

type Reputation struct {
	Score		int
	Count		int
	Correct 	int
}


//Map that maps node names to the first time that node was seen
type Uptimes struct{
	Uptimes map[string]Uptime
}

type Uptime struct{
	time time.Time
}

func newUptimes(name string)Uptimes{
	var ret Uptimes
	ret.Uptimes = make(map[string]Uptime)
	var tmp Uptime
	tmp.time = time.Now()
	ret.Uptimes[name] = tmp
	return ret
}

func newRepMetrics(name string)RepMetrics{
	var ret RepMetrics
	ret.CurrentMetrics = make(map[string]Reputation)
	var tmp Reputation
	tmp.Correct = 0
	tmp.Count = 0
	tmp.Score = 0
	ret.CurrentMetrics[name] = tmp
	return ret
}

//If a node doesn't currently have an uptime, add one
func updateUptime(ut Uptimes, name string, newtime Uptime) bool{
	_, ok := ut.Uptimes[name]
	if !ok{
		ut.Uptimes[name] = newtime
		return true
	}
	return false
}

//Remove the uptime associated with a node
func clearUptime(ut Uptimes, name string){
	delete(ut.Uptimes, name)
}

func getLongestUptime(ut Uptimes) (string,Uptime){
	longest := time.Now()
	lname := ""
	for k,v := range ut.Uptimes{
		if longest.After(v.time){
			longest = v.time
			lname = k
		}
	}
	return lname, ut.Uptimes[lname]
}

//Scorer should take in the current reputation and the new result and update the reputation as a result
func updateReputation(repmets RepMetrics, newmet metric, node string, scorer func(nm metric, rp Reputation)) bool{
	rep, ok := repmets.CurrentMetrics[node]
	if !ok{
		return false
	}
	scorer(newmet, rep)
	return true
}

//The score for hashing is the average time it takes to generate a collision
//It doesn't use correctness currently
func hashScorer(met metric, rep Reputation){
	rep.Count += 1
	newscore := rep.Score / rep.Count
	newscore += int(met.hPerf)
	newscore = newscore/ rep.Count
	rep.Score = newscore
}

//the score for primality is the average number of correct assessments out of 100,000
//The score is score = correct/count
func primeScorer(met metric, rep Reputation){
	i, _ := strconv.ParseInt(met.Val,10,64)
	test := big.NewInt(i)
	if met.IsPrime == testPrime(*test).IsPrime{
		rep.Correct += 1
	}
	rep.Count += 1
	rep.Score = rep.Correct / rep.Count
}