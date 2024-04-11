# This is my solution for the inf273 Pickup and Delivery Problem
## Structure
Currently all in one package. See types.go for most datatypes. Most of the logic is contained in operators.go and find_feasible.go.
### Data directory
This contains all the case files.

### results directory
This directory contains the data for every run named by the time it was run.

It also contains `results.txt` wich is an automatically generated file in the format of the assigments.

## Dependencies
- The go runtime
- python (to automatically get the tables for the latest run)

## How to run
Write `go run .` in the root directory.

## To run test
`go test -v`
