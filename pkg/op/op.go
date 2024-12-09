package op

type Op string

// ARM64 instruction set
const (
	// No Operation does nothing, other than advance the value of program counter by 4.
	// Instruction can be used for instruction alignment purposes.
	NOP Op = "NOP"
	// Invalid operation
	INVALID = "INVALID"
	// Supervisor call
	SVC Op = "SVC"
	// Load memory eg R0 = [0x1234]
	LOAD Op = "LOAD"
	// Store memory eg [0x1234] = R0
	STORE Op = "STORE"
	// Store Pair
	STP Op = "STP"
	// Load Pair
	LDP Op = "LDP"
	// Load byte eg R0 = [0xFF]
	LDRB Op = "LDRB"
	// Store byte eg [0xFF] = R0
	STRB Op = "STRB"
	// Load halfword eg R0 = [0x1234]
	LDRH Op = "LDRH"
	// Store halfword eg [0x1234] = R0
	STRH Op = "STRH"
	// Conditional select eg R0 = R1 if condition else R0
	CSEL Op = "CSEL"
	// Prefetch memory eg [0x1234]
	PRFM Op = "PRFM"
	// Address of Page
	ADRP Op = "ADRP"

	// Atomic Memory instructions eg for channels and locks

	// Load Exclusive Register (atomic read)
	LDXR Op = "LDXR"
	// Store Exclusive Register (atomic write)
	STXR Op = "STXR"

	// Register instructions

	// Move register eg R0 = R1
	MOV Op = "MOV"
	// Move register with zero eg R0 = 0x1234
	MOVZ Op = "MOVZ"
	// Move register with negation eg R0 = ~0x1234
	MOVN Op = "MOVN"
	// Move register with keep eg R0 = R0 | 0x1234
	MOVK Op = "MOVK"

	// Atomic Operations

	// Load-Acquire and Add
	LDADD Op = "LDADD" // eg R0 = [R1] + R2
	// Load-Acquire and Clear
	LDCLR Op = "LDCLR" // eg R0 = [R1] &^ R2
	// Load-Acquire and Exclusive OR
	LDEOR Op = "LDEOR" // eg R0 = [R1] ^ R2
	// Load-Acquire and Set
	LDSET Op = "LDSET" // eg R0 = [R1] | R2
	// Swap
	SWP Op = "SWP" // eg R0 = [R1]; [R1] = R2

	// Conditional instructions

	// Compare registers
	CMP Op = "CMP" // eg R1 == R2
	// Compare negative
	CMN Op = "CMN" // eg R1 < 0
	// Test bits
	TST Op = "TST" // eg R1 & R2
	// Test equivalence
	TEQ Op = "TEQ" // eg R1 == R2
	// Compare Not Equal
	CMPNE Op = "CMP.NE" // eg R0 = R1 != R2
	// Compare Carry Set
	CMPCS Op = "CMP.CS" // eg R0 = R1 > R2
	// Compare Higher
	CMPHI Op = "CMP.HI" // eg R0 = R1 > R2 unsigned
	// Compare Greater Than or Equal
	CMPGE Op = "CMP.GE" // eg R0 = R1 >= R2
	// Compare Less Than
	CMPLT Op = "CMP.LT" // eg R0 = R1 < R2
	// Compare Less Than or Equal
	CMPLE Op = "CMP.LE" // eg R0 = R1 <= R2
	// Compare Greater Than
	CMPGT Op = "CMP.GT" // eg R0 = R1 > R2

	// Conditional Branch Instructions
	// Branch if Equal (Z=1)
	BEQ Op = "B.EQ"
	// Branch if Not Equal (Z=0)
	BNE Op = "B.NE"
	// Branch if Carry Set/Higher or Same (C=1)
	BCSHS Op = "B.CS/HS"
	// Branch if Carry Clear/Lower (C=0)
	BCCLO Op = "B.CC/LO"
	// Branch if Minus/Negative (N=1)
	BMI Op = "B.MI"
	// Branch if Plus/Positive or Zero (N=0)
	BPL Op = "B.PL"
	// Branch if Overflow Set (V=1)
	BVS Op = "B.VS"
	// Branch if Overflow Clear (V=0)
	BVC Op = "B.VC"
	// Branch if Higher (unsigned) (C=1 && Z=0)
	BHI Op = "B.HI"
	// Branch if Lower or Same (unsigned) (C=0 || Z=1)
	BLS Op = "B.LS"
	// Branch if Greater than or Equal (signed) (N=V)
	BGE Op = "B.GE"
	// Branch if Less Than (signed) (N!=V)
	BLT Op = "B.LT"
	// Branch if Greater Than (signed) (Z=0 && N=V)
	BGT Op = "B.GT"
	// Branch if Less than or Equal (signed) (Z=1 || N!=V)
	BLE Op = "B.LE"
	// Branch Always (can be written as just B)
	BAL Op = "B.AL" // eg B.AL my_label

	// Unconditional Branch Instructions
	// Branch
	B Op = "B" // eg B my_label
	// Branch with Link (for subroutine calls)
	BL Op = "BL" // eg BL my_label
	// Branch to Register
	BR Op = "BR" // eg BR R1
	// Branch with Link to Register
	BLR Op = "BLR" // eg BLR R1
	// Return from subroutine (specialized branch)
	RET Op = "RET" // eg RET R1

	// Compare and Branch Instructions

	// Compare and Branch if Zero
	CBZ Op = "CBZ" // eg if R0 == 0
	// Compare and Branch if Not Zero
	CBNZ Op = "CBNZ" // eg if R0 != 0
	// Test bit and Branch if Zero
	TBZ Op = "TBZ" // eg if R0 & 0x1
	// Test bit and Branch if Not Zero
	TBNZ Op = "TBNZ" // eg if R0 & 0x1 == 0

	// Arithmetic instructions
	ADD  Op = "ADD"  // Addition eg R0 = R1 + R2
	MADD Op = "MADD" // Multiply and add eg R0 = R1 * R2 + R3
	SUB  Op = "SUB"  // Subtraction eg R0 = R1 - R2
	MSUB Op = "MSUB" // Multiply and subtract eg R0 = R1 * R2 - R3
	MUL  Op = "MUL"  // Multiplication eg R0 = R1 * R2
	ADC  Op = "ADC"  // Add with carry eg R0 = R1 + R2 + carry flag
	SBC  Op = "SBC"  // Subtract with carry eg R0 = R1 - R2 - carry flag
	RSB  Op = "RSB"  // Reverse subtract eg R0 = R2 - R1
	RSC  Op = "RSC"  // Reverse subtract with carry eg R0 = R2 - R1 - carry flag
	MLA  Op = "MLA"  // Multiply and accumulate eg R0 = R1 * R2 + R3
	MLS  Op = "MLS"  // Multiply and subtract eg R0 = R1 * R2 - R3
	SDIV Op = "SDIV" // Signed divide and check for divide by zero
	UDIV Op = "UDIV" // Unsigned divide and check for divide by zero

	// Saturation Arithmetic Instructions
	QADD    Op = "QADD"    // Saturating add eg R0 = R1 + R2
	QSUB    Op = "QSUB"    // Saturating subtract eg R0 = R1 - R2
	QDAD    Op = "QDAD"    // Saturating double add eg R0 = 2*R1 + R2
	QDSB    Op = "QDSB"    // Saturating double subtract eg R0 = 2*R1 - R2
	SSAT    Op = "SSAT"    // Signed saturate eg R0 = sat(R1)
	USAT    Op = "USAT"    // Unsigned saturate eg R0 = sat(R1)]
	QRDMLAH Op = "QRDMLAH" // Saturating Rounding Doubling Multiply Accumulate Returning High Half
	QRDMLSH Op = "QRDMLSH" // Saturating Rounding Doubling Multiply Subtract Returning High Half

	// Logical instructions
	AND  Op = "AND"  // Bitwise AND eg R0 = R1 & R2
	OR   Op = "OR"   // Bitwise OR eg R0 = R1 | R2
	ORR  Op = "ORR"  // Bitwise OR with flags eg R0 = R1 | R2 | flags
	EOR  Op = "EOR"  // Bitwise XOR
	XOR  Op = "XOR"  // Bitwise XOR (alias EOR Exclusive OR) eg R0 = R1 ^ R2
	BIC  Op = "BIC"  // Bitwise AND NOT eg R0 = R1 &^ R2
	BICS Op = "BICS" // Bitwise AND NOT with flags eg R0 = R1 &^ R2 | flags
	MVN  Op = "MVN"  // Move NOT (bitwise NOT) eg R0 = ^R1
	CLZ  Op = "CLZ"  // Count leading zeros

	// Shift instructions
	ASR Op = "ASR" // Arithmetic shift right
	LSR Op = "LSR" // Logical shift right eg R0 = R1 >> 2
	LSL Op = "LSL" // Logical shift left eg R0 = R1 << 2
	ROR Op = "ROR" // Rotate right eg R0 = R1 rotated right by 2
	RRX Op = "RRX" // Rotate right with extend eg R0 = R1 rotated right by 1 with carry flag

	// Bit Manipulation Instructions
	BFXIL Op = "BFXIL" // Bitfield Extract and Insert Low eg R0 = R1[7:0]
	BFI   Op = "BFI"   // Bitfield Insert eg R0 = R1[7:0] << 8
	BFX   Op = "BFX"   // Bitfield Extract eg R0 = R1[7:0] >> 8

	// System Instructions
	SYS  Op = "SYS"  // System instruction
	SYSL Op = "SYSL" // System instruction with result

	// Cache Maintenance Instructions
	DC Op = "DC" // Data Cache operation eg flush
	IC Op = "IC" // Instruction Cache operation eg invalidate

	// SIMD instructions
	// Arithmetic Operations
	MLAS    Op = "MLAS"    // Multiply Accumulate/Subtract eg R0 = R1 * R2 + R3
	FMAS    Op = "FMAS"    // Floating-point Multiply Accumulate/Subtract eg R0 = R1 * R2 + R3
	NEG     Op = "NEG"     // Negate eg R0 = -R1
	ABS     Op = "ABS"     // Absolute eg R0 = |R1|
	MAX     Op = "MAX"     // Maximum eg R0 = max(R1, R2)
	MIN     Op = "MIN"     // Minimum eg R0 = min(R1, R2)
	ACGE    Op = "ACGE"    // Absolute Compare Greater Than or Equal eg R0 = |R1| >= |R2|
	ACGT    Op = "ACGT"    // Absolute Compare Greater Than eg R0 = |R1| > |R2|
	FMA     Op = "FMA"     // Fused Multiply Accumulate eg R0 = R1 * R2 + R3
	FMS     Op = "FMS"     // Fused Multiply Subtract eg R0 = R1 * R2 - R3
	FNMA    Op = "FNMA"    // Fused Negate Multiply Accumulate eg R0 = -R1 * R2 + R3
	FNMS    Op = "FNMS"    // Fused Negate Multiply Subtract eg R0 = -R1 * R2 - R3
	NMLA    Op = "NMLA"    // Negative Multiply Accumulate eg R0 = -R1 * R2 + R3
	NMLS    Op = "NMLS"    // Negative Multiply Subtract eg R0 = -R1 * R2 - R3
	NMUL    Op = "NMUL"    // Negative Multiply eg R0 = -R1 * R2
	SQRT    Op = "SQRT"    // Square Root eg R0 = sqrt(R1)
	QDMULH  Op = "QDMULH"  // Saturating Doubling Multiply Returning High Half
	QRDMULH Op = "QRDMULH" // Saturating Rounding Doubling Multiply Returning High Half

	// Load/Store Operations
	LDR Op = "LDR" // Load register eg R0 = [R1]
	STR Op = "STR" // Store register eg [R0] = R1
	LD2 Op = "LD2" // Load 2 eg R0 = [R1, R2]
	LD4 Op = "LD4" // Load 4 eg R0 = [R1, R2, R3, R4]
	ST1 Op = "ST1" // Store 1 eg [R0] = R1
	ST2 Op = "ST2" // Store 2 eg [R0, R1] = R2
	ST4 Op = "ST4" // Store 4 eg [R0, R1, R2, R3] = R4

	// Reduction Operations
	MAXV  Op = "MAXV"  // Maximum Across Vector
	MINV  Op = "MINV"  // Minimum Across Vector
	MAXAV Op = "MAXAV" // Maximum Absolute Across Vector
	MINAV Op = "MINAV" // Minimum Absolute Across Vector

	// Conversion Operations
	CVT  Op = "CVT"  // Convert
	CVTA Op = "CVTA" // Convert with Round to Nearest with Ties to Away
	CVTN Op = "CVTN" // Convert with Round to Nearest with Ties to Even
	CVTP Op = "CVTP" // Convert with Round towards Plus Infinity
	CVTM Op = "CVTM" // Convert with Round towards Minus Infinity

	// Predication and Masking
	PST   Op = "PST"   // Predicate Stack Push
	PSEL  Op = "PSEL"  // Predicated Select
	PNOT  Op = "PNOT"  // Predicate NOT
	PRFOP Op = "PRFOP" // Prefetch Op

	// Matrix Ops
	MMLAV Op = "MMLAV" // Matrix Multiply and Accumulate
	MMLAR Op = "MMLAR" // Matrix Multiply and Reduce Add
	MMLSR Op = "MMLSR" // Matrix Multiply and Reduce Subtract

	// Advanced SIMD (NEON) instructions
	// Arithmetic Operations
	ABA    Op = "ABA"    // Absolute Difference and Accumulate
	ABAL   Op = "ABAL"   // Absolute Difference and Accumulate Long
	ABD    Op = "ABD"    // Absolute Difference
	ABDL   Op = "ABDL"   // Absolute Difference Long
	DUP    Op = "DUP"    // Duplicate
	CLE    Op = "CLE"    // Compare Less Than or Equal
	CLT    Op = "CLT"    // Compare Less Than
	CGE    Op = "CGE"    // Compare Greater Than or Equal
	CGT    Op = "CGT"    // Compare Greater Than
	CEQ    Op = "CEQ"    // Compare Equal
	EXT    Op = "EXT"    // Extract
	HADD   Op = "HADD"   // Halving Add
	HSUB   Op = "HSUB"   // Halving Subtract
	MLAL   Op = "MLAL"   // Multiply Accumulate Long
	MLSL   Op = "MLSL"   // Multiply Subtract Long
	PADAL  Op = "PADAL"  // Pairwise Add and Accumulate Long
	PADD   Op = "PADD"   // Pairwise Add
	PMAX   Op = "PMAX"   // Pairwise Maximum
	PMIN   Op = "PMIN"   // Pairwise Minimum
	QMOVN  Op = "QMOVN"  // Saturating Move and Narrow
	RADDHN Op = "RADDHN" // Rounding Add and Narrow
	RHADD  Op = "RHADD"  // Rounding Halving Add
	REV    Op = "REV"    // Reverse
	TBL    Op = "TBL"    // Table Lookup
	TBX    Op = "TBX"    // Table Extension
	TRN    Op = "TRN"    // Transpose
	UZP    Op = "UZP"    // Unzip
	ZIP    Op = "ZIP"    // Zip

	// Shift Operations
	SHL  Op = "SHL"  // Shift Left
	SHR  Op = "SHR"  // Shift Right
	SRA  Op = "SRA"  // Shift Right and Accumulate
	RSRA Op = "RSRA" // Rounding Shift Right and Accumulate
	SLI  Op = "SLI"  // Shift Left and Insert
	SRI  Op = "SRI"  // Shift Right and Insert
	RSHR Op = "RSHR" // Rounding Shift Right
	RSHL Op = "RSHL" // Rounding Shift Left

	// Debug instructions
	BKPT  Op = "BKPT"  // Breakpoint instruction
	DBG   Op = "DBG"   // Debug hint
	HLT   Op = "HLT"   // Halt instruction
	SEV   Op = "SEV"   // Send event
	WFE   Op = "WFE"   // Wait for event
	WFI   Op = "WFI"   // Wait for interrupt
	YIELD Op = "YIELD" // Yield hint to threading system

	// ITM (Instrumentation Trace Macrocell) instructions
	ITR Op = "ITR" // Insert trace record
	DWT Op = "DWT" // Data watchpoint and trace

	// Debug barriers and synchronization
	DSB Op = "DSB" // Data synchronization barrier
	DMB Op = "DMB" // Data memory barrier
	ISB Op = "ISB" // Instruction synchronization barrier

	// Breakpoint and watchpoint instructions
	BRK  Op = "BRK"  // Software breakpoint
	WPT  Op = "WPT"  // Set watchpoint
	DWPT Op = "DWPT" // Set data watchpoint
	IWPT Op = "IWPT" // Set instruction watchpoint

	// Performance monitoring
	PMU        Op = "PMU"        // Performance monitor unit operation
	PMCCNTR    Op = "PMCCNTR"    // Read cycle counter
	PMCR       Op = "PMCR"       // Performance monitor control register
	PMCNTENSET Op = "PMCNTENSET" // Counter enable set register
	PMCNTENCLR Op = "PMCNTENCLR" // Counter enable clear register
	PMOVSR     Op = "PMOVSR"     // Overflow status register
	PMINTENSET Op = "PMINTENSET" // Interrupt enable set register
	PMINTENCLR Op = "PMINTENCLR" // Interrupt enable clear register
	PMXEVTYPER Op = "PMXEVTYPER" // Event type select register
	PMXEVCNTR  Op = "PMXEVCNTR"  // Event count register

	// CoreSight debug infrastructure
	CSETM  Op = "CSETM"  // CoreSight ETM control
	CSDWT  Op = "CSDWT"  // CoreSight DWT control
	CSITM  Op = "CSITM"  // CoreSight ITM control
	CSTPIU Op = "CSTPIU" // CoreSight TPIU control

	// Debug authentication
	DBGAUTHSTATUS Op = "DBGAUTHSTATUS" // Debug authentication status
	DBGBCR        Op = "DBGBCR"        // Breakpoint control register
	DBGBVR        Op = "DBGBVR"        // Breakpoint value register
	DBGWCR        Op = "DBGWCR"        // Watchpoint control register
	DBGWVR        Op = "DBGWVR"        // Watchpoint value register
	DBGDSCR       Op = "DBGDSCR"       // Debug status and control register

	// Debug communications
	DBGITR   Op = "DBGITR"   // Debug instruction transfer register
	DBGDTRRX Op = "DBGDTRRX" // Debug data transfer register receive
	DBGDTRTX Op = "DBGDTRTX" // Debug data transfer register transmit

	// Secure interrupt handling
	SPSR_MON Op = "SPSR_MON" // Saved Program Status Register (Monitor mode)
	LR_MON   Op = "LR_MON"   // Link Register (Monitor mode)
	GICD_S   Op = "GICD_S"   // Generic Interrupt Controller Distributor (Secure)
	GICC_S   Op = "GICC_S"   // Generic Interrupt Controller CPU Interface (Secure)

	// RME (Realm Management Extension)

	// Realm Entry/Exit
	ERET  Op = "ERET"  // Exit from realm
	RETAA Op = "RETAA" // Return from realm with authentication
	RETAB Op = "RETAB" // Return from realm with authentication variant B

	// Realm Management
	ERETA Op = "ERETA" // Exit from realm with authentication
	ERETB Op = "ERETB" // Exit from realm with authentication variant B

	// Realm State Management
	RPAA Op = "RPAA" // Return to protected level with authentication
	RPAB Op = "RPAB" // Return to protected level with authentication variant B
	DRPS Op = "DRPS" // Direct request to protected state
	MRS  Op = "MRS"  // Move to realm state
	MSR  Op = "MSR"  // Move from realm state

	RMEM Op = "RMEM" // Realm Memory Extension
	RMEN Op = "RMEN" // Realm Memory Enable

	// Realm Protection
	TLBI  Op = "TLBI"  // TLB invalidate by realm (cache maintenance)
	PACDA Op = "PACDA" // Pointer authentication for realm
	AUTDA Op = "AUTDA" // Authenticate pointer for realm

	// ARMv8 Cryptographic Extensions
	AES       Op = "AES"       // Advanced Encryption Standard instructions
	AESD      Op = "AESD"      // AES single round decryption
	AESE      Op = "AESE"      // AES single round encryption
	AESIMC    Op = "AESIMC"    // AES inverse mix columns
	AESMC     Op = "AESMC"     // AES mix columns
	SHA1      Op = "SHA1"      // Secure Hash Algorithm 1 instructions
	SHA1C     Op = "SHA1C"     // SHA1 hash update (choose)
	SHA1H     Op = "SHA1H"     // SHA1 fixed rotate
	SHA1M     Op = "SHA1M"     // SHA1 hash update (majority)
	SHA1P     Op = "SHA1P"     // SHA1 hash update (parity)
	SHA1SU0   Op = "SHA1SU0"   // SHA1 schedule update 0
	SHA1SU1   Op = "SHA1SU1"   // SHA1 schedule update 1
	SHA256    Op = "SHA256"    // SHA256 hash instructions
	SHA256H   Op = "SHA256H"   // SHA256 hash update (part 1)
	SHA256H2  Op = "SHA256H2"  // SHA256 hash update (part 2)
	SHA256SU0 Op = "SHA256SU0" // SHA256 schedule update 0
	SHA256SU1 Op = "SHA256SU1" // SHA256 schedule update 1
	PMULL     Op = "PMULL"     // Polynomial multiply long

	// ARMv8.2 Half-precision Floating-point
	FCVT Op = "FCVT" // floating-point convert
	FMAL Op = "FMAL" // floating-point multiply-accumulate long
	FMSL Op = "FMSL" // floating-point multiply-subtract long

	// ARMv8.3 Complex Number Extensions
	CADD Op = "CADD" // complex add
	CMLA Op = "CMLA" // complex multiply-accumulate
	CMUL Op = "CMUL" // complex multiply

	// ARMv8.3 JavaScript Conversion
	JCVT Op = "JCVT" // JavaScript convert

	// ARMv8.4 Dot Product
	DOT   Op = "DOT"   // dot product
	SDOT  Op = "SDOT"  // signed dot product
	UDOT  Op = "UDOT"  // unsigned dot product
	USDOT Op = "USDOT" // mixed sign dot product

	// ARMv8.5 Memory Tagging Extension
	STG   Op = "STG"   // Store Allocation Tag
	STZG  Op = "STZG"  // Store Zero and Allocation Tag
	ST2G  Op = "ST2G"  // Store Allocation Tag to Two Granules
	STZ2G Op = "STZ2G" // Store Zero and Allocation Tag to Two Granules
	LDG   Op = "LDG"   // Load Allocation Tag
	LDGM  Op = "LDGM"  // Load Allocation Tag Multiple
	STGM  Op = "STGM"  // Store Allocation Tag Multiple

	// ARMv8.6 Brain Float 16 (BFloat16)
	BFCVT   Op = "BFCVT"   // BFloat16 convert
	BFMLAL  Op = "BFMLAL"  // BFloat16 multiply-accumulate long
	BFMLALB Op = "BFMLALB" // BFloat16 multiply-accumulate long bottom
	BFMLALT Op = "BFMLALT" // BFloat16 multiply-accumulate long top

	// ARMv8.6 Enhanced Counter
	CNTP   Op = "CNTP"   // Counter-timer physical timer
	CNTV   Op = "CNTV"   // Counter-timer virtual timer
	CNTVCT Op = "CNTVCT" // Counter-timer virtual count

	// ARMv8.7 Pointer Authentication
	PACGA Op = "PACGA" // Pointer Authentication Code Generic
	XPAC  Op = "XPAC"  // Strip Pointer Authentication Code

	// SVE (Scalable Vector Extension)
	// Add across vector (SVE)
	ADDV Op = "ADDV" // eg R0 = R1 + R2
	// AND across vector (SVE)
	ANDV Op = "ANDV"
	// Break after first true condition (SVE)
	BRKB Op = "BRKB"
	// Break after last true condition (SVE)
	BRKA Op = "BRKA"
	// Complex dot product (SVE)
	CDOT Op = "CDOT"
	// Compare vectors for equality (SVE)
	CMPEQ Op = "CMPEQ"
	// Floating-point add accumulate (SVE)
	FADDA Op = "FADDA"
	// Floating-point fused multiply-add (SVE)
	FMLA Op = "FMLA"
	// Floating-point fused multiply-subtract (SVE)
	FMLS Op = "FMLS"

	// SME (Scalable Matrix Extension) Instructions
	// While less than (SME)
	WHILELT Op = "WHILELT"
	// Predicate true (SME)
	PTRUE Op = "PTRUE"
	// Predicate false (SME)
	PFALSE Op = "PFALSE"
	// Predicate next (SME)
	PNEXT Op = "PNEXT"
	// Predicate unpack high (SME)
	PUNPKHI Op = "PUNPKHI"
	// Predicate unpack low (SME)
	PUNPKLO Op = "PUNPKLO"

	// Matrix Multiply-Accumulate Operations
	// Signed matrix multiply-accumulate (SME)
	SMMLA Op = "SMMLA.S" // eg R0 = R1 * R2 + R3
	// Unsigned matrix multiply-accumulate (SME)
	UMMLA Op = "UMMLA.S" // eg R0 = R1 * R2 + R3
	// Mixed sign matrix multiply-accumulate (SME)
	USMMLA Op = "USMMLA" // eg R0 = R1 * R2 + R3
	// Signed matrix outer product and accumulate (SME)
	SMOPA Op = "SMOPA" // eg R0 = R1 * R2 + R3
	// Unsigned matrix outer product and accumulate (SME)
	UMOPA Op = "UMOPA" // eg R0 = R1 * R2 + R3
	// Signed-unsigned matrix outer product and accumulate (SME)
	SUMOPA Op = "SUMOPA" // eg R0 = R1 * R2 + R3
	// Unsigned-signed matrix outer product and accumulate (SME)
	USMOPA Op = "USMOPA" // eg R0 = R1 * R2 + R3
	// Zero matrix tile (SME)
	ZERO Op = "ZERO" // eg R0 = 0
	// Fill matrix tile (SME)
	FILL Op = "FILL" // eg R0 = 0x1234

	// Matrix Control Operations

	// Start streaming mode (SME)
	SMSTART Op = "SMSTART"
	// Stop streaming mode (SME)
	SMSTOP Op = "SMSTOP"
	// Read streaming vector length (SME)
	RDSVL Op = "RDSVL"
	// Add streaming vector length (SME)
	ADDSVL Op = "ADDSVL"
	// Reverse elements in matrix tile (SME)
	REVD Op = "REVD"

	// Matrix Tile Operations

	// Move to matrix tile
	MOVA Op = "MOVA"
	// Move from matrix tile
	MOVS Op = "MOVS"
	// Add horizontally across matrix tile
	ADDHA Op = "ADDHA"
	// Add vertically across matrix tile
	ADDVA Op = "ADDVA"
	// Signed clamp matrix tile elements
	SCLAMP Op = "SCLAMP"
	// Unsigned clamp matrix tile elements
	UCLAMP Op = "UCLAMP"

	// Trustzone instructions
	// Mode switching instructions
	// todo

	// TrustZone-specific coprocessor instructions

	// Move to System Register (Secure)
	MSR_S Op = "MSR_S"
	// Move from System Register (Secure)
	MRS_S Op = "MRS_S"

	// TrustZone cryptographic instructions (if implemented)
	// AES cryptographic operation (Secure)
	AES_S Op = "AES_S"
	// SHA cryptographic operation (Secure)
	SHA_S Op = "SHA_S"
	// Polynomial Multiply Long (Secure)
	PMULL_S Op = "PMULL_S"

	// TrustZone debug and trace
	// Debug Authentication Status (Secure)
	DBGAUTHSTATUS_S Op = "DBGAUTHSTATUS_S"
	// Breakpoint Control Register (Secure)
	DBGBCR_S Op = "DBGBCR_S"
	// Breakpoint Value Register (Secure)
	DBGBVR_S Op = "DBGBVR_S"
	// Watchpoint Control Register (Secure)
	DBGWCR_S Op = "DBGWCR_S"
	// Watchpoint Value Register (Secure)
	DBGWVR_S Op = "DBGWVR_S"
	// Debug Status and Control Register (Secure)
	DBGDSCR_S Op = "DBGDSCR_S"
)

func (op Op) String() string {
	return string(op)
}

func (op Op) IsBranch() bool {
	switch op {
	case B, BL, BEQ, BNE, BGT, BLT, BGE, BLE, CBZ, CBNZ, TBZ, TBNZ, BAL, BR, BLR, RET:
		return true
	default:
		return false
	}
}
