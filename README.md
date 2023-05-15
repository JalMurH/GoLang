# GoLang
In this repository you can find the projects developed in Go
## GoGUI
In this project the typical structure of compiled languages is used where we have a src folder where we will find the main file and its modules.
### General objective
A graphical interface was created to monitor the status of hardware components of a server.
### Components
* It has its dockerfile for portability in different operating systems
* The main file main.go where it is sought to make it as clean as possible to be able to keep track of the program correctly
* The data obtained is passed to the parmas.go module where, in its GetParams function, it saves the data obtained thanks to the gopsutil library in a structure for later use in the graph.go module.
* In graph.go the Graph function makes use of the data in the params structure to graph it using the fyne.io library available in the GO documentation
* In the gpu.go module, the capacity of Go is used to use system commands such as lshw to obtain information about the graphics cards
* In users.go, a list of users will be obtained by accessing the system log to keep track of user connections and their respective permissions in the system.
