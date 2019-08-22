# Point to Point Links

## Supported modules

- Perfect Link
- Best Effort Broadcast

## How to build

Run the command below to generate an executable file named `p2plinks`.

`go build`

## How to run

The executable file generated from the building step above launches a process listening to the localhost port
corresponding to the processes ID which must be a number from `0` to the IP addresses file `number of lines - 1`.

`./p2plinks <id-number> <ip-addresses-file-name>`

If the first line in the IP addresses file is `127.0.0.1:5000` and the process ID is `0`, then the process is going to
listen to the port 5000 on the local machine.
