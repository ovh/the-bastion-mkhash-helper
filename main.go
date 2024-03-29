// vim: set ts=4 sw=4 sts=4 noet:
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
	"regexp"
	"runtime"
	"strings"
)

// CompileTime constants
var version = "undefined"
var date = "undefined"
var commit = "undefined"

// both salt and base64-encoded hash must only contain these chars
const allowedChars = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var validSalt = regexp.MustCompile(`^[` + allowedChars + `]+\z`)

// specific Encoding with a different alphabet ordering than StdEncoding,
// used to encode type8 and type9 passwords
var CiscoEncoding = base64.NewEncoding(allowedChars).WithPadding(base64.NoPadding)

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

func usage() {
	fmt.Fprintf(os.Stderr,
		`Usage: %s [OPTIONS]

  --salt-type8 SALT   Specify a fixed salt for type8
  --salt-type9 SALT   Specify a fixed salt for type9
  --version           Show version info

Type8 is a 256 bits PBKDF2-derived key with 16384 iterations using SHA256 for HMAC
Type9 is a 256 bits Scrypt-derived key with N=20000, r=1 and p=1

To hash a password, simply start the program, enter the password to hash and press ENTER.

Check the README for other ways to invoke the program in non-interactive ways with hints
to avoid leaking the password to tools such as 'ps' or your shell history.

`, os.Args[0])
}

func mustValidateSalt(salt string) string {
	if len(salt) < 4 || len(salt) > 32 {
		fmt.Fprintln(os.Stderr, "Salt length must be between 4 and 32 chars")
		os.Exit(-1)
	}
	if !validSalt.MatchString(salt) {
		fmt.Fprintf(os.Stderr, "Invalid salt, contains a char which is not part of the allowed chars list %q\n",
			allowedChars)
		os.Exit(-1)
	}
	return salt
}

func main() {
	// set default salts
	saltType8 := salt(14)
	saltType9 := salt(14)

	// parse the args
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--version" {
			fmt.Printf("Version %q (%s-%s) build at %s\n", version, commit, runtime.Version(), date)
			os.Exit(0)
		} else if os.Args[i] == "--salt-type8" && len(os.Args) >= i+2 {
			saltType8 = mustValidateSalt(os.Args[i+1])
			i++
		} else if os.Args[i] == "--salt-type9" && len(os.Args) >= i+2 {
			saltType9 = mustValidateSalt(os.Args[i+1])
			i++
		} else {
			usage()
			os.Exit(-1)
		}
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
			type8(password, saltType8),
			type9(password, saltType9),
			len(password),
		})

	if err != nil {
		panic(err)
	} else {
		fmt.Println(string(out))
	}
}
