package args

func FindFlag(args []string, name string) (string, bool) {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-"+name {
			return args[i+1], true
		}
	}

	return "", false
}

func FlagExist(args []string, name string) bool {
	for _, arg := range args {
		if arg == "-"+name {
			return true
		}
	}

	return false
}
