## Installation

Inside the source directory, run the following command:

```sh
source install.sh zsh 
# or if you use bash run
source install.sh bash
```

It will build the binary file and put it in a bin folder on your `$HOME` dir.

## Usage

```bash
❯ grip deter .      

./test/br/col/sk/cool_essay.txt
7:'which makes the output <deter>ministic but means that for'
```

```
Usage:
  
 	grip <searchString> ( <searchDir> | . )

	- searchString	  The desired text you want to search for

	- searchDir   	  The directory in which you'd like to search for the specified text
	
	Alternatevely you can use [ . ] to search in the current directory

```