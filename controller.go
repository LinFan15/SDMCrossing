package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TrafficGroup struct {
	Lights map[string]int
	Sensors map[string]bool
	ExcludedGroups []string
	BaseScore float64
	TimeScore float64
}


type TrafficSolution struct {
	TrafficGroups []string
	ExcludedGroups []string
	Score float64
}


type ByScore []TrafficSolution

func (s ByScore) Len() int {
	return len(s)
}

func (s ByScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByScore) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}


type Controller struct {
	TrafficGroups map[string]*TrafficGroup
	Mutex *sync.Mutex
	CurrentSolution TrafficSolution
	InMotion bool
	ElapsedTime int64
}

func NewController() Controller{
	return Controller {
		map[string]*TrafficGroup {
			"/motor_vehicle/1": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
				},
				[]string {
					"/motor_vehicle/5",
					"/motor_vehicle/8",
				},
				0,
				0,
			},
			"/motor_vehicle/2": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
				},
				[]string {
					"/motor_vehicle/5",
					"/motor_vehicle/6",
					"/motor_vehicle/8",
					"/motor_vehicle/9",
					"/motor_vehicle/10",
					"/motor_vehicle/11",
				},
				0,
				0,
			},
			"/motor_vehicle/3": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
				},
				[]string {
					"/motor_vehicle/5",
					"/motor_vehicle/6",
					"/motor_vehicle/7",
					"/motor_vehicle/8",
					"/motor_vehicle/10",
					"/motor_vehicle/11",
				},
				0,
				0,
			},
			"/motor_vehicle/4": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
				},
				[]string {
					"/motor_vehicle/8",
					"/motor_vehicle/11",
				},
				0,
				0,
			},
			"/motor_vehicle/5": {
				map[string]int {
					"/light/1": 0,
					"/light/2": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
					"/sensor/4": false,
					"/sensor/5": false,
					"/sensor/6": false,
				},
				[]string {
					"/motor_vehicle/1",
					"/motor_vehicle/2",
					"/motor_vehicle/3",
					"/motor_vehicle/8",
					"/motor_vehicle/11",
				},
				0,
				0,
			},
			"/motor_vehicle/6": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
				},
				[]string {
					"/motor_vehicle/2",
					"/motor_vehicle/3",
					"/motor_vehicle/8",
					"/motor_vehicle/9",
					"/motor_vehicle/10",
				},
				0,
				0,
			},
			"/motor_vehicle/7": {
				map[string]int {
					"/light/1": 0,
					"/light/2": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
					"/sensor/4": false,
				},
				[]string {
					"/motor_vehicle/3",
					"/motor_vehicle/10",
				},
				0,
				0,
			},
			"/motor_vehicle/8": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
				},
				[]string {
					"/motor_vehicle/1",
					"/motor_vehicle/2",
					"/motor_vehicle/3",
					"/motor_vehicle/4",
					"/motor_vehicle/5",
					"/motor_vehicle/6",
					"/motor_vehicle/10",
					"/motor_vehicle/11",
				},
				0,
				0,
			},
			"/motor_vehicle/9": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
				},
				[]string {
					"/motor_vehicle/2",
					"/motor_vehicle/6",
				},
				0,
				0,
			},
			"/motor_vehicle/10": {
				map[string]int {
					"/light/1": 0,
					"/light/2": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
					"/sensor/4": false,
					"/sensor/5": false,
					"/sensor/6": false,
				},
				[]string {
					"/motor_vehicle/2",
					"/motor_vehicle/3",
					"/motor_vehicle/6",
					"/motor_vehicle/7",
					"/motor_vehicle/8",
				},
				0,
				0,
			},
			"/motor_vehicle/11": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
					"/sensor/2": false,
					"/sensor/3": false,
				},
				[]string {
					"/motor_vehicle/2",
					"/motor_vehicle/3",
					"/motor_vehicle/4",
					"/motor_vehicle/5",
					"/motor_vehicle/8",
				},
				0,
				0,
			},
			"/cycle/1": {
				map[string]int {
					"/light/1": 0,
				},
				map[string]bool {
					"/sensor/1": false,
				},
				[]string{},
				0,
				0,
			},
		},
		&sync.Mutex{},
		TrafficSolution{},
		false,
		0,
	}
}

func (c *Controller) SetSensorState(sensorName string, state bool) {
	splitSensorName := strings.SplitN(sensorName, "/", 4)
	splitTrafficGroup := splitSensorName[:len(splitSensorName) - 1]

	trafficGroup := strings.Join(splitTrafficGroup, "/")
	relSensorName := "/" + splitSensorName[len(splitSensorName) - 1]

	c.Mutex.Lock()
	c.TrafficGroups[trafficGroup].Sensors[relSensorName] = state
	c.Mutex.Unlock()
}

func (c *Controller) SetTrafficLightState(trafficLight string, state int) {
	splitTrafficLight := strings.SplitN(trafficLight, "/", 4)
	splitTrafficGroup := splitTrafficLight[:len(splitTrafficLight) - 1]

	trafficGroup := strings.Join(splitTrafficGroup, "/")
	relTrafficLight := "/" +  splitTrafficLight[len(splitTrafficLight) - 1]

	c.Mutex.Lock()
	c.TrafficGroups[trafficGroup].Lights[relTrafficLight] = state
	c.Mutex.Unlock()
}

func (c *Controller) updateScores() {
	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		c.Mutex.Lock()
		c.TrafficGroups[trafficGroupName].BaseScore = 0
		c.Mutex.Unlock()

		if trafficGroup.Lights["/lights/1"] > 0 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore = 0
			c.Mutex.Unlock()

			continue
		}

		for i := 0; i < len(trafficGroup.Sensors); i++ {
			if trafficGroup.Sensors["/sensor/" + strconv.Itoa(i+1)] {
				c.Mutex.Lock()
				c.TrafficGroups[trafficGroupName].BaseScore += 1 / float64(i + 1)
				c.Mutex.Unlock()
			} else {
				break
			}
		}

		if trafficGroup.BaseScore >= 1 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore += 0.003472
			c.Mutex.Unlock()
		}
	}
}

func (c Controller) GenerateSolutions() []TrafficSolution {
	var solutions []TrafficSolution

	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		score := trafficGroup.BaseScore + trafficGroup.TimeScore
		solution := TrafficSolution {
			[]string { trafficGroupName },
			trafficGroup.ExcludedGroups,
			score,
		}

		solutions = append(solutions, c.generateSolutionsRecurse(solution))
	}

	sort.Sort(sort.Reverse(ByScore(solutions)))

	return solutions
}

func (c Controller) generateSolutionsRecurse(currentSolution TrafficSolution) TrafficSolution {
	currentTrafficGroupName := currentSolution.TrafficGroups[len(currentSolution.TrafficGroups) - 1]

	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		if trafficGroupName != currentTrafficGroupName &&
			!stringInSlice(trafficGroupName, currentSolution.ExcludedGroups) &&
			!stringInSlice(trafficGroupName, currentSolution.TrafficGroups) {

			currentSolution.TrafficGroups = append(currentSolution.TrafficGroups, trafficGroupName)
			currentSolution.ExcludedGroups = append(currentSolution.ExcludedGroups, trafficGroup.ExcludedGroups...)
			currentSolution.Score += trafficGroup.BaseScore + trafficGroup.TimeScore

			currentSolution = c.generateSolutionsRecurse(currentSolution)
		}
	}

	return currentSolution
}

func (c *Controller) process() []mqttMessage {
	var messages []mqttMessage

	c.updateScores()

	if c.InMotion {
		c.ElapsedTime += 500

		if c.ElapsedTime >= 5000 && c.TrafficGroups[c.CurrentSolution.TrafficGroups[0]].Lights["/light/1"] == 2 {
			for _, trafficGroupName := range c.CurrentSolution.TrafficGroups {
				for light := range c.TrafficGroups[trafficGroupName].Lights {
					c.Mutex.Lock()
					c.TrafficGroups[trafficGroupName].Lights[light] = 1
					c.Mutex.Unlock()

					messages = append(messages, mqttMessage {
						"5" + trafficGroupName + light,
						"1",
					})
				}
			}
		}

		if c.ElapsedTime >= 7000 {
			for _, trafficGroupName := range c.CurrentSolution.TrafficGroups {
				for light := range c.TrafficGroups[trafficGroupName].Lights {
					c.Mutex.Lock()
					c.TrafficGroups[trafficGroupName].Lights[light] = 0
					c.Mutex.Unlock()

					messages = append(messages, mqttMessage {
						"5" + trafficGroupName + light,
						"0",
					})
				}
			}

			c.CurrentSolution = TrafficSolution{}
			c.ElapsedTime = 0
			c.InMotion = false
		}

	} else {
		c.ElapsedTime = 0

		solutions := c.GenerateSolutions()
		solution := solutions[0]

		if solution.Score > 0 {
			for _, rangeSolution := range(solutions) {
				fmt.Println("Groups:", rangeSolution.TrafficGroups, "- Score:", rangeSolution.Score)
			}
			fmt.Println()

			for _, trafficGroupName := range solution.TrafficGroups {
				for light := range c.TrafficGroups[trafficGroupName].Lights {
					c.Mutex.Lock()
					c.TrafficGroups[trafficGroupName].Lights[light] = 2
					c.Mutex.Unlock()

					messages = append(messages, mqttMessage {
						"5" + trafficGroupName + light,
						"2",
					})
				}
			}

			c.CurrentSolution = solution
			c.InMotion = true
		} else {
			fmt.Println("Waiting for Simulator..")
		}

	}

	return messages
}

func (c Controller) Loop(mc chan []mqttMessage) {
	for {
		time.Sleep(500 * time.Millisecond)
		mc <- c.process()
	}
}