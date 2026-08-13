// Package docs embeds the hand-written OpenAPI specification so it ships
// inside the compiled binary (and the Docker image) without depending on
// the working directory the process is started from.
package docs

import _ "embed"

//go:embed openapi.json
var OpenAPISpec []byte
