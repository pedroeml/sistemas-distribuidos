# Point to Point Links

## Supported modules

- Perfect Link
- Best Effort Broadcast

## How to build

This project was developed with Go version 1.12.9. If you have installed this version or higher of Go, just run the command below inside the `p2plinks` folder to generate an executable file named `p2plinks`.

```bash
$ go build
```

If you would like to install an extra version of Go, you can install it as follows:

```bash
$ go get golang.org/dl/go1.12.9
$ go1.12.9 download
```

Then, inside the `p2plinks` folder, build this project with:

```bash
$ go1.12.9 build
```

## How to run

The executable file generated from the building step above launches a process listening to the localhost port
corresponding to the processes ID which must be a number from `0` to the IP addresses file `number of lines - 1`.

```bash
$ ./p2plinks <id-number> <ip-addresses-file-name>
```

If the first line in the IP addresses file is `127.0.0.1:5000` and the process ID is `0`, then the process is going to
listen to the port 5000 on the local machine.
