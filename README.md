# Processes-shell

## Overview

This project is a part of the **OSTEP course project [(Operating Systems: Three Easy Pieces)](https://pages.cs.wisc.edu/~remzi/Classes/537/Fall2021/)**, but it is implemented using Golang instead of C. The goal of this project is to build a simple Unix shell that mimics the functionality of basic Unix shells such as bash. The shell operates by reading commands, creating new processes to execute them, and managing the execution environment.

The shell operates in two modes:

1. **Interactive Mode:** The shell prints a prompt `(wish> )` and waits for user input.
2. **Batch Mode:** The shell reads commands from a file and executes them without a prompt.

Key features include process creation, redirection, parallel execution, and built-in commands such as `exit`, `cd`, and `path`.

You can find more details about the original project on the [OSTEP Projects](https://github.com/remzi-arpacidusseau/ostep-projects/tree/master/processes-shell).

## Main Functions
**1. Command Parsing**
- The shell reads input from the user (interactive mode) or from a file (batch mode).
- It splits the input into a command and its arguments.

**2. Process Execution**
- The shell executes external commands by creating a new process.
- Golang's os/exec package is used to handle process creation and command execution.

**3. Built-in Commands**
- **exit:** Terminates the shell.
- **cd:** Changes the current working directory.
- **path:** Sets the directories in which the shell looks for executables.

**4. Redirection**
- The shell supports output redirection using `>`. The output of a command can be redirected to a file by specifying the file after the `>` symbol.

**5. Parallel Commands**
- Users can execute multiple commands in parallel by separating them with `&`. The shell runs these commands simultaneously and waits for all of them to finish before continuing.

## Usage

### Running in Interactive Mode
To run the shell in interactive mode:
```shell
./wish
```
This will launch the shell and display the prompt `wish> ` for entering commands.

### Running in Batch Mode
To run the shell with a batch file:
```shell
./wish batch.txt
```
This reads commands from `batch.txt` and executes them.

## Testing

```shell
./test-wish.sh
```
![Screenshot from 2024-10-11 04-29-30](https://github.com/user-attachments/assets/e6c98210-1800-4a7d-a3fe-199f84cdb92e)

