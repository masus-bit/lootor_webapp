package utils

const (
	EventTargetCollection     = "collection"
	EventTargetUser           = "user"
	EventTargetCollectionItem = "collectionItem"
	EventTargetWL             = "wishListItem"
)

const (
	EventActionCreate    = "create"
	EventActionUpdate    = "update"
	EventActionDelete    = "delete"
	EventActionSubscribe = "subscribe"
	EventActionLike      = "like"
)

const (
	Seymor          = "seymor"
	Admin           = "admin"
	GodSeymor       = "godSeymor"
	EgRifYkkurIBita = "eg_rif_ykkur_i_bita"
)

func GetAdminLogins() []string {
	return []string{Seymor, Admin, GodSeymor, EgRifYkkurIBita}
}
