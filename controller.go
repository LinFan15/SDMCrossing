package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type TrafficSolution struct {
	TrafficGroups  []string
	ExcludedGroups []string
	Score          float64
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
	TrafficGroups   map[string]*TrafficGroup
	Mutex           *sync.Mutex
	CurrentSolution TrafficSolution
	InMotion        bool
	ElapsedTime     int64
}

func NewController() Controller {
	return Controller{
		trafficModel,
		&sync.Mutex{},
		TrafficSolution{},
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

func (c *Controller) SetTrafficLightState(trafficLight string, state int) {
	parsedTrafficLight := parseAbsoluteName(trafficLight)

	c.Mutex.Lock()
	c.TrafficGroups[parsedTrafficLight[0]].Lights[parsedTrafficLight[1]] = state
	c.Mutex.Unlock()
}

func (c *Controller) updateScores() {
	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		c.Mutex.Lock()
		c.TrafficGroups[trafficGroupName].BaseScore = 0
		c.Mutex.Unlock()

		if trafficGroup.Lights["/light/1"] > 0 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore = 0
			c.Mutex.Unlock()

			continue
		}

		i := 1

		for sensorName := range trafficGroup.Sensors {
			canonicalSensorName := getCanonicalSensorName(trafficGroupName, sensorName)

			if c.TrafficGroups[canonicalSensorName[0]].Sensors[canonicalSensorName[1]] {
				c.Mutex.Lock()
				c.TrafficGroups[trafficGroupName].BaseScore += 1 / float64(i)
				c.Mutex.Unlock()
			} else {
				break
			}

			i += 1
		}

		if trafficGroup.BaseScore >= 1 {
			c.Mutex.Lock()
			c.TrafficGroups[trafficGroupName].TimeScore += 0.01219 // 0.0104 0.004167
			c.Mutex.Unlock()
		}
	}
}

func (c Controller) GenerateSolutions() []TrafficSolution {
	var solutions []TrafficSolution

	for trafficGroupName, trafficGroup := range c.TrafficGroups {

		isValid := true

		for _, exceptionSensor := range trafficGroup.ExceptionSensors {
			parsedExceptionSensor := parseAbsoluteName(exceptionSensor)

			if c.TrafficGroups[parsedExceptionSensor[0]].Sensors[parsedExceptionSensor[1]] {
				isValid = false
				break
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
		solution = c.generateSolutionsRecurse(solution)
		// solution.Score = solution.Score / float64(len(solution.TrafficGroups))
		solutions = append(solutions, solution)
	}

	sort.Sort(sort.Reverse(ByScore(solutions)))

	return solutions
}

func (c Controller) generateSolutionsRecurse(currentSolution TrafficSolution) TrafficSolution {
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

			if !isValid {
				continue
			}

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

					messages = append(messages, mqttMessage{
						"5" + trafficGroupName + light,
						"1",
					})
				}
			}
		}

		if c.ElapsedTime >= 7000 {
			for _, trafficGroupName := range c.CurrentSolution.TrafficGroups {

				var state int

				if c.TrafficGroups[trafficGroupName].LightDirection == "RTG" {
					state = 0
				} else if c.TrafficGroups[trafficGroupName].LightDirection == "GTR" {
					state = 2
				}

				for light := range c.TrafficGroups[trafficGroupName].Lights {
					c.Mutex.Lock()
					c.TrafficGroups[trafficGroupName].Lights[light] = state
					c.Mutex.Unlock()

					messages = append(messages, mqttMessage{
						"5" + trafficGroupName + light,
						strconv.Itoa(state),
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
			//for name, trafficGroup := range c.TrafficGroups {
			//	fmt.Println(name + ":", trafficGroup.TimeScore, "-", trafficGroup.BaseScore, "|", trafficGroup.Sensors)
			//}

			for _, rangeSolution := range solutions {
				fmt.Println("Groups:", rangeSolution.TrafficGroups, "- Score:", rangeSolution.Score)
			}
			fmt.Println()

			for _, trafficGroupName := range solution.TrafficGroups {

				var state int

				if c.TrafficGroups[trafficGroupName].LightDirection == "RTG" {
					state = 2
				} else if c.TrafficGroups[trafficGroupName].LightDirection == "GTR" {
					state = 0
				}

				for light := range c.TrafficGroups[trafficGroupName].Lights {
					c.Mutex.Lock()
					c.TrafficGroups[trafficGroupName].Lights[light] = state
					c.Mutex.Unlock()

					messages = append(messages, mqttMessage{
						"5" + trafficGroupName + light,
						strconv.Itoa(state),
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
