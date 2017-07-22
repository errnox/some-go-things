// Rob Pike, Self-referenctial Functions and the Design of Options
//
// (Blog Post)
//
//
// See here for more info:
//
//   https://commandcenter.blogspot.de/2014/01/self-referential-functions-and-design.html







// ---------------------------------------------------------------------


// Method 1: Basic


// Option sets the options specified.
func (f *Foo) Option(opts ...option) {
    for _, opt := range opts {
        opt(f)
    }
}

// Verbosity sets Foo's verbosity level to v.
func Verbosity(v int) option {
    return func(f *Foo) {
        f.verbosity = v
    }
}

foo.Option(pkg.Verbosity(3))







// ---------------------------------------------------------------------

// Method 2: Enhanced


type option func(*Foo) interface{}

// Verbosity sets Foo's verbosity level to v.
func Verbosity(v int) option {
    return func(f *Foo) interface{} {
        previous := f.verbosity
        f.verbosity = v
        return previous
    }
}

// Option sets the options specified.
// It returns the previous value of the last argument.
func (f *Foo) Option(opts ...option) (previous interface{}) {
    for _, opt := range opts {
        previous = opt(f)
    }
    return previous
}







// ---------------------------------------------------------------------

// Method 2 allows restoring the previous value.

prevVerbosity := foo.Option(pkg.Verbosity(3))
foo.DoSomeDebugging()
foo.Option(pkg.Verbosity(prevVerbosity.(int)))







// ---------------------------------------------------------------------

// Method 3: Even Better


type option func(f *Foo) option

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (f *Foo) Option(opts ...option) (previous option) {
    for _, opt := range opts {
        previous = opt(f)
    }
    return previous
}

// Verbosity sets Foo's verbosity level to v.
func Verbosity(v int) option {
    return func(f *Foo) option {
        previous := f.verbosity
        f.verbosity = v
        return Verbosity(previous)
    }
}

prevVerbosity := foo.Option(pkg.Verbosity(3))
foo.DoSomeDebugging()
foo.Option(prevVerbosity)

func DoSomethingVerbosely(foo *Foo, verbosity int) {
    // Could combine the next two lines,
    // with some loss of readability.
    prev := foo.Option(pkg.Verbosity(verbosity))
    defer foo.Option(prev)
    // ... do some stuff with foo under high verbosity.
}
