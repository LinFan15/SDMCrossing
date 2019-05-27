package main

import (
	"fmt"
	"strconv"
	"time"
)

type TrafficSolution struct {
	TrafficGroups  []string
	ExcludedGroups []string
	Score          float64
}

type Controller struct {
	TrafficGroups     map[string]*TrafficGroup
	ChangesChannel    chan TrafficChange
	CurrentSolution   TrafficSolution
	Pending           int
	AssociatedPending []string
	InMotion          bool
	ElapsedTime       int64
}

func NewController(ch chan TrafficChange) Controller {
	return Controller{
		trafficModel,
		ch,
		TrafficSolution{},
		0,
		[]string{},
		false,
		0,
	}
}

func (c *Controller) SetSensorState(sensorName string, state bool) {
	parsedSensorName := parseAbsoluteName(sensorName)

	newState := 0

	if state {
		newState = 1
	}

	c.ChangesChannel <- TrafficChange{
		parsedSensorName,
		"sensor",
		float64(newState),
	}
}

func (c *Controller) SetTrafficItemState(trafficGroup string, trafficItem string, state int) []mqttMessage {
	var messages []mqttMessage

	c.ChangesChannel <- TrafficChange{
		[2]string{trafficGroup, trafficItem},
		"item",
		float64(state),
	}

	messages = append(messages, mqttMessage{
		"5" + trafficGroup + trafficItem,
		strconv.Itoa(state),
	})

	return messages
}

func (c *Controller) SetTrafficGroupState(trafficGroup string, setState int) []mqttMessage {
	var messages []mqttMessage

	for trafficItemName := range c.TrafficGroups[trafficGroup].Items {
		state := -1

		itemType := c.TrafficGroups[trafficGroup].Items[trafficItemName].Type

		if setState == 0 {
			if itemType == "RTG" || itemType == "GR" {
				state = 0
			} else if itemType == "GTR" {
				state = 2
			} else if itemType == "RG" {
				state = 1
			}
		} else if setState == 1 {
			if itemType == "RTG" || itemType == "GTR" {
				state = 1
			} else {
				state = c.TrafficGroups[trafficGroup].Items[trafficItemName].State
			}
		} else if setState == 2 {
			if itemType == "RTG" {
				state = 2
			} else if itemType == "GTR" || itemType == "RG" {
				state = 0
			} else if itemType == "GR" {
				state = 1
			}
		}

		if state != -1 {
			messages = append(messages, c.SetTrafficItemState(trafficGroup, trafficItemName, state)...)
		}
	}

	return messages
}

func (c *Controller) SetVesselGroupState(trafficGroupName string, setState int) []mqttMessage {
	var messages []mqttMessage

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		state := -1

		trafficItem := c.TrafficGroups[trafficGroupName].Items[trafficItemName]

		if setState == 0 {
			if trafficItem.Type == "RTG" && trafficItem.State > 0 {
				state = 0
			} else if trafficItem.Type == "GTR" && trafficItem.State < 2 {
				state = 2
			} else if trafficItem.Type == "RG" && trafficItem.State == 1 {
				state = 0
			} else if trafficItem.Type == "GR" && trafficItem.State == 0 {
				state = 1
			}
		} else if setState == 1 {
			if trafficItem.Type == "RTG" || trafficItem.Type == "GTR" {
				state = 1
			} else {
				state = c.TrafficGroups[trafficGroupName].Items[trafficItemName].State
			}
		} else if setState == 2 {
			if trafficItem.Type == "RTG" && trafficItem.State < 2 {
				state = 2
			} else if trafficItem.Type == "GTR" && trafficItem.State > 0 {
				state = 0
			} else if trafficItem.Type == "RG" && trafficItem.State == 0 {
				state = 1
			} else if trafficItem.Type == "GR" && trafficItem.State == 1 {
				state = 0
			}
		}
		if state != -1 {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, state)...)
		}
	}

	return messages
}

func (c *Controller) SetTrafficItemsInitialState() []mqttMessage {
	var messages []mqttMessage

	for trafficGroupName := range c.TrafficGroups {
		messages = append(messages, c.SetTrafficGroupState(trafficGroupName, 0)...)
	}

	c.CurrentSolution = TrafficSolution{}
	c.ElapsedTime = 0
	c.InMotion = false

	return messages
}

func (c *Controller) checkTrafficGroupSensors(trafficGroupName string) bool {
	for _, sensorState := range c.TrafficGroups[trafficGroupName].Sensors {
		if sensorState {
			return true
		}
	}

	return false
}

func (c *Controller) handleTrafficGroup(trafficGroupName string, mc chan []mqttMessage) {
	var messages []mqttMessage

	if trafficGroupName[:7] == "/bridge" {
		c.handleBridge(trafficGroupName, mc)
		if stringInSlice(trafficGroupName, c.AssociatedPending) {
			c.AssociatedPending = delFromSlice(indexOfStringInSlice(trafficGroupName, c.AssociatedPending), c.AssociatedPending)
		}
		return
	}

	if len(c.TrafficGroups[trafficGroupName].AssociatedGroups) > 0 {
		for _, associatedGroupName := range c.TrafficGroups[trafficGroupName].AssociatedGroups {
			if !stringInSlice(associatedGroupName, c.AssociatedPending) {
				c.AssociatedPending = append(c.AssociatedPending, associatedGroupName)

				go c.handleAssociatedTrafficGroup(associatedGroupName, mc)

				if associatedGroupName[:7] == "/bridge" {
					c.Pending -= 1
					return
				}
			} else if associatedGroupName[:7] == "/bridge" {
				c.Pending -= 1
				return
			}
		}
	}

	messages = c.SetTrafficGroupState(trafficGroupName, 2)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[0])

	messages = c.SetTrafficGroupState(trafficGroupName, 1)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	messages = c.SetTrafficGroupState(trafficGroupName, 0)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[2])

	c.Pending -= 1
	return
}

func (c *Controller) handleAssociatedTrafficGroup(trafficGroupName string, mc chan []mqttMessage) {
	var messages []mqttMessage

	if trafficGroupName[:7] == "/bridge" {
		c.handleBridge(trafficGroupName, mc)
		if stringInSlice(trafficGroupName, c.AssociatedPending) {
			c.AssociatedPending = delFromSlice(indexOfStringInSlice(trafficGroupName, c.AssociatedPending), c.AssociatedPending)
		}
		return
	}

	if len(c.TrafficGroups[trafficGroupName].AssociatedGroups) > 0 {
		for _, associatedGroupName := range c.TrafficGroups[trafficGroupName].AssociatedGroups {
			if !stringInSlice(associatedGroupName, c.AssociatedPending) {
				c.AssociatedPending = append(c.AssociatedPending, associatedGroupName)
			}

			go c.handleAssociatedTrafficGroup(associatedGroupName, mc)

			if associatedGroupName[:7] == "/bridge" {
				return
			}
		}
	}

	messages = c.SetTrafficGroupState(trafficGroupName, 2)

	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[0])

	messages = c.SetTrafficGroupState(trafficGroupName, 1)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	messages = c.SetTrafficGroupState(trafficGroupName, 0)
	mc <- messages

	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[2])

	if stringInSlice(trafficGroupName, c.AssociatedPending) {
		c.AssociatedPending = delFromSlice(indexOfStringInSlice(trafficGroupName, c.AssociatedPending), c.AssociatedPending)
	}
}

func (c *Controller) handleBridge(trafficGroupName string, mc chan []mqttMessage) {
	var messages []mqttMessage

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		if trafficItemName[:6] == "/light" {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, 0)...)
		}
	}

	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[0])

	messages = c.SetTrafficItemState(trafficGroupName, "/gate/1", 1)
	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	for c.TrafficGroups[trafficGroupName].Sensors["/sensor/1"] {
		time.Sleep(100 * time.Millisecond)
	}

	messages = c.SetTrafficItemState(trafficGroupName, "/gate/2", 1)
	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	messages = []mqttMessage{}

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		if trafficItemName[:5] == "/deck" {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, 0)...)
		}
	}

	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[2])

	var vesselSensors [][]string
	var vesselGroups []string

	vesselMainSensor := []string{""}

	for sensorName := range c.TrafficGroups[trafficGroupName].Sensors {
		canonicalSensorName := getCanonicalName(trafficGroupName, sensorName)
		if canonicalSensorName[0][:7] == "/vessel" {
			vesselNr, _ := strconv.Atoi(canonicalSensorName[0][8:])
			vesselMainSensor = []string{"/vessel/" + canonicalSensorName[0][8:], "/sensor/1"}
			for i := vesselNr - 1; i >= vesselNr-2; i-- {
				sensorGroupName := "/vessel/" + strconv.Itoa(i)
				vesselGroups = append(vesselGroups, sensorGroupName)

				for sensorName := range c.TrafficGroups[sensorGroupName].Sensors {
					vesselSensors = append(vesselSensors, []string{sensorGroupName, sensorName})
				}
			}
			break
		}
	}

	for _, vesselGroup := range vesselGroups {
		wasGreen := false

		for c.checkTrafficGroupSensors(vesselGroup) || c.TrafficGroups[vesselMainSensor[0]].Sensors[vesselMainSensor[1]] {
			wasGreen = true
			messages = c.SetVesselGroupState(vesselGroup, 2)
			mc <- messages

			time.Sleep(100 * time.Millisecond)
		}

		if wasGreen {
			messages = c.SetVesselGroupState(vesselGroup, 0)
			mc <- messages
		}
	}

	messages = []mqttMessage{}

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		if trafficItemName[:5] == "/deck" {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, 1)...)
		}
	}

	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[2])

	messages = []mqttMessage{}

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		if trafficItemName[:5] == "/gate" {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, 0)...)
		}
	}

	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[1])

	messages = []mqttMessage{}

	for trafficItemName := range c.TrafficGroups[trafficGroupName].Items {
		if trafficItemName[:6] == "/light" {
			messages = append(messages, c.SetTrafficItemState(trafficGroupName, trafficItemName, 2)...)
		}
	}

	mc <- messages
	time.Sleep(c.TrafficGroups[trafficGroupName].Duration[3])
}

func (c *Controller) updateScores() {
	for trafficGroupName, trafficGroup := range c.TrafficGroups {
		if len(trafficGroup.Items) == 0 || trafficGroupName[:7] == "/bridge" {
			continue
		}

		c.ChangesChannel <- TrafficChange{
			[2]string{trafficGroupName, ""},
			"baseScore",
			0,
		}

		trafficItem := getTrafficItemFromGroup(trafficGroupName)

		if (trafficItem.Type == "RTG" || trafficItem.Type == "GR") && trafficItem.State > 0 {
			c.ChangesChannel <- TrafficChange{
				[2]string{trafficGroupName, ""},
				"timeScore",
				0,
			}

			continue
		} else if (trafficItem.Type == "GTR") && trafficItem.State < 2 {
			c.ChangesChannel <- TrafficChange{
				[2]string{trafficGroupName, ""},
				"timeScore",
				0,
			}

			continue
		} else if trafficItem.Type == "RG" && trafficItem.State != 1 {
			c.ChangesChannel <- TrafficChange{
				[2]string{trafficGroupName, ""},
				"timeScore",
				0,
			}
			continue
		}

		for sensorName := range trafficGroup.Sensors {
			canonicalSensorName := getCanonicalName(trafficGroupName, sensorName)
			sensorNumber, _ := strconv.Atoi(canonicalSensorName[1][len(canonicalSensorName[1])-1:])

			if c.TrafficGroups[canonicalSensorName[0]].Sensors[canonicalSensorName[1]] {
				c.ChangesChannel <- TrafficChange{
					[2]string{trafficGroupName, ""},
					"baseScore",
					c.TrafficGroups[trafficGroupName].BaseScore + 1/float64(sensorNumber),
				}
			}
		}

		if trafficGroup.BaseScore >= 1 {
			c.ChangesChannel <- TrafficChange{
				[2]string{trafficGroupName, ""},
				"timeScore",
				c.TrafficGroups[trafficGroupName].TimeScore + 0.01219,
			}
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

		if isValid {
			for _, sensor := range trafficGroup.Sensors {
				if sensor {
					isValid = true
					break
				} else {
					isValid = false
				}
			}
		}

		if isValid {
			for trafficItemName := range trafficGroup.Items {
				if trafficItemName[:6] == "/light" {
					isValid = true
					break
				} else {
					isValid = false
				}
			}
		}

		if len(trafficGroup.Items) == 0 {
			isValid = false
		}

		if trafficGroupName[:7] == "/bridge" {
			isValid = false
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

		//fmt.Println("Group:", trafficGroupName, " - Score:", trafficGroup.BaseScore, trafficGroup.TimeScore)

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

			if isValid {
				for _, sensor := range trafficGroup.Sensors {
					if sensor {
						isValid = true
						break
					} else {
						isValid = false
					}
				}
			}

			if isValid {
				for trafficItemName := range trafficGroup.Items {
					if trafficItemName[:6] == "/light" {
						isValid = true
						break
					} else {
						isValid = false
					}
				}
			}

			if len(trafficGroup.Items) == 0 {
				isValid = false
			}

			if trafficGroupName[:7] == "/bridge" {
				isValid = false
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

	//fmt.Println(c.Pending)

	if c.Pending == 0 {
		c.InMotion = false
	}

	if c.InMotion {
		c.ElapsedTime += 500
	} else {
		c.ElapsedTime = 0

		solution := c.GenerateSolution()

		if solution.Score > 0 {
			//for name, trafficGroup := range c.TrafficGroups {
			//	fmt.Println(name + ":", trafficGroup.TimeScore, "-", trafficGroup.BaseScore, "|", trafficGroup.Sensors)
			//}

			fmt.Println("Groups:", solution.TrafficGroups, "- Score:", solution.Score)
			//fmt.Println()

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
