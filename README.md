# key-rotation

Abstraction + bindings for gracefully rotating sets of keys.

The **Interesting Parts**:
* A [KeyStore](rotation/store.go) represents the target system and user for which a set of keys is to be managed.
* [KeyRotation](rotation/rotation.go) uses a two phase plan & apply approach.  Plan will decide the state of the current
  keys ( states enumerated below ) and apply will make those changes happen.
  * **Valid** - A key is younger than the start of the grace period.  Use valid keys as your primary active keys
    withing client systems.
  * **Grace Period** - A key past it's prime but still usable.  A grace period provides overlap to allow applications to
    transition to newer _valid_ keys without interrupting existing services.  Like milk past it's prime so no cereal but
    maybe you'll use it in mac'n'cheese.
  * **Expired** - A key well past it's prime.  `key-rotation` will delete these keys upon apply.

## Bindings
* [AWS](awskeystore)


## Development

By default `key-rotation` will build a CLI capable of interacting with bound key stores.  Build and check out the help!
```bash
go build .
```

### Programmatically

```go
package somepackage

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/truewhitespace/key-rotation/awskeystore"
	"github.com/truewhitespace/key-rotation/rotation"
	"time"
)

type Output struct {
	ID     string
	Secret *string
}

func DoKeyRotation(ctx context.Context, username string, iamSystem *iam.IAM) (*Output, error) {
	var err error
	keystore := awskeystore.NewAWSUserKeyStore(username, iamSystem)
	rotator, err := rotation.NewKeyRotation(72 * time.Hour, 48 * time.Hour)
	if err != nil { return nil, err }

	plan, err := rotator.Plan(ctx, keystore)
	if err != nil { return nil, err }
	
	keys, err := plan.Apply(ctx,keystore)
	if err != nil { return nil, err }
	
	key := keys[0].(*awskeystore.AWSAccessKey)
	return &Output{
		ID:     key.ID,
		Secret: key.Secret,
	}, nil
}
```
