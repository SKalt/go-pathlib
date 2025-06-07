https://godbolt.org/z/Kvx4YPoeW
```go
package main

import "fmt"

type Foo string

func (f Foo) Bar() string {
	return fmt.Sprintf("%s bar", string(f))
}
func foo(f string) string {
	return fmt.Sprintf("%s bar", f)
}

func main() {
	a := Foo("foo")
	fmt.Println(a.Bar())
	b := "not foo"
	fmt.Println(foo(b))
}
```
<details><summary>full output go ASM</summary>

```asm
main_Foo_Bar_pc0:
        TEXT    main.Foo.Bar(SB), ABIInternal, $36-16
        MOVW    8(g), R1
        PCDATA  $0, $-2
        CMP     R1, R13
        BLS     main_Foo_Bar_pc128
        PCDATA  $0, $-1
        MOVW.W  R14, -40(R13)
        FUNCDATA        $0, gclocals·Q1sFeS4dfNXE44R7gm5hqA==(SB)
        FUNCDATA        $1, gclocals·dIWD4RTVaLBTQr5cmlm6eA==(SB)
        FUNCDATA        $2, main.Foo.Bar.stkobj(SB)
        FUNCDATA        $5, main.Foo.Bar.arginfo1(SB)
        MOVW    $0, R0
        MOVW    R0, main..autotmp_2-8(SP)
        MOVW    $0, R0
        MOVW    R0, main..autotmp_2-4(SP)
        MOVW    main.f(FP), R0
        MOVW    R0, 4(R13)
        MOVW    main.f+4(FP), R0
        MOVW    R0, 8(R13)
        PCDATA  $1, $1
        CALL    runtime.convTstring(SB)
        MOVW    12(R13), R0
        MOVW    $type:string(SB), R1
        MOVW    R1, main..autotmp_2-8(SP)
        MOVW    R0, main..autotmp_2-4(SP)
        MOVW    $go:string."%s bar"(SB), R0
        MOVW    R0, 4(R13)
        MOVW    $6, R0
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_2-8(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $2
        CALL    fmt.Sprintf(SB)
        MOVW    24(R13), R0
        MOVW    28(R13), R1
        MOVW    R0, main.~r0+8(FP)
        MOVW    R1, main.~r0+12(FP)
        MOVW.P  40(R13), R15
main_Foo_Bar_pc128:
        NOP
        PCDATA  $1, $-1
        PCDATA  $0, $-2
        MOVW    R14, R3
        CALL    runtime.morestack_noctxt(SB)
        PCDATA  $0, $-1
        JMP     main_Foo_Bar_pc0
        JMP     0(PC)
        WORD    $type:string(SB)
        WORD    $go:string."%s bar"(SB)
main_foo_pc0:
        TEXT    main.foo(SB), ABIInternal, $36-16
        MOVW    8(g), R1
        PCDATA  $0, $-2
        CMP     R1, R13
        BLS     main_foo_pc128
        PCDATA  $0, $-1
        MOVW.W  R14, -40(R13)
        FUNCDATA        $0, gclocals·Q1sFeS4dfNXE44R7gm5hqA==(SB)
        FUNCDATA        $1, gclocals·dIWD4RTVaLBTQr5cmlm6eA==(SB)
        FUNCDATA        $2, main.foo.stkobj(SB)
        FUNCDATA        $5, main.foo.arginfo1(SB)
        MOVW    $0, R0
        MOVW    R0, main..autotmp_2-8(SP)
        MOVW    $0, R0
        MOVW    R0, main..autotmp_2-4(SP)
        MOVW    main.f(FP), R0
        MOVW    R0, 4(R13)
        MOVW    main.f+4(FP), R0
        MOVW    R0, 8(R13)
        PCDATA  $1, $1
        CALL    runtime.convTstring(SB)
        MOVW    12(R13), R0
        MOVW    $type:string(SB), R1
        MOVW    R1, main..autotmp_2-8(SP)
        MOVW    R0, main..autotmp_2-4(SP)
        MOVW    $go:string."%s bar"(SB), R0
        MOVW    R0, 4(R13)
        MOVW    $6, R0
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_2-8(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $2
        CALL    fmt.Sprintf(SB)
        MOVW    24(R13), R0
        MOVW    28(R13), R1
        MOVW    R0, main.~r0+8(FP)
        MOVW    R1, main.~r0+12(FP)
        MOVW.P  40(R13), R15
main_foo_pc128:
        NOP
        PCDATA  $1, $-1
        PCDATA  $0, $-2
        MOVW    R14, R3
        CALL    runtime.morestack_noctxt(SB)
        PCDATA  $0, $-1
        JMP     main_foo_pc0
        JMP     0(PC)
        WORD    $type:string(SB)
        WORD    $go:string."%s bar"(SB)
main_main_pc0:
        TEXT    main.main(SB), ABIInternal, $56-0
        MOVW    8(g), R1
        PCDATA  $0, $-2
        CMP     R1, R13
        BLS     main_main_pc324
        PCDATA  $0, $-1
        MOVW.W  R14, -60(R13)
        FUNCDATA        $0, gclocals·yr4yLQlPnRVKirpHM7hJfw==(SB)
        FUNCDATA        $1, gclocals·OKKa86U2/TXSbXuLfGsN8w==(SB)
        FUNCDATA        $2, main.main.stkobj(SB)
        NOP.EQ
        MOVW    $type:string(SB), R0
        MOVW    R0, main..autotmp_16-8(SP)
        MOVW    $main..stmp_0(SB), R0
        MOVW    R0, main..autotmp_16-4(SP)
        MOVW    $go:string."%s bar"(SB), R0
        MOVW    R0, 4(R13)
        MOVW    $6, R0
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_16-8(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $0
        CALL    fmt.Sprintf(SB)
        MOVW    24(R13), R0
        MOVW    28(R13), R1
        MOVW    $0, R2
        MOVW    R2, main..autotmp_23-16(SP)
        MOVW    $0, R2
        MOVW    R2, main..autotmp_23-12(SP)
        MOVW    R0, 4(R13)
        MOVW    R1, 8(R13)
        PCDATA  $1, $1
        CALL    runtime.convTstring(SB)
        MOVW    12(R13), R0
        MOVW    $type:string(SB), R1
        MOVW    R1, main..autotmp_23-16(SP)
        MOVW    R0, main..autotmp_23-12(SP)
        MOVW    os.Stdout(SB), R0
        NOP.EQ
        MOVW    $go:itab.*os.File,io.Writer(SB), R1
        MOVW    R1, 4(R13)
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_23-16(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $0
        CALL    fmt.Fprintln(SB)
        NOP.EQ
        MOVW    $type:string(SB), R0
        MOVW    R0, main..autotmp_16-8(SP)
        MOVW    $main..stmp_1(SB), R0
        MOVW    R0, main..autotmp_16-4(SP)
        MOVW    $go:string."%s bar"(SB), R0
        MOVW    R0, 4(R13)
        MOVW    $6, R0
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_16-8(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        CALL    fmt.Sprintf(SB)
        MOVW    24(R13), R0
        MOVW    28(R13), R1
        MOVW    $0, R2
        MOVW    R2, main..autotmp_28-24(SP)
        MOVW    $0, R2
        MOVW    R2, main..autotmp_28-20(SP)
        MOVW    R0, 4(R13)
        MOVW    R1, 8(R13)
        PCDATA  $1, $2
        CALL    runtime.convTstring(SB)
        MOVW    12(R13), R0
        MOVW    $type:string(SB), R1
        MOVW    R1, main..autotmp_28-24(SP)
        MOVW    R0, main..autotmp_28-20(SP)
        MOVW    os.Stdout(SB), R0
        NOP.EQ
        MOVW    $go:itab.*os.File,io.Writer(SB), R1
        MOVW    R1, 4(R13)
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_28-24(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $0
        CALL    fmt.Fprintln(SB)
        MOVW.P  60(R13), R15
main_main_pc324:
        NOP
        PCDATA  $1, $-1
        PCDATA  $0, $-2
        MOVW    R14, R3
        CALL    runtime.morestack_noctxt(SB)
        PCDATA  $0, $-1
        JMP     main_main_pc0
        JMP     0(PC)
        WORD    $type:string(SB)
        WORD    $main..stmp_0(SB)
        WORD    $go:string."%s bar"(SB)
        WORD    os.Stdout(SB)
        WORD    $go:itab.*os.File,io.Writer(SB)
        WORD    $main..stmp_1(SB)
main_Foo_Bar_pc0_1:
        TEXT    main.(*Foo).Bar(SB), DUPOK|WRAPPER|ABIInternal, $36-12
        MOVW    8(g), R1
        PCDATA  $0, $-2
        CMP     R1, R13
        BLS     main_Foo_Bar_pc160_1
        PCDATA  $0, $-1
        MOVW.W  R14, -40(R13)
        MOVW    16(g), R1
        CMP     $0, R1
        BNE     main_Foo_Bar_pc172_1
main_Foo_Bar_pc28_1:
        NOP
        FUNCDATA        $0, gclocals·nOkmq/HNV55TvJpJVI/Ozw==(SB)
        FUNCDATA        $1, gclocals·dIWD4RTVaLBTQr5cmlm6eA==(SB)
        FUNCDATA        $2, main.(*Foo).Bar.stkobj(SB)
        FUNCDATA        $5, main.(*Foo).Bar.arginfo1(SB)
        MOVW    main.f(FP), R0
        CMP     $0, R0
        BEQ     main_Foo_Bar_pc152_1
        MOVW    4(R0), R1
        MOVW    (R0), R0
        NOP.EQ
        MOVW    $0, R2
        MOVW    R2, main..autotmp_4-8(SP)
        MOVW    $0, R2
        MOVW    R2, main..autotmp_4-4(SP)
        MOVW    R0, 4(R13)
        MOVW    R1, 8(R13)
        PCDATA  $1, $1
        CALL    runtime.convTstring(SB)
        MOVW    12(R13), R0
        MOVW    $type:string(SB), R1
        MOVW    R1, main..autotmp_4-8(SP)
        MOVW    R0, main..autotmp_4-4(SP)
        MOVW    $go:string."%s bar"(SB), R0
        MOVW    R0, 4(R13)
        MOVW    $6, R0
        MOVW    R0, 8(R13)
        MOVW    $main..autotmp_4-8(SP), R0
        MOVW    R0, 12(R13)
        MOVW    $1, R0
        MOVW    R0, 16(R13)
        MOVW    R0, 20(R13)
        PCDATA  $1, $2
        CALL    fmt.Sprintf(SB)
        MOVW    24(R13), R0
        MOVW    28(R13), R1
        MOVW    R0, main.~r0+4(FP)
        MOVW    R1, main.~r0+8(FP)
        MOVW.P  40(R13), R15
main_Foo_Bar_pc152_1:
        CALL    runtime.panicwrap(SB)
        AND.EQ  R0, R0
main_Foo_Bar_pc160_1:
        NOP
        PCDATA  $1, $-1
        PCDATA  $0, $-2
        MOVW    R14, R3
        CALL    runtime.morestack_noctxt(SB)
        PCDATA  $0, $-1
        JMP     main_Foo_Bar_pc0_1
main_Foo_Bar_pc172_1:
        MOVW    (R1), R2
        ADD     $44, R13, R3
        CMP     R2, R3
        BNE     main_Foo_Bar_pc28_1
        ADD     $4, R13, R4
        MOVW    R4, (R1)
        JMP     main_Foo_Bar_pc28_1
        JMP     0(PC)
        WORD    $type:string(SB)
        WORD    $go:string."%s bar"(SB)
        TEXT    type:.eq.sync/atomic.Pointer[os.dirInfo](SB), DUPOK|LEAF|NOFRAME|ABIInternal, $-4-12
        FUNCDATA        $0, gclocals·TswRR9Pia9Wsluv5u1sUnA==(SB)
        FUNCDATA        $1, gclocals·J26BEvPExEQhJvjp9E8Whg==(SB)
        FUNCDATA        $5, type:.eq.sync/atomic.Pointer[os.dirInfo].arginfo1(SB)
        MOVW    main.q+4(FP), R0
        MOVW    (R0), R0
        MOVW    main.p(FP), R1
        MOVW    (R1), R1
        CMP     R1, R0
        MOVW    $0, R0
        MOVW.EQ $1, R0
        MOVB    R0, main.r+8(FP)
        JMP     (R14)
```

</details>

The go IR/assembly is pretty much the same for both the vanilla string and
```diff
--- main_Foo_Bar_pc0
+++ main_foo_pc0
@@ -1,15 +1,15 @@
--- ./main_Foo_Bar_pc0	2025-06-07 12:05:27.476501223 +0000
+++ ./main_foo_pc0	2025-06-07 12:05:49.188501701 +0000
@@ -1,15 +1,15 @@
-main_Foo_Bar_pc0:
-        TEXT    main.Foo.Bar(SB), ABIInternal, $36-16
+main_foo_pc0:
+        TEXT    main.foo(SB), ABIInternal, $36-16
         MOVW    8(g), R1
         PCDATA  $0, $-2
         CMP     R1, R13
-        BLS     main_Foo_Bar_pc128
+        BLS     main_foo_pc128
         PCDATA  $0, $-1
         MOVW.W  R14, -40(R13)
         FUNCDATA        $0, gclocals·Q1sFeS4dfNXE44R7gm5hqA==(SB)
         FUNCDATA        $1, gclocals·dIWD4RTVaLBTQr5cmlm6eA==(SB)
-        FUNCDATA        $2, main.Foo.Bar.stkobj(SB)
-        FUNCDATA        $5, main.Foo.Bar.arginfo1(SB)
+        FUNCDATA        $2, main.foo.stkobj(SB)
+        FUNCDATA        $5, main.foo.arginfo1(SB)
         MOVW    $0, R0
         MOVW    R0, main..autotmp_2-8(SP)
         MOVW    $0, R0
@@ -40,14 +40,14 @@
         MOVW    R0, main.~r0+8(FP)
         MOVW    R1, main.~r0+12(FP)
         MOVW.P  40(R13), R15
-main_Foo_Bar_pc128:
+main_foo_pc128:
         NOP
         PCDATA  $1, $-1
         PCDATA  $0, $-2
         MOVW    R14, R3
         CALL    runtime.morestack_noctxt(SB)
         PCDATA  $0, $-1
-        JMP     main_Foo_Bar_pc0
+        JMP     main_foo_pc0
         JMP     0(PC)
         WORD    $type:string(SB)
         WORD    $go:string."%s bar"(SB)

```


what happens when you use an inlinable function? Nothing much :shrug: https://godbolt.org/z/o5GovoqfK
https://godbolt.org/z/Ej6coM116 making Foo.Bar() private ass Foo.bar() helps inlining
