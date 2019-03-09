package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/CyCoreSystems/go-meteor/star"
)

// Port indicates the port number on which Meteor should run.  It is optional,
// and if not defined, it will not be set within the Dockerfile
var Port string

var appDir string

var fullBase bool

// Config describes the configuration by which the Dockerfile should be generated
type Config struct {

	// MeteorVersion is the version of Meteor to use (for the build tool)
	MeteorVersion string

	// NodeVersion is the version of node to use for the base image
	NodeVersion string

	// SlimBase indicates that the base image for the final container image should be based off of the slim node base image rather than the full build.
	SlimBase bool

	// Port (if set) indicates what port should be EXPOSED and the value to which PORT will be set
	Port string
}

func init() {
	flag.StringVar(&appDir, "appDir", ".", "path to the base directory of the Meteor application")
	flag.StringVar(&Port, "port", "", "value of PORT variable to be set inside Dockerfile (optional)")
	flag.BoolVar(&fullBase, "full", false, "base the final image off the full node base image rather than the slim image")
}

func main() {
	flag.Parse()

	t := template.Must(template.New("Dockerfile").Parse(tmplString))

	s, err := star.Find(appDir)
	if err != nil {
		fmt.Printf("failed to locate star.json within %s: %v\n", appDir, err)
		os.Exit(1)
	}

	err = t.Execute(os.Stdout, &Config{
		MeteorVersion: s.MeteorVersion(),
		NodeVersion:   s.NodeVersion,
		SlimBase:      !fullBase,
		Port:          Port,
	})
	if err != nil {
		fmt.Println("failed to execute template:", err)
		os.Exit(1)
	}
	os.Exit(0)
}

var tmplString = `
FROM node:{{.NodeVersion}} AS builder

ENV BUNDLE_DIR /home/node/bundle
ENV SRC_DIR /home/node/src
ENV TMP_DIR /home/node/tmp

USER node:node

RUN mkdir -p $SRC_DIR $BUNDLE_DIR $TMP_DIR
COPY --chown=node:node . $SRC_DIR

RUN curl -o $TMP_DIR/meteor.sh 'https://install.meteor.com?release={{.MeteorVersion}}'; sh $TMP_DIR/meteor.sh

ENV PATH="/home/node/.meteor:${PATH}"
WORKDIR $SRC_DIR
RUN meteor npm install --production
RUN meteor build --server-only --directory $BUNDLE_DIR
RUN cd ${BUNDLE_DIR}/bundle/programs/server && npm install

FROM node:{{.NodeVersion}}{{if .SlimBase}}-slim{{end}}

ENV APP_DIR /home/node/app
ENV BUNDLE_DIR /home/node/bundle

USER node:node

COPY --from=builder $BUNDLE_DIR $APP_DIR
WORKDIR $APP_DIR/bundle

{{with .Port}}
ENV PORT {{.}}
EXPOSE {{.}}
{{end}}

CMD ["node", "./main.js"]
`
