/*
 * pixeldrain_test.go
 *
 * Copyright (c) 2018-2019 Junpei Kawamoto
 *
 * This software is released under the MIT License.
 *
 * http://opensource.org/licenses/mit-license.php
 */

package pixeldrain

import (
	"github.com/jkawamoto/go-pixeldrain/client"
	"os"
	"testing"
)

func TestNew(t *testing.T) {

	pd := New()
	if pd.Client != client.Default {
		t.Errorf("Client is %v, expects %v", pd.Client, client.Default)
	}
	if pd.Stderr != os.Stderr {
		t.Errorf("Stderr is %v, expects %v", pd.Stderr, os.Stderr)
	}

}
