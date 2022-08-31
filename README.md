## Installation

Inside the source directory, run the following command:

```sh
source install.sh zsh 
# or if you use bash run
source install.sh bash
```
It will build the binary file and put it in a bin folder on your `$HOME` dir.

---

**Note**: If you have already run the command and your *.rc file is already modified, you can run the command again but without the zsh/bash part in order to not have the function repeated in it. 


## Usage

```bash
❯ grip deter .      

./test/br/col/sk/cool_essay.txt
7:'which makes the output <deter>ministic but means that for'
```

```
Usage:

grip <searchString> ( <searchDir> | . ) [-opt]
	
Arguments:

	<searchString>	  The desired text you want to search for

	<searchDir>   	  The directory in which you'd like to search. Use '.' to search in the current directory

Options:
	
	-h 			  Show hidden folders and files

```