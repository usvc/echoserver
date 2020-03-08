package main

func init() {
	conf.LoadFromEnvironment()
	command = GetCommand()
}

func main() {
	command.Execute()
}
