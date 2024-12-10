### Program:

```go
package test

func addInt32s() []int32 {
	var a, b, result [4]int32

	a = [4]int32{1, 2, 3, 4}
	b = [4]int32{5, 6, 7, 8}

	for i := 0; i < 4; i++ {
		result[i] = a[i] + b[i]
	}

	return result[:]
}
```

```asm
addInt32s:
    // Load the first vector of 4 int32 from memory into V0
    LDP Q0, Q1, [X0]
    
    // Load the second vector of 4 int32 from memory into V1
    LDP Q2, Q3, [X1]
    
    // Perform vector addition using SIMD instructions
    ADD V0.4S, V0.4S, V2.4S
    
    // Store the result back to memory
    STP Q0, Q1, [X2]
    
    RET
```

### Breakdown:

```asm
LDP Q0, Q1, [X0]: Load Pair instruction in ARM64
```
Loads two 128-bit SIMD registers (Q0 and Q1) from the memory address in X0
This captures the first vector (slice a)

```asm
LDP Q2, Q3, [X1]: Load Pair for the second vector
```
- Loads two 128-bit SIMD registers (Q2 and Q3) from the memory address in X1
- This captures the second vector (slice b)

```asm
ADD V0.4S, V0.4S, V2.4S: Vector Addition
```
- .4S indicates 4x 32-bit signed integers
- Performs element-wise addition
- Result stored in V0 (equivalent to Q0 in 32-bit mode)

```asm
STP Q0, Q1, [X2]: Store Pair instruction
```
- Stores the result back to memory at address contained in X2
- Stores two 128-bit registers

```asm
RET // Return from the function
```


### Memory Layout

	•	Assume R0, R1, and R2 hold the starting addresses of arrays a, b, and result, respectively.
	•	Register Q0 and Q1 hold four 32-bit integers each, as shown below:

|Q0 register (initial) | Q1 register (initial) | Result in Q0 (after VADD)|
|----------------------|-----------------------|--------------------------|
| 1 | 5 |  6 |
| 2 | 6 |  8 |
| 3 | 7 | 10 |
| 4 | 8 | 12 |

### AST
```
   237  .  .  .  .  .  .  .  List: []ast.Stmt (len = 1) {
   238  .  .  .  .  .  .  .  .  0: *ast.AssignStmt {
   239  .  .  .  .  .  .  .  .  .  Lhs: []ast.Expr (len = 1) {
   240  .  .  .  .  .  .  .  .  .  .  0: *ast.IndexExpr {
   241  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   242  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:3
   243  .  .  .  .  .  .  .  .  .  .  .  .  Name: "result"
   244  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 163)
   245  .  .  .  .  .  .  .  .  .  .  .  }
   246  .  .  .  .  .  .  .  .  .  .  .  Lbrack: test/loop.go:10:9
   247  .  .  .  .  .  .  .  .  .  .  .  Index: *ast.Ident {
   248  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:10
   249  .  .  .  .  .  .  .  .  .  .  .  .  Name: "i"
   250  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 195)
   251  .  .  .  .  .  .  .  .  .  .  .  }
   252  .  .  .  .  .  .  .  .  .  .  .  Rbrack: test/loop.go:10:11
   253  .  .  .  .  .  .  .  .  .  .  }
   254  .  .  .  .  .  .  .  .  .  }
   255  .  .  .  .  .  .  .  .  .  TokPos: test/loop.go:10:13
   256  .  .  .  .  .  .  .  .  .  Tok: =
   257  .  .  .  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
   258  .  .  .  .  .  .  .  .  .  .  0: *ast.BinaryExpr {
   259  .  .  .  .  .  .  .  .  .  .  .  X: *ast.IndexExpr {
   260  .  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   261  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:15
   262  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
   263  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 53)
   264  .  .  .  .  .  .  .  .  .  .  .  .  }
   265  .  .  .  .  .  .  .  .  .  .  .  .  Lbrack: test/loop.go:10:16
   266  .  .  .  .  .  .  .  .  .  .  .  .  Index: *ast.Ident {
   267  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:17
   268  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "i"
   269  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 195)
   270  .  .  .  .  .  .  .  .  .  .  .  .  }
   271  .  .  .  .  .  .  .  .  .  .  .  .  Rbrack: test/loop.go:10:18
   272  .  .  .  .  .  .  .  .  .  .  .  }
   273  .  .  .  .  .  .  .  .  .  .  .  OpPos: test/loop.go:10:20
   274  .  .  .  .  .  .  .  .  .  .  .  Op: +
   275  .  .  .  .  .  .  .  .  .  .  .  Y: *ast.IndexExpr {
   276  .  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   277  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:22
   278  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
   279  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 108)
   280  .  .  .  .  .  .  .  .  .  .  .  .  }
   281  .  .  .  .  .  .  .  .  .  .  .  .  Lbrack: test/loop.go:10:23
   282  .  .  .  .  .  .  .  .  .  .  .  .  Index: *ast.Ident {
   283  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: test/loop.go:10:24
   284  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "i"
   285  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 195)
   286  .  .  .  .  .  .  .  .  .  .  .  .  }
   287  .  .  .  .  .  .  .  .  .  .  .  .  Rbrack: test/loop.go:10:25
   288  .  .  .  .  .  .  .  .  .  .  .  }
   289  .  .  .  .  .  .  .  .  .  .  }
   290  .  .  .  .  .  .  .  .  .  }
   291  .  .  .  .  .  .  .  .  }
   292  .  .  .  .  .  .  .  }
   293  .  .  .  .  .  .  .  Rbrace: test/loop.go:11:2
   294  .  .  .  .  .  .  }
   ```

[NEXT -> GC](gc.md)

<div align="center">
  <img src="docs/img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
        <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>