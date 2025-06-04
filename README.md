# container-from-scratch
Go Code to implement containers from scratch

Added features:
1. Container process has capabilities similar to how runc sets them 

Inspired by the very-famous cfs videos of the one and only Liz Rice

To run:
` go run main.go run <command>`

Ex: go run main.go run sh

Note: This might not run in WSL because of how its configured leading to permission issues
