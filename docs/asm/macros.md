## Macros

Let's take the famous ternary example 

```arm
cmp w0, #0
csel w0, w1, w2, eq       // w0 = (w0 == 0) ? w1 : w2
```

and turn it into a macro

```arm
.macro SELECT_IF_ZERO dst, test, true_val, false_val
    cmp     \test, #0
    csel    \dst, \true_val, \false_val, eq
.endm
//use it like so
SELECT_IF_ZERO w0, w0, w1, w2    // w0 = (w0 == 0) ? w1 : w2
```
This macro takes 4 parameters: destination register, test register, value if true, value if false. Makes the comparison with zero and
performs the conditional select based on equality.

Now let's make it more flexible and compare against any value, not just zero to make it feel more like a C-style ternary operator (condition ? true_val : false_val)

```arm
.macro TERNARY dst_cond, question, true_val, false_val
    cmp     \dst_cond, \question    // condition ?
    csel    \dst_cond, \true_val, \false_val, eq   // true_val : false_val
.endm
//use it like so
TERNARY w0, #0, w1, w2    // w0 = (w0 == 0) ? w1 : w2
TERNARY w0, w1, w2, w3    // w0 = (w0 == w1) ? w2 : w3
```
Much cleaner! The first parameter acts as the value being tested and the destination for result, just like a typical ternary expression.