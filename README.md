## First thoughts
- It shouldn't attempt to read a large string into memory before processing it as this obviously wouldn't scale.
- A good level of unit testing should cover the statistic generation.
- The app should be composed of simple, easily testable functions. DI may be used.
- Care with race conditions needs to be taken when accessing state. Go channels are likely the best option. They can be used to implement a kind of worker queue.
- It may be possible to use the same data structure as Redis' uses for the realtime leaderboard I built last year at Chroma: [https://en.wikipedia.org/wiki/Skip_list](https://en.wikipedia.org/wiki/Skip_list)
- There could be extreme situations where a payload will need to be rejected due to scarce resources (i.e. full queue). In this case, an appropriate HTTP response code like 503 should be returned.

## MVP
An initial MVP was built using mutexes instead of channels with the aim of meeting the requirements without investing much time in optimisations.

It has a good level of test coverage that will help spot any regressions during future optimisations or adding of new features.

## Optimisation
The original approach of using a `map[string]int` to store the occurances of words/letters had very fast writes, but much slower queries.

```sh
BenchmarkWrite-4   	 2000000	       858 ns/op	      15 B/op	       0 allocs/op
```

```sh
BenchmarkQuery-4   	      20	  73913326 ns/op	     160 B/op	       2 allocs/op
```

I looked for an efficient sorted set package and found [wangjia184's](https://godoc.org/github.com/wangjia184/sortedset). It has O(log(N)) complexity for the  operations this program uses. This results in slightly slower writes, but a much faster ability to calculate the most common _N_ words during queries.

```sh
BenchmarkWrite-4   	 1000000	      1889 ns/op	      69 B/op	       0 allocs/op
```

```sh
BenchmarkQuery-4   	  500000	      2157 ns/op	     720 B/op	      16 allocs/op
```