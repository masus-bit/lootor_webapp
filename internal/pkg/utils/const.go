package utils

const (
	EventTargetCollection     = "collection"
	EventTargetUser           = "user"
	EventTargetCollectionItem = "collectionItem"
	EventTargetWL             = "wishListItem"
	EventTargetPost           = "post"
)

const (
	EventActionCreate    = "create"
	EventActionUpdate    = "update"
	EventActionDelete    = "delete"
	EventActionSubscribe = "subscribe"
	EventActionLike      = "like"
)

const (
	NotificationTypeCollection     = "collection"
	NotificationTypeUser           = "user"
	NotificationTypeCollectionItem = "collectionItem"
	NotificationTypePost           = "post"
	NotificationTypeComment        = "comment"
	NotificationTypeAnswer         = "answer"
)

const (
	NotificationActionSubscribe = "subscribe"
	NotificationActionLike      = "like"
	NotificationActionComment   = "comment"
	NotificationActionReact     = "react"
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
	TagAttachExp          = 1
	UserSelfSubExp        = 15
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
	TagExp                = 2
	PostCreateExp         = 10
	PostReactExt          = 1
)

func GetAdminLogins() []string {
	return []string{Seymor, Admin, GodSeymor, EgRifYkkurIBita, GSeymor}
}
