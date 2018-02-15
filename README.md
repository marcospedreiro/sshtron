# [sshtron](https://github.com/zachlatta/sshtron)

**Note:** All the credit goes to the authors of the original [zachlatta/sshtron](https://github.com/zachlatta/sshtron)!

I just rewrote the code in a different format as a learning exercise because I've never written any game code.

Below is a slightly modified version of the original README

----

SSHTron is a multiplayer lightcycle game that runs through SSH. Just run the command below and you'll be playing in seconds:

```bash
$ ssh sshtron.zachlatta.com
```

## Controls

- WASD or vim keybindings to move (do not use your arrow keys).
- `]` to accelerate, `[` to decelerate
- Escape or Ctrl+C to exit.

## Want to choose color yourself?
There are total 7 colors to choose from: Red, Green, Yellow, Blue, Magenta, Cyan and White

```bash
$ ssh red@sshtron.zachlatta.com
```

If the color you picked is already taken in all open games, you'll randomly be assigned a color.

## Running Your Own Copy

Clone the project and cd into its directory. These instructions assume that you have your GOPATH setup correctly.

```bash
# Create an RSA public/private keypair in the current directory for the server
# to use. Don't give it a passphrase.
$ ssh-keygen -t rsa -f sshtron.pem

# Download dependencies and compile the project
$ go get && make build

# Run it! Configuration is driven out of the ./config/resources/config.json by default!
$ make run
```

## Running under a Docker container

See the original [zachlatta/sshtron](https://github.com/zachlatta/sshtron)!

## [CVE-2016-0777](https://www.qualys.com/2016/01/14/cve-2016-0777-cve-2016-0778/openssh-cve-2016-0777-cve-2016-0778.txt)

CVE-2016-0777 revealed two SSH client vulnerabilities that can be exploited by a malicious SSH server. While SSHTron does not exploit these vulnerabilities, you should still patch your client before you play. SSHTron is open source, but the server could always be running a modified version of SSHTron that does exploit the vulnerabilities described in CVE-2016-0777.

[If you haven't yet patched your SSH client, you can follow these instructions to do so now.](https://www.jacobtomlinson.co.uk/quick%20tip/2016/01/15/fixing-ssh-vulnerability-CVE-2016-0777/)

## License

SSHTron is licensed under the MIT License. See the full license text in LICENSE.
