package main
/*age main*/

/*import (*/
	/*"fmt"*/
	/*"testing"*/
/*)*/


/*func TestGenerateInitialSolution(t *testing.T) {*/
	/*p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")*/
    /*if err != nil {*/
        /*t.Error(err)*/
    /*}*/

    /*s := p.GenerateInitialSolution()*/
	/*solution := s.Solution*/
    /*fmt.Println(solution)*/

	/*expected := []int{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7}*/
	/*for i := range solution {*/
		/*if solution[i] != expected[i] {*/
			/*t.Errorf("Expected %d, got %d", expected[i], solution[i])*/
		/*}*/
	/*}*/
/*}*/

/*func TestGenerateRandomSolution(t *testing.T) {*/
	/*p := Problem{*/
        /*NumberOfVehicles: 2, */
        /*NumberOfCalls: 3,*/
        /*CallVehicleMap: map[int][]int{1: {1, 2}, 2: {1, 2}, 3: {1, 2}},*/
    /*}*/
	/*solution := p.GenerateRandomSolution().Solution*/

	/*if len(solution) != 8 {*/
		/*t.Errorf("Expected 8, got %d", len(solution))*/
	/*}*/
/*}*/

/*func TestMoveFromOutsource(t *testing.T) {*/
	/*p := Problem{NumberOfVehicles: 1, NumberOfCalls: 1, CallVehicleMap: map[int][]int{1: {1}}}*/
    /*s := p.GenerateInitialSolution()*/

    /*solution := s.Solution*/
	/*expected := []int{1, 1, 0}*/

	/*s.moveFromOutsource([]int{1, 2}, []int{0})*/

	/*for i := range solution {*/
		/*if solution[i] != expected[i] {*/
			/*t.Errorf("Expected %d, got %d", expected[i], solution[i])*/
		/*}*/
	/*}*/

	/*p = Problem{NumberOfVehicles: 2, NumberOfCalls: 2, CallVehicleMap: map[int][]int{1: {1, 2}, 2: {2}}}*/

	/*s = p.GenerateInitialSolution()*/
    /*solution = s.Solution*/
	/*expectedOptions := [][]int{*/
		/*{0, 1, 1, 0, 2, 2},*/
		/*{1, 1, 0, 0, 2, 2},*/
		/*{0, 2, 2, 0, 1, 1},*/
		/*{2, 2, 0, 0, 1, 1},*/
	/*}*/

	/*s.moveFromOutsource([]int{2, 3}, []int{0, 1})*/

	/*found := matchAnySlice(solution, expectedOptions)*/
	/*if !found {*/
		/*t.Errorf("No match found")*/
	/*}*/

/*}*/

/*func TestMoveCallInVehicle(t *testing.T) {*/
	/*p := &Problem{NumberOfVehicles: 1, NumberOfCalls: 2}*/

    /*s := Solution{*/
        /*Problem: p,*/
        /*Solution: []int{1, 2, 1, 2, 0},*/
    /*}*/

    /*solution := s.Solution*/
	/*expectedOptions := [][]int{*/
		/*{2, 1, 1, 2, 0},*/
		/*{1, 1, 2, 2, 0},*/
		/*{1, 2, 2, 1, 0},*/
	/*}*/
	/*s.moveCallInVehicle([]int{0, 2}, []int{4})*/

	/*found := matchAnySlice(solution, expectedOptions)*/
	/*if !found {*/
		/*t.Errorf("No match found")*/
	/*}*/
/*}*/

/*func TestOneReinsert(t *testing.T) {*/
    /*p := &Problem{NumberOfVehicles: 2, NumberOfCalls: 1}*/

    /*s := Solution{*/
        /*Problem: p,*/
        /*Solution: []int{0, 1, 1, 0},*/
    /*}*/

    /*solution := s.Solution*/
    /*expected := []int{0,0,1,1}*/

    /*s.OneReinsert()*/
    
    /*if !matchSlice(solution, expected) {*/
        /*t.Errorf("Expected %v, got %v", expected, solution)*/
    /*}*/
/*}*/

