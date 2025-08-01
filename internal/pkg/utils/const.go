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
	GSeymor         = "gseymor"
)

const (
	CollectionExp         = 50
	CIExp                 = 1
	EntityAttachExp       = 1
	UserSelfSubExp        = 15
	CollectionSelfSubExp  = 10
	CollectionSelfLikeExp = 5
	CISelfLikeExp         = 1
	CIPurchaseDateExp     = 1
	CIPurchasePriceExp    = 1
	CIRatingExp           = 1
	PictureExp            = 1
	MonthlyDonateExp      = 100
	YearlyDonateExp       = 1000
	AvatarAddExt          = 50
	BannerAddExp          = 50
	CICopyNumberExp       = 1
)

func GetAdminLogins() []string {
	return []string{Seymor, Admin, GodSeymor, EgRifYkkurIBita, GSeymor}
}
