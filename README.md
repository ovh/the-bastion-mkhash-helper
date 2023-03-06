# The Bastion mkhash helper

This helper is an optional companion to [The Bastion](https://github.com/ovh/the-bastion),
it generate hashes from a password specified on STDIN, and outputs the result in JSON on STDOUT. 

Currently, two hashes formats are returned, commonly used on network devices:

- So-called "type 8" passwords, where a key is derived using PBKDF2-HMAC-SHA256 with 20000 rotations
- So-called "type 9" passwords, where a key is derived using scrypt with 16384 rotations

# Installation

If you're using [The Bastion](https://github.com/ovh/the-bastion), this helper
is installed automatically for you as part of the standard install procedure.

If you want to install it manually, `deb` and `rpm` packages are provided, pick the proper one for your distro.

Alternatively, static binaries are also provided, so you can simply drop them in the proper
location depending on your OS, most of the time `/usr/local/bin` is a good pick.

# Howto's

**This helper doesn't expect to be called manually.**

However if you insist doing it, take extra-care avoiding the revelation of
your password to other parties that might be present on the system.
This is why passing the password as a command-line parameter will never be implemented,
as it would like it through ``ps`` or similar tools, for the fraction of a second
the program might be running.

Using ``bash``, for example, usually any command starting with a ``space`` will not be
saved in the history, as long as ``ignorespace`` or ``ignoreboth`` is present in your
``HISTCONTROL`` environment variable, which is usually the default.
Using the ``<<<`` bashism will use the following string as the called command's STDIN,
without leaking anything to ``ps``.

A lot of this depends on the specificities of your system, which are out of the
scope of this readme.

# Example

```bash
$ echo $HISTCONTROL
ignoreboth
$  ./the-bastion-mkhash-helper <<< "my complex password" | jq .
{
  "Type8": "$8$WdC8ZM7L0COoOO$1u2qMXybH90nzp2QADbJrMLWQL1Dk6n3pkEstOi4NFs",
  "Type9": "$9$94iBImyN51ac3p$SzhA.gzdov6j6vjj8MpMbhMWRLbF.8xv5GHx8udb.62",
  "PasswordLen": 19
}
```

# License

Copyright 2023 OVH SAS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
