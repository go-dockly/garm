### Conditionals

If the first condition evaulates to true (c equals zero), only then is the second condition evaluated. To implement the above, one could use the following
```arm
    CMP    c, 0
    B.NE   false

    CMP    x, y
    B.NE   false
true:
    @ body of if statement
false:
    @ end of if statement
```
armv7+
```arm
    cmp     c, 0
    cmp.eq  x, y
    b.ne    false
```
armv8+
```arm
    cmp    c, 0
    ccmp   x, y, 0, eq
    b.ne    false

    @ conditions are true:
false:
```
Ternary operator
```arm
	bEqual = (c == 0) ? (x == y) : 0;
```

Below code is a loop which runs until the counter in X1 hits
zero, at which point the condition code NE (not equal to zero) controlling
the branch becomes false and the loop terminates
```arm
		MOV    X1, #10
	loop
		...
		SUBS   X1, X1, #1
		b.ne   loop
```

sequence of several conditional instructions
```arm
	CMP    X0, #5   @ if (a == 5)
	MOV.EQ X0, #10  @
	BL.EQ  fn       @   fn(10)
```
Assume a is in X0. Compare X0 to 5. The next two instructions will be
executed only if the compare returns EQual. They move 10 into X0, then call
‘fn’ (branch with link, BL)

Set the flags, then use various condition codes:
```arm
	CMP    X0, #0   ; if (x < 0)
	MOV.LE X0, #0   ;   x = 0;
	MOV.GT X0, #1   ; else x = 1;
```

Use conditional compare instructions:
```arm
	CMP    X0, #'A' ; if (c == 'A'
	CMP.NE X0, #'B' ;  || c == 'B')
	MOV.EQ X1, #1   ;   y = 1;
```

A sequence which doesn’t use conditional execution:
```arm
	CMP   X3, #0
	B.EQ  next
	ADD   X0, X0, X1
	SUB   X0, X0, X2
next
	...
```

By transforming the sequence with conditional execution an instruction can be
```arm
removed:
	CMP    X3, #0
	ADD.NE X0, X0, X1
	SUB.NE X0, X0, X2
	...
```

[NEXT -> loops](loop.md)

<div align="center">
	<img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
	<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>