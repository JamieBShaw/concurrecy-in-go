package main

func forSelectLoop() {
	for { // Either loop infinitely or range over something
		select {
		// Do some work with channels
		}
	}
}

//for {
//	select {
//	case <-done:
//		return
//	default:
//	}
	// Do non-preemptable work
//}

//for {
//	select {
//	case <-done:
//	return
//	default:
	// Do non-preemptable work
//	}
//}
