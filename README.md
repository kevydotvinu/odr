## OVN kubernetes database reader - odr
A minimal tool that can read the OVN Kubernetes database file and print the output in JSON format. It will not come closer to `ovn-nbctl` commands but may help read the OVN-K database file and look for a specific piece of information fastly without any dependencies.

### Prerequisites
- OVN kubernetes database file called `leader_nbdb` in any sub-directory of `$(pwd)`.

### Installation

##### Linux
```
curl -#Lo odr $(curl -s https://api.github.com/repos/kevydotvinu/odr/releases/latest \
| jq -r '.assets | .[] | select(.name | contains("linux")) | .browser_download_url') && \
mv odr ~/.local/bin/ && chmod +x ~/.local/bin/odr && PATH=${PATH}:~/.local/bin
```

##### Mac
```
curl -#Lo odr $(curl -s https://api.github.com/repos/kevydotvinu/odr/releases/latest \
| jq -r '.assets | .[] | select(.name | contains("darwin")) | .browser_download_url') && \
mv odr ~/.local/bin/ && chmod +x ~/.local/bin/odr && PATH=${PATH}:~/.local/bin
```

### Usage

##### List all
```
odr --all
```

##### Search [table](https://man7.org/linux/man-pages/man5/ovn-nb.5.html#TABLE_SUMMARY)
```
odr --search Logical_Router
```

##### Show the details of a specific UUID
```
odr --uuid 571bb084-d5f2-4665-927e-03ace3a6b2b6
```
