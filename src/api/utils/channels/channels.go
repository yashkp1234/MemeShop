package channels

//OK checks if a channel is done
func OK(done chan bool) bool {
	select {
	case ok := <-done:
		if ok {
			return true
		}
	}
	return false
}
