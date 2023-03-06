/*
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this
 * file except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
 * ANY KIND, either express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 */

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"io"
	"os"
	"runtime"
	"strings"
)

// CompileTime constants
var version = "undefined"
var date = "undefined"
var commit = "undefined"

// specific Encoding with a different alphabet ordering than StdEncoding,
// used to encode type8 and type9 passwords
var CiscoEncoding = base64.NewEncoding("./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz").WithPadding(base64.NoPadding)

// type8 and type9 crypto parameters taken from https://github.com/openwall/john/issues/711
func type8(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 20000, 32, sha256.New)
	return "$8$" + salt + "$" + CiscoEncoding.EncodeToString(dk)
}

func type9(password string, salt string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 1, 1, 32)
	return "$9$" + salt + "$" + CiscoEncoding.EncodeToString(dk)
}

func salt(size int) string {
	if size < 1 {
		panic("salt size must be > 0")
	}
	// size is the expected encoded size, so we need to generate a smaller
	// amount of random data to get exactly this number of encoded chars in the end,
	// however it might be impossible to get exactly this size due to how base64
	// works, so if this is the case, generate a bit too much data then truncate
	// the string to the target size after it's encoded
	randSize := CiscoEncoding.DecodedLen(size)
	for CiscoEncoding.EncodedLen(randSize) < size {
		randSize++
	}

	buf := make([]byte, randSize)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	// we might have generated too long a string, slice it to the requested size
	return CiscoEncoding.EncodeToString(buf)[:size]
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Version %q (%s-%s) build at %s\n", version, commit, runtime.Version(), date)
		os.Exit(0)
	}

	if len(os.Args) != 1 {
		fmt.Fprintf(os.Stderr,
			`Usage: %s [--version]

To hash a password, simply push the password to STDIN, taking care
to not leak it to tools such as 'ps', e.g. using bash:

  %s <<< "my password"

or

  echo "my password" | %s

as 'echo' is usually a shell builtin, hence not appearing in 'ps'

`, os.Args[0], os.Args[0], os.Args[0])
		os.Exit(-1)
	}

	// read the password and trim the ending \n
	password, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}
	password = strings.TrimSuffix(password, "\n")

	// generate the hashes and print them
	out, err := json.Marshal(
		struct {
			Type8       string
			Type9       string
			PasswordLen int
		}{
			type8(password, salt(14)),
			type9(password, salt(14)),
			len(password),
		})

	if err != nil {
		panic(err)
	} else {
		fmt.Println(string(out))
	}
}
