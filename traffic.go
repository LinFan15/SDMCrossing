package main

import (
	"strings"
	"time"
)

type TrafficModel map[string]*TrafficGroup

type TrafficGroup struct {
	Items            map[string]*TrafficItem
	Sensors          map[string]bool
	Duration         []time.Duration
	ExcludedGroups   []string
	AssociatedGroups []string
	ExceptionSensors []string
	BaseScore        float64
	TimeScore        float64
}

type TrafficItem struct {
	State int
	Type  string
}

var trafficModel TrafficModel = map[string]*TrafficGroup{
	"/motor_vehicle/1": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/2": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/3": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/4": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/5": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
			"/sensor/5": false,
			"/sensor/6": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/6": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/7": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
		},
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
		[]string{},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/8": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/9": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/10": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
			"/sensor/4": false,
			"/sensor/5": false,
			"/sensor/6": false,
		},
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
		[]string{},
		[]string{
			"/motor_vehicle/14/sensor/1",
		},
		0,
		0,
	},
	"/motor_vehicle/11": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
			"/sensor/3": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/12": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/motor_vehicle/14": {
		map[string]*TrafficItem{},
		map[string]bool{
			"/sensor/1": false,
		},
		[]time.Duration{},
		[]string{},
		[]string{},
		[]string{},
		0,
		0,
	},
	"/cycle/1": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/cycle/2": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/cycle/3": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/cycle/4": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/1": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/2": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/3": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/4": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
		[]time.Duration{
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
		[]string{},
		0,
		0,
	},
	"/foot/5": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/6": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/7": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/foot/8": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
			"/light/2": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
			"/sensor/2": false,
		},
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
		[]string{},
		0,
		0,
	},
	"/vessel/1": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
		[]time.Duration{
			10 * time.Second,
			0 * time.Second,
			0 * time.Second,
		},
		[]string{
			"/motor_vehicle/13",
			"/cycle/5",
			"/foot/9",
		},
		[]string{
			"/bridge/1",
		},
		[]string{},
		0,
		0,
	},
	"/vessel/2": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"RTG",
			},
		},
		map[string]bool{
			"/sensor/1": false,
		},
		[]time.Duration{
			10 * time.Second,
			0 * time.Second,
			0 * time.Second,
		},
		[]string{
			"/motor_vehicle/13",
			"/cycle/5",
			"/foot/9",
		},
		[]string{
			"/bridge/1",
		},
		[]string{},
		0,
		0,
	},
	"/vessel/3": {
		map[string]*TrafficItem{},
		map[string]bool{
			"/sensor/1": false,
		},
		[]time.Duration{},
		[]string{},
		[]string{},
		[]string{},
		0,
		0,
	},
	"/bridge/1": {
		map[string]*TrafficItem{
			"/light/1": &TrafficItem{
				0,
				"GR",
			},
			"/gate/1": &TrafficItem{
				0,
				"GR",
			},
			"/gate/2": &TrafficItem{
				0,
				"GR",
			},
			"/deck/1": &TrafficItem{
				0,
				"RG",
			},
		},
		map[string]bool{
			"/sensor/1":          false,
			"/vessel/3/sensor/1": false,
		},
		[]time.Duration{
			6 * time.Second,
			4 * time.Second,
			10 * time.Second,
		},
		[]string{},
		[]string{},
		[]string{},
		0,
		0,
	},
}

func getTrafficItemFromGroup(groupName string) TrafficItem {
	for _, trafficItem := range trafficModel[groupName].Items {
		return *trafficItem
	}

	return TrafficItem{}
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

func getCanonicalName(trafficGroup string, itemName string) [2]string {
	splitItemName := strings.Split(itemName, "/")

	if len(splitItemName) > 3 {
		return parseAbsoluteName(itemName)
	}

	return [2]string{
		trafficGroup,
		itemName,
	}
}
