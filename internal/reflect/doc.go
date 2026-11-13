// Package reflect provides functions based on reflection.
package reflect

//revive:disable:line-length-limit // go:generate line length.

//go:generate mockgen -package=reflect_test -destination=mock_callback_test.go -source=walker_test.go Callback

//revive:enable:line-length-limit
