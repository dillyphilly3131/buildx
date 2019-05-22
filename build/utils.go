/*
   Copyright The buildx Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package build

import (
	"archive/tar"
	"bytes"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// archiveHeaderSize is the number of bytes in an archive header
const archiveHeaderSize = 512

func isLocalDir(c string) bool {
	st, err := os.Stat(c)
	return err == nil && st.IsDir()
}

func isArchive(header []byte) bool {
	for _, m := range [][]byte{
		{0x42, 0x5A, 0x68},                   // bzip2
		{0x1F, 0x8B, 0x08},                   // gzip
		{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}, // xz
	} {
		if len(header) < len(m) {
			continue
		}
		if bytes.Equal(m, header[:len(m)]) {
			return true
		}
	}

	r := tar.NewReader(bytes.NewBuffer(header))
	_, err := r.Next()
	return err == nil
}

// toBuildkitExtraHosts converts hosts from docker key:value format to buildkit's csv format
func toBuildkitExtraHosts(inp []string) (string, error) {
	if len(inp) == 0 {
		return "", nil
	}
	hosts := make([]string, 0, len(inp))
	for _, h := range inp {
		parts := strings.Split(h, ":")

		if len(parts) != 2 || parts[0] == "" || net.ParseIP(parts[1]) == nil {
			return "", errors.Errorf("invalid host %s", h)
		}
		hosts = append(hosts, parts[0]+"="+parts[1])
	}
	return strings.Join(hosts, ","), nil
}
