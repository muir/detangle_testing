# testout

A tiny utility to process the output of Go tests to make it easier to see
a linear flow of what contributed to a failed test.

When you combine t.Run() with t.Parallel(), the output of Go tests can be
very hard to follow because the output from many tests can be intermixed.

This utiltity detangles the output.
