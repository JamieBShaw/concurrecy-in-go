package main

// CSP (Communicating Sequential Processes)

// As we discussed in the section “The Difference Between Concurrency and
// Parallelism” for modeling concurrent problems, it’s common for languages to
// end their chain of abstraction at the level of the OS thread and memory access
// synchronization. Go takes a different route and supplants this with the
// concept of goroutines and channels.
// If we were to draw a comparison between concepts in the two ways of
// abstracting concurrent code, we’d probably compare the goroutine to a
// thread, and a channel to a mutex (these primitives only have a passing
// resemblance, but hopefully the comparison helps you get your bearings)

// ---------------------------------------------------------------------------------

// Goroutines are for executing code off the main routine
// Channels are for communicating between go routines

// Go follows a model of concurrency called the fork-join model
// The word fork refers to the fact that at any point in the program, it can split off a child
// branch of execution to be run concurrently with its parent. The word join
// refers to the fact that at some point in the future, these concurrent branches of
// execution will join back together. Where the child rejoins the parent is called
// a join point.
// Creating a goroutine with go command is a fork

// ---------------------------------------------------------------------------------

// CHANNELS

// Remember that when we discussed blocking, we said that writes to a channel block if a channel is
// full, and reads from a channel block if the channel is empty? “Full” and
// “empty” are functions of the capacity, or buffer size. An unbuffered channel
// has a capacity of zero and so it’s already full before any writes. A buffered
// channel with no receivers and a capacity of four would be full after four
// writes, and block on the fifth write since it has nowhere else to place the fifth
// element. Like unbuffered channels, buffered channels are still blocking; the
// preconditions that the channel be empty or full are just different. In this way,
// buffered channels are an in-memory FIFO queue for concurrent processes to communicate over.

// Unbuffered Channels

// Unbuffered channels is a channel that initially has no capacity.
// Unbuffered Channel will block the goroutine whenever it is empty and waiting to be filled
// By default communication over the channels is sync, when you send some value there must
// be a receiver. Otherwise, fatal error: all goroutines are asleep - deadlock!

// Buffered Channels

//- Buffered Channel will block the goroutine either when it is empty and waiting to be filled or
//  it's full-capacity and there's a statement that wants to be sent into the channel.
//- While there is space in the buffered channel, sending into the channel is a NON-BLOCKING operation.

// Tips for using channels
//
// Channel 'Ownership', which goroutine is responsible for the channel.
// Unidirectional channel declarations are the tool that will allow us to distinguish between
// goroutines that own channels and those that only utilize them: channel
// owners have a write-access view into the channel (chan or chan<-), and
// channel utilizes only have a read-only view into the channel (<-chan)

// The goroutine that owns a channel should:
// - Instantiate the channel.
// - Perform writes, or pass ownership to another goroutine
// - Close the channel
// - Encapsulate the previous three things in this list and expose them via a read channel

// As a consumer of a channel:
// - Knowing when a channel is blocked
// - Responsibility handling blocking for any reason

// Sync package

// You can think of a WaitGroup like a concurrent-safe counter: calls to Add
// increment the counter by the integer passed in, and calls to Done decrement
// the counter by one. Calls to Wait block until the counter is zero.
// Notice that the calls to Add are done outside the goroutines they’re helping to
// track. If we didn’t do this, we would have introduced a race condition,
// because remember from “Goroutines” that we have no guarantees about
// when the goroutines will be scheduled; we could reach the call to Wait before
// either of the goroutines begin. Had the calls to Add been placed inside the
// goroutines’ closures, the call to Wait could have returned without blocking at
// all because the calls to Add would not have taken place

// ---------------------------------------------------

// Utilising concurrency

// When working with concurrent code, there are a few different options for
// safe operation. We’ve gone over two of them:

// - Synchronization primitives for sharing memory (e.g., sync.Mutex)
// - Synchronization via communicating (e.g., channels)

// However, there are a couple of other options that are implicitly safe within
// multiple concurrent processes:

// - Immutable data
// - Data protected by confinement

// Confinement

// Confinement is the simple yet powerful idea of ensuring information is only
// ever available from one concurrent process. When this is achieved, a
// concurrent program is implicitly safe and no synchronization is needed.
// There are two kinds of confinement possible: ad hoc and lexical.

// Ad-hoc is more about setting community standards, repo standards, work standards
// around the best practice to write concurrent code, however, this isn't particularly strict

// Lexical confinement, means writing code that makes the concurrent part lexically separate or
// clearer to current and future developers how the concurrency is working in the code.
// e.g. returning read only channels, creating functions who are clear owners and consumers of
// channels etc.

// The for-select loop

// Sending iteration variables out on a channel
// Oftentimes you’ll want to convert something that can be iterated over
// into values on a channel. This is nothing fancy, and usually looks
// something like this:
// for _, s := range []string{"a", "b", "c"} {
//     select {
//     case <-done:
//         return
//     case stringStream <- s:
//     }
// }

// Looping infinitely waiting to be stopped
// It’s very common to create goroutines that loop infinitely until they’re
// stopped. There are a couple variations of this one. Which one you
// choose is purely a stylistic preference.
// The first variation keeps the select statement as short as possible:
//  for {
//      select {
//      case <-done:
//          return
//      default:
//      }
//  }

// If the done channel isn’t closed, we’ll exit the select statement and
// continue on to the rest of our for loop’s body.
// The second variation embeds the work in a default clause of the
// select statement:
//  for {
//      select {
//      case <-done:
//          return
//      default:
//      Do non-preemptable work
//      }
//  }
// When we enter the select statement, if the done channel hasn’t been
// closed, we’ll execute the default clause instead.
// There’s nothing more to this pattern, but it shows up all over the place,
// and so it’s worth mentioning.

// Preventing Goroutine Leaks

// We could even represent this interconnectedness as a graph: whether or not a
// child goroutine should continue executing might be predicated on knowledge
// of the state of many other goroutines. The parent goroutine (often the main
// goroutine) with this full contextual knowledge should be able to tell its child
// goroutines to terminate.

// Now that we know how to ensure goroutines don’t leak, we can stipulate a
// convention:
// - If a goroutine is responsible for creating a goroutine, it is also
//	 responsible for ensuring it can stop the goroutine.

// Pipelines

// That’s interesting; what are the properties of a pipeline stage?
// - A stage consumes and returns the same type.
// - A stage must be reified2 by the language so that it may be passed
//	 around. Functions in Go are reified and fit this purpose nicely.

// Best Practices for Constructing Pipelines


// Fan-Out, Fan-In
// Fan-out is a term to describe the process of starting multiple goroutines to
// handle input from the pipeline, and fan-in is a term to describe the process of
// combining multiple results into one channel.

//  You might consider fanning out one of your stages if both of the following apply:
//  It doesn’t rely on values that the stage had calculated before.
//  It takes a long time to run.
//  The property of order-independence is important because you have no
//  guarantee in what order concurrent copies of your stage will run, nor in what
//  order they will return.

// Fan-In:
// Fanning in involves creating the multiplexed channel consumers
// will read from, and then spinning up one goroutine for each incoming
// channel, and one goroutine to close the multiplexed channel when the
// incoming channels have all been closed. Since we’re going to be creating a
// goroutine that is waiting on N other goroutines to complete, it makes sense to
// create a sync.WaitGroup to coordinate things. The multiplex function also
// notifies the WaitGroup that it’s done.

// The or-done-channel
// At times you will be working with channels from disparate parts of your
// system. Unlike with pipelines, you can’t make any assertions about how a
// channel will behave when code you’re working with is canceled via its done
// channel. That is to say, you don’t know if the fact that your goroutine was
// canceled means the channel you’re reading from will have been canceled. For
// this reason, as we laid out in “Preventing Goroutine Leaks”, we need to wrap
// our read from the channel with a select statement that also selects from a
// done channel. This is perfectly fine, but doing so takes code that’s easily read
// like this:
// for val := range myChan {
//      Do something with val
// }
// And explodes it out into this:
// loop:
// for {
//     select {
//     case <-done:
//         break loop
//     case maybeVal, ok := <-myChan:
//         if ok == false {
//             return // or maybe break from for
//         }
//         // Do something with val
//     }
// }
// The below is used as a helper function to resolve this:
// func orDone(done, c <-chan interface{}) <-chan interface{} {
//	valStream := make(chan interface{})
//	go func() {
//		defer close(valStream)
//		for {
//			select {
//			case <-done:
//				return
//			case v, ok := <-c:
//				if !ok {
//					return
//				}
//				select {
//				case valStream <- v:
//				case <-done:
//				}
//			}
//		}
//	}()
//	return valStream
//}

// The tee-channel

// Much like the unix tee command, tee-channel patterns takes a value from a channel
// and then enters it into multiple other channels so they can be used elsewhere