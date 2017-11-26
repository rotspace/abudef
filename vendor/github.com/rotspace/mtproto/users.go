package mtproto

import "fmt"

func (m *MTProto) UsersGetFullUsers(id TL) (*TL_userFull, error) {
	var user TL_userFull
	x := <-m.InvokeAsync(TL_users_getFullUser{
		Id: id,
	})
	if x.err != nil {
		return nil, x.err
	}

	switch x.data.(type) {
	case TL_userFull:
		user = x.data.(TL_userFull)
	default:
		return nil, fmt.Errorf("Got: %T", x)
	}

	return &user, nil
}
