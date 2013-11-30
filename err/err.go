/**
 * User: jackong
 * Date: 11/12/13
 * Time: 11:14 AM
 */
package err

type HandlerError struct {
	Code int32
	Msg string
}

func (this HandlerError) Error() string{
    return this.Msg
}
