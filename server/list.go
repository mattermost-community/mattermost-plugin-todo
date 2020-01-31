package main

type list struct {
	listID string
	userID string
	p      *Plugin
}

const (
	// MyListKey is the key used to store the order of the owned todos
	MyListKey = ""
	// InListKey is the key used to store the order of received todos
	InListKey = "_in"
	// OutListKey is the key used to store the order of sent todos
	OutListKey = "_out"
)

func (p *Plugin) getInListForUser(userID string) *list {
	return &list{
		listID: InListKey,
		userID: userID,
		p:      p,
	}
}

func (p *Plugin) getOutListForUser(userID string) *list {
	return &list{
		listID: OutListKey,
		userID: userID,
		p:      p,
	}
}

func (p *Plugin) getMyListForUser(userID string) *list {
	return &list{
		listID: MyListKey,
		userID: userID,
		p:      p,
	}
}

func (p *Plugin) getUserListForItem(userID string, itemID string) (*list, *OrderElement, string) {
	itemList := p.getMyListForUser(userID)
	oe, _, _ := itemList.getOrderForItem(itemID)
	if oe != nil {
		return itemList, oe, MyListKey
	}

	itemList = p.getOutListForUser(userID)
	oe, _, _ = itemList.getOrderForItem(itemID)
	if oe != nil {
		return itemList, oe, OutListKey
	}

	itemList = p.getInListForUser(userID)
	oe, _, _ = itemList.getOrderForItem(itemID)
	if oe != nil {
		return itemList, oe, InListKey
	}

	return nil, nil, ""
}

func (l *list) getItems() ([]*Item, error) {
	return l.p.getItemListForUser(l.userID, l.listID)
}

func (l *list) add(itemID string, foreignItemID string, foreignUserID string) error {
	return l.p.addToListForUser(l.userID, itemID, l.listID, foreignItemID, foreignUserID)
}

func (l *list) remove(itemID string) error {
	return l.p.removeFromListForUser(l.userID, itemID, l.listID)
}

func (l *list) getOrderForItem(itemID string) (*OrderElement, int, error) {
	return l.p.getOrderForItem(l.userID, itemID, l.listID)
}
