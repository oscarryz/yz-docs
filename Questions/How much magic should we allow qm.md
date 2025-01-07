

The driving design is to avoid "magic" by making the syntax as self explanatory as possible ,  and avoid surprises.  By magic I mean things that the programmer didn't write and still are there present (not to confuse with library code that will be explicitly available in the library code ).

It is however impractical require the programmer to hand write convenient functionality all the time, a clear example is the keywork `self` or `this`  other languages have or autoloading the standard library with an import. 

For these two specific examples currently the strategy is to explicitly create a self variable and set it upon creation (not terrible). And manually importing 


