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
package progress

import (
	"time"

	"github.com/moby/buildkit/client"
)

func ResetTime(in Writer) Writer {
	w := &pw{Writer: in, status: make(chan *client.SolveStatus), tm: time.Now()}
	go func() {
		for {
			select {
			case <-in.Done():
				return
			case st, ok := <-w.status:
				if !ok {
					close(in.Status())
					return
				}
				if w.diff == nil {
					for _, v := range st.Vertexes {
						if v.Started != nil {
							d := v.Started.Sub(w.tm)
							w.diff = &d
						}
					}
				}
				if w.diff != nil {
					for _, v := range st.Vertexes {
						if v.Started != nil {
							d := v.Started.Add(-*w.diff)
							v.Started = &d
						}
						if v.Completed != nil {
							d := v.Completed.Add(-*w.diff)
							v.Completed = &d
						}
					}
					for _, v := range st.Statuses {
						if v.Started != nil {
							d := v.Started.Add(-*w.diff)
							v.Started = &d
						}
						if v.Completed != nil {
							d := v.Completed.Add(-*w.diff)
							v.Completed = &d
						}
						v.Timestamp = v.Timestamp.Add(-*w.diff)
					}
					for _, v := range st.Logs {
						v.Timestamp = v.Timestamp.Add(-*w.diff)
					}
				}
				in.Status() <- st
			}
		}
	}()
	return w
}

type pw struct {
	Writer
	tm     time.Time
	diff   *time.Duration
	status chan *client.SolveStatus
}

func (p *pw) Status() chan *client.SolveStatus {
	return p.status
}
