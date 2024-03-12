package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load_problem(path string) (Problem, error) {
	file, err := os.Open(path)
	if err != nil {
		return Problem{}, err
	}
	defer file.Close()

	var p Problem
    p.Name = path
	scanner := bufio.NewScanner(file)

	// Read number of nodes
	scanner.Scan()

    scanner.Scan()
	p.NumberOfNodes, err = strconv.Atoi(scanner.Text())
    if err != nil {
        return Problem{}, fmt.Errorf("failed to parse number of nodes")
    }
    
    // Read number of vehicles
    scanner.Scan()

    scanner.Scan()
    p.NumberOfVehicles, err = strconv.Atoi(scanner.Text())
    if err != nil {
        return Problem{}, fmt.Errorf("failed to parse number of vehicles")
    }
    

    // Read vehicle details
    scanner.Scan()

    p.Vehicles = make([]Vehicle, p.NumberOfVehicles+1)
    for i := 1; i <= p.NumberOfVehicles; i++ {
        currentVehicle := Vehicle{}

        scanner.Scan()
        vehicleDetails := strings.Split(scanner.Text(), ",")
        currentVehicle.Index, err = strconv.Atoi(vehicleDetails[0])
        currentVehicle.HomeNode, err = strconv.Atoi(vehicleDetails[1])
        currentVehicle.StartingTime, err = strconv.Atoi(vehicleDetails[2])
        currentVehicle.Capacity, err = strconv.Atoi(vehicleDetails[3])
        
        currentVehicle.TravelTimes = make([][]int, p.NumberOfNodes+1)
        currentVehicle.TravelCosts = make([][]int, p.NumberOfNodes+1)

        for j := 1; j <= p.NumberOfNodes; j++ {
            currentVehicle.TravelTimes[j] = make([]int, p.NumberOfNodes+1)
            currentVehicle.TravelCosts[j] = make([]int, p.NumberOfNodes+1)
        }

        p.Vehicles[i] = currentVehicle
    }

    // Read number of calls
    scanner.Scan()

    scanner.Scan()
    p.NumberOfCalls, _ = strconv.Atoi(scanner.Text())
    p.CallVehicleMap = make(map[int][]int)

    for i := 1; i <= p.NumberOfCalls; i++ {
        p.CallVehicleMap[i] = make([]int, 0)
    }

    // Read compatible calls with vehicles
    scanner.Scan()

    for i := 0; i < p.NumberOfVehicles; i++ {
        scanner.Scan()

        vehicleDetails := strings.Split(scanner.Text(), ",")
        vehicle_index, _ := strconv.Atoi(vehicleDetails[0])

        for _, call := range vehicleDetails[1:] {
            callIndex, err := strconv.Atoi(call)
            if err != nil {
                return Problem{}, fmt.Errorf("failed to parse call index")
            }

            p.CallVehicleMap[callIndex] = append(p.CallVehicleMap[callIndex], vehicle_index)
        }
    }


    // Read calls
    scanner.Scan()

    p.Calls = make([]Call, p.NumberOfCalls+1)
    for i := 1; i <= p.NumberOfCalls; i++ {
        currentCall := Call{}

        scanner.Scan()
        callDetails := strings.Split(scanner.Text(), ",")

        currentCall.Index, _ = strconv.Atoi(callDetails[0])
        currentCall.OriginNode, _ = strconv.Atoi(callDetails[1])
        currentCall.OriginCostForVehicle = make([]int, p.NumberOfVehicles+1)
        currentCall.OriginTimeForVehicle = make([]int, p.NumberOfVehicles+1)
        currentCall.DestinationNode, _ = strconv.Atoi(callDetails[2])
        currentCall.DestinationCostForVehicle = make([]int, p.NumberOfVehicles+1)
        currentCall.DestinationTimeForVehicle = make([]int, p.NumberOfVehicles+1)
        currentCall.Size, _ = strconv.Atoi(callDetails[3])
        currentCall.CostOfNotTransporting, _ = strconv.Atoi(callDetails[4])

        currentCall.PickupTimeWindow.LowerBound, _ = strconv.Atoi(callDetails[5])
        currentCall.PickupTimeWindow.UpperBound, _ = strconv.Atoi(callDetails[6])

        currentCall.DeliveryTimeWindow.LowerBound, _ = strconv.Atoi(callDetails[7])
        currentCall.DeliveryTimeWindow.UpperBound, _ = strconv.Atoi(callDetails[8])

        if i != currentCall.Index {
            return Problem{}, fmt.Errorf("call index mismatch")
        }

        p.Calls[i] = currentCall
    }
    
    // Read travel times
    scanner.Scan()

    for i := 0; i < p.NumberOfNodes * p.NumberOfNodes * p.NumberOfVehicles; i++ {
        scanner.Scan()
        travelTimes := strings.Split(scanner.Text(), ",")

        vehicleIndex, _ := strconv.Atoi(travelTimes[0])
        originNode, _ := strconv.Atoi(travelTimes[1])
        destinationNode, _ := strconv.Atoi(travelTimes[2])
        travelTime, _ := strconv.Atoi(travelTimes[3])
        travelCost, _ := strconv.Atoi(travelTimes[4])

        p.Vehicles[vehicleIndex].TravelTimes[originNode][destinationNode] = travelTime
        p.Vehicles[vehicleIndex].TravelCosts[originNode][destinationNode] = travelCost
    }

    // Read node times and costs
    scanner.Scan()

    for i := 0; i < p.NumberOfCalls * p.NumberOfVehicles; i++ {
        scanner.Scan()
        callTimes := strings.Split(scanner.Text(), ",")

        vehicleIndex, _ := strconv.Atoi(callTimes[0])
        callIndex, _ := strconv.Atoi(callTimes[1])
        originTime, _ := strconv.Atoi(callTimes[2])
        originCost, _ := strconv.Atoi(callTimes[3])
        destinationTime, _ := strconv.Atoi(callTimes[4])
        destinationCost, _ := strconv.Atoi(callTimes[5])

        p.Calls[callIndex].OriginTimeForVehicle[vehicleIndex] = originTime
        p.Calls[callIndex].OriginCostForVehicle[vehicleIndex] = originCost
        p.Calls[callIndex].DestinationTimeForVehicle[vehicleIndex] = destinationTime
        p.Calls[callIndex].DestinationCostForVehicle[vehicleIndex] = destinationCost
    }

	return p, nil
}

