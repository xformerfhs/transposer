# Transposer

A program to encrypt text files by transposing characters.

## Introduction

Most ciphers are substitution ciphers.
They replace one character or group of characters with another one.
However, substitution ciphers can be attacked by statistical analysis.
If they are not done right, the count of characters or words can give away the key.
Modern ciphers like [ChaCha20](https://en.wikipedia.org/wiki/Salsa20#ChaCha_variant) or [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) cannot be attacked this way, but classical character substitution ciphers can.

There is another method of encryption where characters are not replaced, but their positions are changed.
These encryption methods are called "[transposition ciphers](https://en.wikipedia.org/wiki/Transposition_cipher)" or "transpositions".
Transpositions cannot be broken by statistical analysis.
In fact, they can only be attacked by brute-force or guessing of the keys.

This program implements a transposition cipher where the key is derived from one or more passwords.
In short, the characters of the clear text are filled into a table line by line and read out column by column where the order of the columns is determined by the passwords.
This is the easiest form of a transposition.
There are other ones.
This program implements the line-by-line to column-by-column transposition.

A secure transposition cipher needs multiple transpositions.
I.e., the transposed characters need to be transposed again (probably several times).
A single transposition leaves too much structure in the transposed characters.
So, this program allows for multiple passwords which means multiple transpositions.

For this to work securely, it is necessary that the lengths of the passwords must not have a common factor.
Here is an example:
Assume that the two passwords `abcdef` and `anotherone` are used.
The first one has a length of 6 characters, the second one has a length of 10 characters.
These lengths have a factor in common:
Both lengths are divisible by 2.
This common factor means that the combined transposition with both passwords also has this common factor in the transposed characters.
This weakens the encryption considerably.

Using passwords with lengths that have a common factor will produce an error message.

Transposition ciphers are only secure for long texts and long passwords.
This program enforces a minimum password length of 6 characters.
However, the length of the clear text is up to the user.

It should be at least a thousand characters, and the passwords should be long and random, like e.g. these two:
`fmejivhguqvlsd1esl5f2mbo0vowl5c3oh93deacosymu1vk0qes4064bqoemuiwxsz2c0sru97yzfmzg9i58jczg`
`wm2zmxrvz0idg7iox2k8fb9suzezcdqfrejjrj3ed2qzrl73c769eutl7ipihmlii3zdcfg7zamo3lt`

This program has been created for educational purposes.
I use it in my cryptography teachings to show how transposition works and what the strengths and weaknesses are.

## Calls

The program has four commands.
Two of them have options.

The commands are:

| Command   | Meaning                                     |
|-----------|---------------------------------------------|
| `decrypt` | Decrypt an encrypted text file.             |
| `encrypt` | Encrypt a text file.                        |
| `help`    | Show information about calling the program. |
| `version` | Show version information.                   |

When the program finishes, the return code is set according to the following table:

| Code | Meaning                   |
|------|---------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line |
| `2`  | Error while processing    |

### Encryption

A text file is encrypted like this:

```
transposer encrypt [-encoding {input[:output]}] {-passwords {transpositionPasswords} | -passwords-file {filename}] [-case {upper|lower}] [-onlyletters] {inputFilePath...}
```

The options have the following meaning:

| Option           | Meaning                                                                                                                                   |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| `encoding`       | Specifies the encoding of the input and the output file.                                                                                  |
| `passwords`      | A list of passwords separated by either `:` or `,`. Mutually exclusive with `passwords-file`.                                             |
| `passwords-file` | The path of a file that contains the passwords. Each line has one password. Blank lines are ignored. Mutually exclusive with `passwords`. |
| `case`           | The letters of the source file are converted to either lower- or upper-case when reading the source files.                                |
| `onlyletters`    | Only letters are read from the source files. All other characters are ignored.                                                            |

- One of the options `passwords` or `passwords-file` must be specified.
- The default encoding is `win1252` on Windows and `utf8` on all other operating systems.
- See below for an explanation of encodings.
- If only the input encoding is given, the output encoding is the same as the input encoding.
- If one of the encodings is empty, the default encoding is used.

Each option can be started by either `-` or `--`.

Encrypting produces a file which has `_trans` added before the file extension.

### Decryption

An encrypted text file is decrypted like this:

```
transposer decrypt [-encoding {input[:output]}] {-passwords {transpositionPasswords} | -passwords-file {filename}] {inputFilePath...}
```
I.e., it has the same options as the encryption except `case` and `onlyletters`.

The options have the following meaning:

| Option           | Meaning                                                                                                                                   |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| `encoding`       | Specifies the encoding of the input and the output file.                                                                                  |
| `passwords`      | A list of passwords separated by either `:` or `,`. Mutually exclusive with `passwords-file`.                                             |
| `passwords-file` | The path of a file that contains the passwords. Each line has one password. Blank lines are ignored. Mutually exclusive with `passwords`. |

- One of the options `passwords` or `passwords-file` must be specified.
- The default encoding is `win1252` on Windows and `utf8` on all other operating systems.
- See below for an explanation of encodings.
- If only the input encoding is given, the output encoding is the same as the input encoding.
- If one of the encodings is empty, the default encoding is used.

Each option can be started by either `-` or `--`.

Decrypting produces a file which has `_untrans` added before the file extension.

### Encodings

Text files are just a stream of bytes with no inherent meaning.
There has to be a mapping between bytes and characters.
I.e., each character is represented by a mapping into a sequence of bytes.
This is called a [character encoding](https://en.wikipedia.org/wiki/Character_encoding).
This encoding needs to be specified.
The list of valid file encodings is given when the program is started with no arguments or the `help` option..

Here are the most important ones:

| Name        | Meaning                                                          |
|-------------|------------------------------------------------------------------|
| `cp437`     | [IBM Code Page 437](https://en.wikipedia.org/wiki/Code_page_437) |
| `cp850`     | [IBM Code Page 850](https://en.wikipedia.org/wiki/Code_page_850) |
| `cp852`     | [IBM Code Page 852](https://en.wikipedia.org/wiki/Code_page_852) |
| `iso88591`  | [ISO 8859-1](https://en.wikipedia.org/wiki/ISO/IEC_8859-1)       |
| `iso885915` | [ISO 8859-15](https://en.wikipedia.org/wiki/ISO/IEC_8859-15)     |
| `utf16be`   | [UTF-16BE](https://en.wikipedia.org/wiki/UTF-16)                 |
| `utf16le`   | [UTF-16LE](https://en.wikipedia.org/wiki/UTF-16)                 |
| `utf8`      | [UTF-8](https://en.wikipedia.org/wiki/UTF-8)                     |
| `win1250`   | [Windows 1250](https://en.wikipedia.org/wiki/Windows-1250)       |
| `win1252`   | [Windows 1252](https://en.wikipedia.org/wiki/Windows-1252)       |

`utf16` is a synonym for `utf16le`.

## Examples

In all examples the source file has the name `Freude.txt`, encoded in UTF-8 with the following content:

```text
Freude, schöner Götterfunken,
Tochter aus Elisium,
Wir betreten feuertrunken
Himmlische, dein Heiligthum.
Deine Zauber binden wieder,
was der Mode Schwerd getheilt;
Bettler werden Fürstenbrüder,
wo dein sanfter Flügel weilt.
```

This is the first verse of Schiller's poem 'An die Freude' ('Ode to Joy').

In the first example this file is encrypted with the password `schiller` like this:

```
transposer encrypt -encoding utf8 -passwords schiller Freude.txt
```

The result is a file with the name `Freude_trans.txt` with the following content encoded in UTF-8:

```text
rcökh ,tekldl
and edllee,nr ,rfTai nrieHue wwMwheerddfgtehtetE
rueieiDuded  tenn
  wuötnelWeensigebereSg;r bwsFedne,riitr
cntien,rce
 Froalieer
 sretHh hnr 
 htBwüü nül  uouub um,em biaoeetrseete.FsGncsmefnm i.Ziesdritdtriel
```

While this is unintelligible, it is easy to break.
The clear text is quite short and the password has only 8 letters.

Decrypting the encrypted file with the following call recovers the original file:

```
transposer decrypt -encoding utf8 -passwords schiller Freude_trans.txt
```

The result is a file with the name `Freude_untrans.txt` with the original content encoded in UTF-8.

Now, the option `case` is added to the call:

```
transposer encrypt -encoding utf8 -passwords schiller -case upper Freude.txt
```

The result is a file with the name `Freude_trans.txt` with the following content encoded in UTF-8:

```text
RCÖKH ,TEKLDL
AND EDLLEE,NR ,RFTAI NRIEHUE WWMWHEERDDFGTEHTETE
RUEIEIDUDED  TENN
  WUÖTNELWEENSIGEBERESG;R BWSFEDNE,RIITR
CNTIEN,RCE
 FROALIEER
 SRETHH HNR 
 HTBWÜÜ NÜL  UOUUB UM,EM BIAOEETRSEETE.FSGNCSMEFNM I.ZIESDRITDTRIEL
```

When decrypting this file, the original file is restored, but with all characters in upper-case.

In the next example the option `onlyletters` is used:

```
transposer encrypt -encoding utf8 -passwords schiller -onlyletters Freude.txt
```

The result is a file with the name `Freude_trans.txt` with the following content encoded in UTF-8:

```text
röeTumerHhiDbnaegBesesllsökeseekinhZneMeieFreewenrosWttieleewsSeertraütuefcEiermdiiridcttdewngdruhlrnumegnbeehhtenofeeGntibfnliteidrwelnbdtlcterituesHuadrorlrüüireFhtnaurencemuewddtwrdnFi
```

There are no longer any line-breaks, blanks or punctuations.

This file can be decrypted with the following call:

```
transposer decrypt -encoding utf8 -passwords schiller Freude_trans.txt
```

The result is an UTF-8 encoded file `Freude_untrans.txt` with the following content:

```text
FreudeschönerGötterfunkenTochterausElisiumWirbetretenfeuertrunkenHimmlischedeinHeiligthumDeineZauberbindenwiederwasderModeSchwerdgetheiltBettlerwerdenFürstenbrüderwodeinsanfterFlügelweilt
```

The decrypted file also has no blanks, punctuation or line-breaks, just letters.

In the next example two passwords are specified:

```
transposer encrypt -encoding utf8 -passwords schiller,beethoven Freude.txt
```

The encrypted file `Freude_trans.txt` has the following content:
```text
rke rMgre Wewier nüomrG icld,iwtud ersinosrüu snitödlrehee weeFt,ar  ubec.dtdn wdEunee ,t r B ,eFndlh
eTuetetösgd
cit ü atmire rnwf
d
lbbriF
hwuetsmr
 aeaereieti;nceeHhluoeeeikllfHehi unSerrle
nbiesZt,n,i dtDnngren
eht me.fse
```

This is harder to break than using only one password, but the clear text is still too short.

In the following example still another password is used:

```
transposer encrypt -encoding utf8 -passwords schiller,beethoven,thalia Freude.txt
```

This leads to the following two error messages:

```text
2025-04-20 15:05:46 +02:00  13  E  Password error: the lengths of the passwords 'schiller' (8) and 'thalia' (6) share a common divisor: 2
2025-04-20 15:05:46 +02:00  13  E  Password error: the lengths of the passwords 'beethoven' (9) and 'thalia' (6) share a common divisor: 3
```

## Important Notes

- Keep your passwords secure — they are required for decryption.
- When using multiple passwords, remember their order.
- The letters-only option is irreversible — removed characters cannot be restored.
- Make sure to use the same options during decryption as used during encryption.

## Security Considerations

- The security strength depends on the complexity and length of the passwords.
- Using multiple passwords increases security.

## Challenge

There is a file `Challenge_trans.txt` in the `challenge` directory.
This file has been created by encrypting a German text encoded in UTF-8 with three passwords.
The passwords were not random and none of them is longer than 29 characters.

If anyone is able to decrypt this file and tell me its correct clear text, I would be very much interested in how this has been done.
Guessing the clear text does not count.
There has to be an algorithm that does only use the file and the information given about the passwords.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
