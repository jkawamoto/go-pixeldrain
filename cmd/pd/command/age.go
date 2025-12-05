// age.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"filippo.io/age"
)

const AgeExt = ".age"

// waitReadCloser is an io.ReadCloser that calls the given wait function when it's closed.
type waitReadCloser struct {
	io.ReadCloser
	wait func()
}

func (r waitReadCloser) Close() error {
	err := r.ReadCloser.Close()
	r.wait()
	return err
}

// waitWriteCloser is an io.WriteCloser that calls the given wait function when it's closed.
type waitWriteCloser struct {
	io.WriteCloser
	wait func()
}

func (w waitWriteCloser) Close() error {
	err := w.WriteCloser.Close()
	w.wait()
	return err
}

func Encrypt(src io.ReadCloser, recipients []age.Recipient) io.ReadCloser {
	r, w := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		defer func() {
			if reason := recover(); reason != nil {
				err = errors.Join(err, fmt.Errorf("recovered: %v", reason))
			}
			_ = w.CloseWithError(errors.Join(err, src.Close()))
		}()

		ew, err := age.Encrypt(w, recipients...)
		if err != nil {
			return
		}
		defer func() {
			err = errors.Join(err, ew.Close())
		}()

		_, err = io.Copy(ew, src)
	}()

	return waitReadCloser{
		ReadCloser: r,
		wait:       wg.Wait,
	}
}

func Decrypt(dest io.WriteCloser, identity []age.Identity) io.WriteCloser {
	r, w := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		defer func() {
			if reason := recover(); reason != nil {
				err = errors.Join(err, fmt.Errorf("recovered: %v", reason))
			}
			_ = r.CloseWithError(errors.Join(err, dest.Close()))
		}()

		dr, err := age.Decrypt(r, identity...)
		if err != nil {
			return
		}

		_, err = io.Copy(dest, dr)
	}()

	return waitWriteCloser{
		WriteCloser: w,
		wait:        wg.Wait,
	}
}
