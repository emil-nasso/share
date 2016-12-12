package main

func noCommandHelpMessage() string {
	return "Missing command.\n" +
		"Allowed:\n" +
		"  get - downloads files\n" +
		"  put - shares files\n" +
		"  server - starts as a broker server"
}

func missingGetArgumentHelpMessage() string {
	return "Invalid command. Missing get argument."
}

func missingPutArgumentHelpMessage() string {
	return "Missing path to file.\n" +
		"Usage: share put path/to/file.conf"
}
