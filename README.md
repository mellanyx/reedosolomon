<!-- markdownlint-configure-file {
  "MD013": {
    "code_blocks": false,
    "tables": false
  },
  "MD033": false,
  "MD041": false
} -->

<div align="center">

# universal encoder/decoder reedosolomon

<img alt="GitHub (Pre-)Release Date" src="https://img.shields.io/github/release-date-pre/mellanyx/reedosolomon">
<img src="https://img.shields.io/github/v/release/mellanyx/reedosolomon">
<img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/mellanyx/reedosolomon">

####

<img alt="GitHub all releases" src="https://img.shields.io/github/downloads/mellanyx/reedosolomon/total">
<img src="https://img.shields.io/github/go-mod/go-version/mellanyx/reedosolomon">

**This library is a universal file encoder and decoder using Ridosolomon codes with error correction**

It remembers which directories you use most frequently, so you can "jump" to
them in just a few keystrokes.<br />
zoxide works on all major shells.

[Getting started](#getting-started) •
[Installation](#installation) •
[Configuration](#configuration) •
[Integrations](#third-party-integrations)

</div>

# Use as compiled file

### *Step 1: Clone repo*

```sh
git clone https://github.com/mellanyx/reedosolomon.git
```

### *Step 2: Build project*

Let's go to the cloned project and execute the following command:

```sh
go build
```
After that the compiled file will appear

![](assetsForReadMe/gifs/build.gif)


### *Step 3: Applications*

#### *Encode file*

```sh
./reedosolomon encode [file path] [primitive] [eccsybmols]
```

primitive poly GF(256) = one of (285, 299, 301, 333, 351, 355, 357, 361, 369, 391, 397, 425, 251, 463, 487, 501)

Number of bugs to be fixed = eccsybmols * 2

For example, if we specify eccsybmols equal to 2, then the maximum number of errors that we can fix will be 4.

![](assetsForReadMe/gifs/encode.gif)


#### *Corrupt file*


```sh
./reedosolomon corrupt [file path] [eccsybmols]
```

eccsybmols - we indicate a value no more than what we indicated during encoding. Once executed, the program itself will corrupt every 256 bits with an error count equal to (eccsybmols / 2)

![](assetsForReadMe/gifs/corrupt.gif)


#### *Decode file*

```sh
./reedosolomon decode [file path] [primitive poligon 285 / 301] [eccsybmols]
```

primitive poligon - we specify the same value as when encoding.

eccsybmols - we specify the same value as when encoding and corrupt.

After execution, we should get our original fully working file

![](assetsForReadMe/gifs/decode.gif)



# Use as a library

### *Step 1: Import in your project*

```sh
go get https://github.com/mellanyx/reedosolomon.git
```

```golang
package main

import "github.com/mellanyx/reedosolomon/functions"
```

### *Step 2: Usage*

#### *Encode file*

primitive poly GF(256) = one of (285, 299, 301, 333, 351, 355, 357, 361, 369, 391, 397, 425, 251, 463, 487, 501)

Number of bugs to be fixed = eccsybmols * 2

For example, if we specify eccsybmols equal to 2, then the maximum number of errors that we can fix will be 4.

```golang
reedosolomon.EncodeFile(filepath, primitive, eccsymbols)
```

#### *Corrupt file*

eccsybmols - we indicate a value no more than what we indicated during encoding. Once executed, the program itself will corrupt every 256 bits with an error count equal to (eccsybmols / 2)

```golang
reedosolomon.CorruptFile(filepath, eccsymbols)
```

#### *Decode file*

```golang
reedosolomon.DecodeAndFixCorruptFile(filepath, primitive, eccsymbols)
```

primitive poligon - we specify the same value as when encoding.

eccsybmols - we specify the same value as when encoding and corrupt.

After execution, we should get our original fully working file

# Testing on a 1GB file

### *Tested on MacBook Air 2020:*
* **Chip: Apple M1**
* **Shared memory: 16 GB**
* **SSD: 512 GB**

#### *Encode 1GB file*

![](assetsForReadMe/screenshots/encode_1GB.png)

#### *Corrupt & Decode 1GB file*

![](assetsForReadMe/screenshots/corrupt_decode_1GB.png)

#### *Result outputs*

![](assetsForReadMe/screenshots/results_1GB.png)

## Authors

* **Alexandr Starikov** - [GitHub](https://github.com/mellanyx)

## License

This project is licensed under the MIT License.