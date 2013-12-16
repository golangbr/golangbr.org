// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

// This file replaces main.go when running godoc under app-engine.
// See README.godoc-app for details.

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func serveError(w http.ResponseWriter, r *http.Request, relpath string, err error) {
	w.WriteHeader(http.StatusNotFound)
	servePage(w, Page{
		Title:    "File " + relpath,
		Subtitle: relpath,
		Body:     applyTemplate(errorHTML, "errorHTML", err), // err may contain an absolute path!
	})
}

// Workaround for Github+GAE =P
func fakeZIP() (z *zip.Reader, err error) {
	var buf bytes.Buffer
	var wr io.Writer
	w := zip.NewWriter(&buf)

	filepath.Walk(zipGoroot, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			wr, err = w.Create(path)
			if err != nil {
				log.Fatalf("%s: %s\n", zipFilename, err)
			}
			ir, _ := ioutil.ReadFile(path)
			wr.Write(ir)
		}
		return nil
	})
	w.Close()
	byt := buf.Bytes()
	rdr := bytes.NewReader(byt)
	return zip.NewReader(rdr, int64(buf.Len()))
}

func init() {
	log.Println("initializing godoc ...")
	log.Printf(".zip file   = %s", zipFilename)
	log.Printf(".zip GOROOT = %s", zipGoroot)
	log.Printf("index files = %s", indexFilenames)

	// initialize flags for app engine
	*goroot = path.Join("/", zipGoroot) // fsHttp paths are relative to '/'
	*indexEnabled = false
	*indexFiles = indexFilenames
	*maxResults = 100    // reduce latency by limiting the number of fulltext search results
	*indexThrottle = 0.3 // in case *indexFiles is empty (and thus the indexer is run)
	*showPlayground = true

	// read .zip file and set up file systems
	rc, err := fakeZIP()
	if err != nil {
		log.Fatalf("%s: %s\n", zipFilename, err)
	}
	// rc is never closed (app running forever)
	fs.Bind("/", NewZipFS(rc, zipFilename), *goroot, bindReplace)
	// initialize http handlers
	readTemplates()
	initHandlers()
	registerPublicHandlers(http.DefaultServeMux)
	registerPlaygroundHandlers(http.DefaultServeMux)

	// initialize default directory tree with corresponding timestamp.
	initFSTree()

	// Immediately update metadata.
	updateMetadata()

	// initialize search index
	if *indexEnabled {
		go indexer()
	}

	log.Println("godoc initialization complete")
}
