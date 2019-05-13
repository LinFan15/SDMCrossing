package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type TrafficSolution struct {
	TrafficGroups  []string
	ExcludedGroups []string
	Score          float64
}

type Controller struct {
	TrafficGroups   map[string]*TrafficGroup
	Mutex           *sync.Mutex
	CurrentSolution TrafficSolution
	Pending         int
	InMotion        bool
	ElapsedTime     int64
}

func NewController() Controller {
	return Controller{
		trafficModel,
		&sync.Mutex{},
		TrafficSolution{},
		0,
		false,
		0,
	}
}

func (c *Controller) SetSensorState(sensorName string, state bool) {
	parsedSensorName := parseAbsoluteName(sensorName)

	c.Mutex.Lock()
	c.TrafficGroups[parsedSensorName[0]].Sensors[parsedSensorName[1]] = state
	c.Mutex.Unlock()
}

func (c *Controller) SetTrafficGroupState(trafficGroupName string, state int) []mqttMessage {
	var messages []mqttMessage

	for light := range c.TrafficGroups[trafficGroupName].Lights {
		c.Mutex.Lock()
		c.TrafficGroups[trafficGroupName].Lights[light] = state
		c.Mutex.Unlock()

		messages = append(messages, mqttMessage{
			"5" + trafficGroupName + light,
			strconv.Itoa(state),
		})
	}

	return messages
}

func (c *Controller) SetTrafficLightsInitialState() []mqttMessage {
	var messages []mqttMessage

	for trafficGroupName, trafficGroup := range c.TrafficGroups {

		var state int

		if trafficGroup.LightDirection == "RTG" {
			state = 0
		} else if trafficGroup.LightDirection == "GTR" {
			state = 2
		}

		messages = c.SetTrafficGroupState(trafficGroupName, state)
	}

	c.CurrentSolution = TrafficSolution{}
	c.ElapsedTime = 0
	c.InMotion = false

	return messages
}

func (c *Controller) handleTrafficGroup(trafficGroupName string, mc chan []mqttMessage) {
	var state int
	var messages []mqttMessage

	if c.TrafficGroups[trafficGroupName].LightDirection == "RTG" {
		state = 2
	} else if c.TrafficGroups[trafficGroupName].LightDirection == "GTR" {
		state = 0
	}

	messages = c.SetTrafficGroupState(trafficGroupName, state)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[0])

	state = 1
	messages = c.SetTrafficGroupState(trafficGroupName, state)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	if c.TrafficGroups[trafficGroupName].LightDirection == "RTG" {
		state = 0
	} else if c.TrafficGroups[trafficGroupName].LightDirection == "GTR" {
		state = 2
	}

	messages = c.SetTrafficGroupState(trafficGroupName, state)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[2])

	c.Mutex.Lock()
	c.Pending -= 1
	c.Mutex.Unlock()

	return
}

func (c *Controller) updateScores() {
	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		if len(trafficGroup.Lights) == 0 {
			continue
		}

		c.Mutex.Lock()
		c.TrafficGroups[trafficGroupName].BaseScore = 0
		c.Mutex.Unlock()

		if trafficGroup.LightDirection == "RTG" && trafficGroup.Lights["/light/1"] > 0 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore = 0
			c.Mutex.Unlock()

			continue
		} else if trafficGroup.LightDirection == "GTR" && trafficGroup.Lights["/light/1"] < 2 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore = 0
			c.Mutex.Unlock()

			continue
		}

		for sensorName := range trafficGroup.Sensors {
			canonicalSensorName := getCanonicalSensorName(trafficGroupName, sensorName)
			sensorNumber, _ := strconv.Atoi(canonicalSensorName[1][len(canonicalSensorName[1])-1:])

			if c.TrafficGroups[canonicalSensorName[0]].Sensors[canonicalSensorName[1]] {
				c.Mutex.Lock()
				c.TrafficGroups[trafficGroupName].BaseScore += 1 / float64(sensorNumber)
				c.Mutex.Unlock()
			}
		}

		if trafficGroup.BaseScore >= 1 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore += 0.01219 // 0.0104 0.004167
			c.Mutex.Unlock()
		}
	}
}

func (c Controller) GenerateSolution() TrafficSolution {
	var currentSolution TrafficSolution

	for trafficGroupName, trafficGroup := range c.TrafficGroups {

		isValid := true

		for _, exceptionSensor := range trafficGroup.ExceptionSensors {
			parsedExceptionSensor := parseAbsoluteName(exceptionSensor)

			if c.TrafficGroups[parsedExceptionSensor[0]].Sensors[parsedExceptionSensor[1]] {
				isValid = false
				break
			}
		}

		for _, sensor := range trafficGroup.Sensors {
			if sensor {
				isValid = true
				break
			} else {
				isValid = false
			}
		}

		if !isValid {
			continue
		}

		score := trafficGroup.BaseScore + trafficGroup.TimeScore

		solution := TrafficSolution{
			[]string{trafficGroupName},
			trafficGroup.ExcludedGroups,
			score,
		}

		fmt.Println("Group:", trafficGroupName, " - Score:", trafficGroup.BaseScore, trafficGroup.TimeScore)

		if solution.Score > currentSolution.Score {
			currentSolution = solution
		}
	}

	solution := c.generateSolutionsRecurse(currentSolution)
	// solution.Score = solution.Score / float64(len(solution.TrafficGroups))

	return solution
}

func (c Controller) generateSolutionsRecurse(currentSolution TrafficSolution) TrafficSolution {
	if currentSolution.Score == 0 {
		return currentSolution
	}

	var newSolution TrafficSolution
	currentTrafficGroupName := currentSolution.TrafficGroups[len(currentSolution.TrafficGroups)-1]

	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		if trafficGroupName != currentTrafficGroupName &&
			!stringInSlice(trafficGroupName, currentSolution.ExcludedGroups) &&
			!stringInSlice(trafficGroupName, currentSolution.TrafficGroups) {

			isValid := true

			for _, exceptionSensor := range trafficGroup.ExceptionSensors {
				parsedExceptionSensor := parseAbsoluteName(exceptionSensor)

				if c.TrafficGroups[parsedExceptionSensor[0]].Sensors[parsedExceptionSensor[1]] {
					isValid = false
					break
				}
			}

			for _, sensor := range trafficGroup.Sensors {
				if sensor {
					isValid = true
					break
				} else {
					isValid = false
				}
			}

			if !isValid {
				continue
			}

			solution := TrafficSolution{
				append(currentSolution.TrafficGroups, trafficGroupName),
				append(currentSolution.ExcludedGroups, trafficGroup.ExcludedGroups...),
				currentSolution.Score + trafficGroup.BaseScore + trafficGroup.TimeScore,
			}

			if solution.Score > newSolution.Score {
				newSolution = solution
			}
		}
	}

	if newSolution.Score == 0 {
		return currentSolution
	}

	newSolution = c.generateSolutionsRecurse(newSolution)

	return newSolution
}

func (c *Controller) process(mc chan []mqttMessage) {
	c.updateScores()

	if c.InMotion {
		c.ElapsedTime += 500

		if c.Pending == 0 {
			mc <- c.SetTrafficLightsInitialState()
		}

	} else {
		c.ElapsedTime = 0

		solution := c.GenerateSolution()

		if solution.Score > 0 {
			//for name, trafficGroup := range c.TrafficGroups {
			//	fmt.Println(name + ":", trafficGroup.TimeScore, "-", trafficGroup.BaseScore, "|", trafficGroup.Sensors)
			//}

			fmt.Println("Groups:", solution.TrafficGroups, "- Score:", solution.Score)
			fmt.Println()

			c.Pending = len(solution.TrafficGroups)

			for _, trafficGroupName := range solution.TrafficGroups {

				go c.handleTrafficGroup(trafficGroupName, mc)
			}

			c.CurrentSolution = solution
			c.InMotion = true
		} else {
			fmt.Println("Waiting for Simulator..")
		}

	}
}

func (c Controller) Loop(mc chan []mqttMessage) {
	for {
		time.Sleep(500 * time.Millisecond)
		c.process(mc)
	}
}
