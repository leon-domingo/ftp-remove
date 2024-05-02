# A cli-command to remove "old" files from a FTP server

## Usage

`ftp-remove <host:port> <user> <password> <max_age_in_days> <file_regex>`

For example,

`ftp-remove ftp.mydomain.com:21 us3r p4sSw0rD 3 "foobar__(?P<year>\d{4})(?P<month>\d{4})(?P<day>\d{2}).tgz"`

This will remove files named *"foobar__YYYYmmdd.tgz"* older than **3 days** in the FTP **ftp.mydomain.com** (**port=21**) whose access credentials are **us3r** / **p4sSw0rD**.

You may notice the *file regex* follows the Python / Go syntax for the **named groups**. You MUST define a *year*, *month* and *day* groups in order to calculate the "age" of the file.
