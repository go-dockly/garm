### Branching

[Azeria on branching](https://azeria-labs.com/arm-conditional-execution-and-branching-part-6/)

Branching forward, to skip over some code:
```arm
	...
	B fwd		@  jump to label 'fwd'
	...
fwd
```
B instruction to unconditionally branch to label at PC-relative offset, with hint that this is not a subroutine call or return

Branching backwards, creating a loop:
```arm
back
	...		
	B back		@ jump to label 'back'
```

Using BL to call a subroutine: 
```arm
	...
	BL  calc	@ call 'calc'
	...

calc			@ function body
	ADD r0,r1,r2	@ do some work here
	MOV pc, r14	@ PC = R14 to return

	ENDs
```
Branch with Link branches to a PC-relative offset, setting register X30 to PC+4 with hint that this is a subroutine call


[NEXT -> memory](memory.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
		<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>