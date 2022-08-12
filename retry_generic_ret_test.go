//   Copyright 2020 Vimeo
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

//go:build go1.18
// +build go1.18

package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTyped(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := make(chan struct{})
	backoff := DefaultBackoff()
	backoff.MinBackoff = time.Microsecond

	go func() {
		type retStruct struct {
			a int
			b string
		}

		q := 0
		r := NewRetryable(18)
		r.B = backoff
		s, err := Typed(ctx, r, func(ctx context.Context) (retStruct, error) {
			q++
			if q == 2 {
				return retStruct{a: 3, b: "fizzlebat"}, nil
			}
			return retStruct{}, fmt.Errorf("foo")
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, q)
		assert.Equal(t, retStruct{a: 3, b: "fizzlebat"}, s)
		close(c)
	}()
	<-c
}
