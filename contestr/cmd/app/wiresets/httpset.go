package wiresets

import (
	"contestr/internal/handlers/http"
	"contestr/internal/services/codeforces"
	"contestr/internal/transport"
	"github.com/google/wire"
)

var handlesSet = wire.NewSet(
	http.NewHelloHandle,

	codeforces.NewService,
	wire.Bind(new(http.CodeforcesService), new(*codeforces.Service)),
	http.NewContestHandle,
)

var HTTPSet = wire.NewSet(
	handlesSet,
	http.NewHandlers,
	transport.NewHTTPServer,
)
