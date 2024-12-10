### Compiling

We need to tell the assembler that we are expecting a processor that has extended features. 
Below should work for most processors with fpu and NEON support
```bash
as -mfpu=neon-vfpv4 -o main.o main.asm
```

[NEXT -> linking](2_linking.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
