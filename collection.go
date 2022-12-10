package pocketbase

type Collection[T any] struct {
	*Client
	Name string
}

func CollectionSet[T any](client *Client, collection string) Collection[T] {
	return Collection[T]{client, collection}
}

func (c Collection[T]) Update(id string, body T) error {
	return c.Client.Update(c.Name, id, body)
}

func (c Collection[T]) Create(body T) (ResponseCreate, error) {
	return c.Client.Create(c.Name, body)
}

func (c Collection[T]) Delete(id string) error {
	return c.Client.Delete(c.Name, id)
}

func (c Collection[T]) List(params ParamsList) (ResponseList[T], error) {
	var response ResponseList[T]
	params.hackResponseRef = &response

	_, err := c.Client.List(c.Name, params)
	return response, err
}
