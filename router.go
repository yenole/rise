package rise

import "github.com/yenole/rise/router"

var WrapGin = router.New

type (
	Router  = router.Router
	Context = router.Context
)
