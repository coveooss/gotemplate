# Arithmetic

## Addition

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(1 + 2);              | {{ add 1 2 }}               | 3      | **Addition**
| @add(4, 5);            | {{ add 4 5 }}               | 9      | *or add*
| @sum(6,7);             | {{ sum 6 7 }}               | 13     | *or sum*
| @(2+3);                | {{ add 2 3 }}               | 5      | Spaces are optional
| @(  8  +  9  );        | {{ add 8 9 }}               | 17     | You can insert an arbitrary number of spaces in expressions
| @sum(1.2, 3.4);        | {{ sum 1.2 3.4 }}           | 4.6    | It also works with floating point numbers
| @sum(1, 2, 3, 4);      | {{ sum 1 2 3 4 }}           | 10     | It is possible to supply multiple arguments to addition operation
| @add(list(1,2,3,4));   | {{ add (list 1 2 3 4) }}    | 10     | this is useful on this line since there is ambiguity on where the expression finish

## Subtraction

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(4 - 2);              | {{ sub 4 2 }}               | 2      | **Subtraction**
| @sub(4, 2);            | {{ sub 4 2 }}               | 2      | *or sub*
| @subtract(4, 2);       | {{ subtract 4 2 }}          | 2      | *or subtract*

## Negative values

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(-23);                | {{ -23 }}                   | -23    | Negative value
| @(2 + -23);            | {{ add 2 -23 }}             | -21    | Operation with negative value
| @(2 + -(5 * 3));       | {{ add 2 (sub 0 (mul 5 3)) }} | -13  | Operation with negative expression

## Product

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(2 * 3);              | {{ mul 2 3 }}               | 6      | **Multiplication**
| @mul(4, 5);            | {{ mul 4 5 }}               | 20     | *or mul*
| @multiply(6, 7);       | {{ multiply 6 7 }}          | 42     | *or multiply*
| @prod(8, 9);           | {{ prod 8 9 }}              | 42     | *or prod*
| @product(10, 11);      | {{ product 10 11 }}         | 110    | *or product*
| @mul(1, 2, 3, 4);      | {{ mul 1 2 3 4 }}           | 24     | It is possible to supply multiple arguments to multiplication operation
| @mul(list(5,6,7,8));   | {{ mul (list 5 6 7 8) }}    | 1680   | or even an array

## Division

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(4 / 2);              | {{ div 4 2 }}               | 2      | **Division**
| @(13 ÷ 3);             | {{ div 13 3 }}              | 4.333333333333333 | *you can use the ÷ character instead of /*
| @div(20, 4);           | {{ div 20 4 }}              | 5      | *or div*
| @divide(10, 4);        | {{ divide 10 4 }}           | 2.5    | *or divide*
| @quotient(22, 10);     | {{ quotient 22 10 }}        | 2.2    | *or quotient*

## modulo

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(4 % 3);              | {{ mod 4 3 }}               | 1      | **Modulo**
| @mod(12, 5);           | {{ mod 12 5 }}              | 2      | *or mod*
| @modulo(20, 6)         | {{ modulo 20 6 }}           | 2      | *or modulo*

## Power

| Razor expression       | Go Template                 | Result | Note
| ---------------------- | --------------------------- | -----: | ----
| @(4 ** 3);             | {{ pow 4 3 }}               | 64     | **Power**
| @pow(12, 5);           | {{ pow 12 5 }}              | 248832 | *or pow*
| @power(3, 8);          | {{ power 3 8 }}             | 6561   | *or power*
| @pow10(3);             | {{ pow10 3 }}               | 1000   | **Power 10**
| @power10(5);           | {{ power10 5 }}             | 100000 | *or power10*
| @(1e+5);               | {{ 1e+5 }}                  | 100000 | Scientific notation (positive)
| @(2e-3);               | {{ 2e-3 }}                  | 0.002  | Scientific notation (negative)

## Bit operators

| Razor expression         | Go Template                 | Result | Note
| ------------------------ | --------------------------- | -----: | ----
| @(1 << 8);               | {{ lshift 1 8 }}            | 256    | **Left shift**
| @lshift(3, 5);           | {{ lshift 3 5 }}            | 96     | *or lshift*
| @leftShift(4, 4);        | {{ leftShift 4 4 }}         | 64     | *or leftShift*
| @(1024 >> 4);            | {{ rshift 1024 4 }}         | 64     | **Right shift**
| @rshift(456, 3);         | {{ rshift 456 3 }}          | 57     | *or rshift*
| @rightShift(72, 1);      | {{ rightShift 72 1 }}       | 36     | *or rightShift*
| @(65535 & 512);          | {{ band 65535 512 }}        | 512    | **Bitwise AND**
| @band(12345, 678);       | {{ band 12345 678 }}        | 32     | *or band*
| @bitwiseAND(222, 111);   | {{ bitwiseAND 222 111 }}    | 78     | *or bitwiseAND*
| @@(1 &#124; 2 &#124; 4); | {{ bor (bor 1 2) 4 }}       | 7      | **Bitwise OR**
| @bor(100, 200, 300);     | {{ bor 100 200 300 }}       | 492    | *or bor*
| @bitwiseOR(64, 256, 4);  | {{ bitwiseOR 64 256 4 }}    | 324    | *or bitwiseOR*
| @(1 ^ 2 ^ 4);            | {{ bxor (bxor 1 2) 4 }}     | 7      | **Bitwise XOR**
| @bxor(100, 200, 300);    | {{ bxor 100 200 300 }}      | 384    | *or bxor*
| @bitwiseXOR(64, 256, 4); | {{ bitwiseXOR 64 256 4 }}   | 324    | *or bitwiseXOR*
| @(255 &^ 4);             | {{ bclear 255 4 }}          | -      | **Bitwise Clear**
| @bclear(0xff, 3, 8);     | {{ bclear 0xff 3 8 }}       | -    | *or bclear*
| @bitwiseClear(0xf, 7);   | {{ bitwiseClear 0xf 7 }}    | -    | *or bitwiseClear*

## Other mathematic functions

### Special cases

There are special behavior for certain operators depending of the arguments:

#### String multiplication

@("*" * 100) will result in {{ mul "*" 100 }} which result in:

****************************************************************************************************
