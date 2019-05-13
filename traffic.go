package main

import (
	"strings"
	"time"
)

type TrafficModel map[string]*TrafficGroup

type TrafficGroup struct {
	Lights           map[string]int
	Sensors          map[string]bool
	LightDirection   string
	Duration         []time.Duration
	ExcludedGroups   []string
	ExceptionSensors []string
	BaseScore        float64
	TimeScore        float64
}

var trafficModel TrafficModel = map[string]*TrafficGroup{
	"/motor_vehicle/1": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/5",
			"/motor_vehicle/8",
			"/cycle/1",
			"/cycle/4",
			"/foot/1",
			"/foot/8",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/2": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/5",
			"/motor_vehicle/6",
			"/motor_vehicle/8",
			"/motor_vehicle/9",
			"/motor_vehicle/10",
			"/motor_vehicle/11",
			"/motor_vehicle/12",
			"/cycle/1",
			"/cycle/3",
			"/foot/1",
			"/foot/6",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/3": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/5",
			"/motor_vehicle/6",
			"/motor_vehicle/7",
			"/motor_vehicle/8",
			"/motor_vehicle/10",
			"/motor_vehicle/11",
			"/cycle/1",
			"/cycle/2",
			"/foot/1",
			"/foot/4",
		},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/4": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/8",
			"/motor_vehicle/11",
			"/cycle/1",
			"/cycle/2",
			"/foot/2",
			"/foot/3",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/5": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
			"/sensor/5": false,
			"/sensor/6": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/8",
			"/motor_vehicle/11",
			"/cycle/2",
			"/cycle/4",
			"/foot/3",
			"/foot/8",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/6": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/8",
			"/motor_vehicle/9",
			"/motor_vehicle/10",
			"/motor_vehicle/12",
			"/cycle/2",
			"/cycle/3",
			"/foot/3",
			"/foot/6",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/7": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/3",
			"/motor_vehicle/10",
			"/motor_vehicle/12",
			"/cycle/2",
			"/cycle/3",
			"/foot/4",
			"/foot/5",
		},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/8": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/4",
			"/motor_vehicle/5",
			"/motor_vehicle/6",
			"/motor_vehicle/10",
			"/motor_vehicle/11",
			"/motor_vehicle/12",
			"/cycle/1",
			"/cycle/3",
			"/cycle/4",
			"/foot/2",
			"/foot/5",
			"/foot/8",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/9": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/6",
			"/motor_vehicle/12",
			"/cycle/3",
			"/cycle/4",
			"/foot/6",
			"/foot/7",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/10": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
			"/sensor/5": false,
			"/sensor/6": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/6",
			"/motor_vehicle/7",
			"/motor_vehicle/8",
			"/cycle/2",
			"/cycle/4",
			"/foot/4",
			"/foot/7",
		},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/11": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/4",
			"/motor_vehicle/5",
			"/motor_vehicle/8",
			"/cycle/1",
			"/cycle/4",
			"/foot/2",
			"/foot/7",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/12": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/7",
			"/motor_vehicle/8",
			"/motor_vehicle/9",
			"/motor_vehicle/2",
			"/motor_vehicle/6",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/13": {
		map[string]int{
			"/light/1": 2,
			"/light/2": 2,
		},
		map[string]bool{
			"/vessel/1/sensor/1": false,
			"/vessel/2/sensor/1": false,
		},
		"GTR",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/vessel/1",
			"/vessel/2",
		},
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/14": {
		map[string]int{},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			7 * time.Second,
			4 * time.Second,
			1 * time.Second,
		},
		[]string{},
		[]string{},
		0,
		0,
	},
	"/cycle/1": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			8 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/2",
			"/motor_vehicle/3",
			"/motor_vehicle/4",
			"/motor_vehicle/8",
			"/motor_vehicle/11",
		},
		[]string{},
		0,
		0,
	},
	"/cycle/2": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			8 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/3",
			"/motor_vehicle/4",
			"/motor_vehicle/5",
			"/motor_vehicle/6",
			"/motor_vehicle/7",
			"/motor_vehicle/10",
		},
		[]string{},
		0,
		0,
	},
	"/cycle/3": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			8 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/6",
			"/motor_vehicle/7",
			"/motor_vehicle/8",
			"/motor_vehicle/9",
		},
		[]string{},
		0,
		0,
	},
	"/cycle/4": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			8 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/5",
			"/motor_vehicle/8",
			"/motor_vehicle/9",
			"/motor_vehicle/10",
			"/motor_vehicle/11",
		},
		[]string{},
		0,
		0,
	},
	"/cycle/5": {
		map[string]int{
			"/light/1": 2,
			"/light/2": 2,
		},
		map[string]bool{
			"/vessel/1/sensor/1": false,
			"/vessel/2/sensor/1": false,
		},
		"GTR",
		[]time.Duration{
			8 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/vessel/1",
			"/vessel/2",
		},
		[]string{},
		0,
		0,
	},
	"/foot/1": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/2",
			"/motor_vehicle/3",
		},
		[]string{},
		0,
		0,
	},
	"/foot/2": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/4",
			"/motor_vehicle/8",
			"/motor_vehicle/11",
		},
		[]string{},
		0,
		0,
	},
	"/foot/3": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/4",
			"/motor_vehicle/5",
			"/motor_vehicle/6",
		},
		[]string{},
		0,
		0,
	},
	"/foot/4": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG", []time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},

		[]string{
			"/motor_vehicle/3",
			"/motor_vehicle/7",
			"/motor_vehicle/10",
		},
		[]string{},
		0,
		0,
	},
	"/foot/5": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/7",
			"/motor_vehicle/8",
		},
		[]string{},
		0,
		0,
	},
	"/foot/6": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/2",
			"/motor_vehicle/6",
			"/motor_vehicle/9",
		},
		[]string{},
		0,
		0,
	},
	"/foot/7": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/9",
			"/motor_vehicle/10",
			"/motor_vehicle/11",
		},
		[]string{},
		0,
		0,
	},
	"/foot/8": {
		map[string]int{
			"/light/1": 0,
			"/light/2": 0,
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		"RTG",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/1",
			"/motor_vehicle/5",
			"/motor_vehicle/8",
		},
		[]string{},
		0,
		0,
	},
	"/foot/9": {
		map[string]int{
			"/light/1": 2,
			"/light/2": 2,
		},
		map[string]bool{
			"/vessel/1/sensor/1": false,
			"/vessel/2/sensor/1": false,
		},
		"GTR",
		[]time.Duration{
			6 * time.Second,
			6 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/vessel/1",
			"/vessel/2",
		},
		[]string{},
		0,
		0,
	},
	"/vessel/1": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			12 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/13",
			"/cycle/5",
			"/foot/9",
		},
		[]string{},
		0,
		0,
	},
	"/vessel/2": {
		map[string]int{
			"/light/1": 0,
		},
		map[string]bool{
			"/sensor/1": false,
		},
		"RTG",
		[]time.Duration{
			12 * time.Second,
			2 * time.Second,
			1 * time.Second,
		},
		[]string{
			"/motor_vehicle/13",
			"/cycle/5",
			"/foot/9",
		},
		[]string{},
		0,
		0,
	},
}

func parseAbsoluteName(s string) [2]string {
	splitEntityName := strings.SplitN(s, "/", 4)
	splitTrafficGroup := splitEntityName[:len(splitEntityName)-1]

	trafficGroup := strings.Join(splitTrafficGroup, "/")
	relEntityName := "/" + splitEntityName[len(splitEntityName)-1]

	return [2]string{
		trafficGroup,
		relEntityName,
	}
}

func getCanonicalSensorName(trafficGroup string, sensorName string) [2]string {
	splitSensorName := strings.Split(sensorName, "/")

	if len(splitSensorName) > 3 {
		return parseAbsoluteName(sensorName)
	}

	return [2]string{
		trafficGroup,
		sensorName,
	}
}
