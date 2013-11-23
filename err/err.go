/**
 * User: jackong
 * Date: 11/12/13
 * Time: 11:14 AM
 */
package err

type AccessError struct {
	Status int
	Msg string
}

func (this AccessError) Error() string{
    return this.Msg
}

const (
	CODE_OK = 0
)
