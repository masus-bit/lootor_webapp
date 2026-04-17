package utils

const (
	EventTargetCollection     = "collection"
	EventTargetUser           = "user"
	EventTargetCollectionItem = "collectionItem"
	EventTargetWL             = "wishListItem"
	EventTargetPost           = "post"
	EventTargetTag            = "tag"
	EventTargetPhoto          = "photo"
)

const (
	EventActionCreate      = "create"
	EventActionUpdate      = "update"
	EventActionDelete      = "delete"
	EventActionSubscribe   = "subscribe"
	EventActionLike        = "like"
	EventActionAddTag      = "tagAdd"
	EventActionDislike     = "dislike"
	EventActionUnsubscribe = "unsubscribe"
)

const (
	NotificationTypeCollection     = "collection"
	NotificationTypeUser           = "user"
	NotificationTypeCollectionItem = "collectionItem"
	NotificationTypePost           = "post"
	NotificationTypeComment        = "comment"
	NotificationTypeAnswer         = "answer"
	NotificationTypePhoto          = "photo"
	NotificationTypeAchievement    = "achievement"
)

const (
	NotificationActionSubscribe   = "subscribe"
	NotificationActionLike        = "like"
	NotificationActionComment     = "comment"
	NotificationActionReact       = "react"
	NotificationActionAchievement = "received"
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

const (
	AchieveBetaTester            = "betaTester"
	AchieveBetaTesterDonate      = "betaTesterDonate"
	AchieveFirstCollectionCreate = "firstCollectionCreate"
	AchieveDonate                = "donate"
	AchieveSubscribers           = "subscribers"
	AchieveYearlyRegister        = "yearlyRegister"
	AchieveCollectionItemsAdded  = "collectionItemsAdded"
	AchievePhotosAdded           = "photosAdded"
	AchieveTagsCreated           = "tagsCreated"
	AchievePostsCreated          = "postsCreated"
	AchieveCollectionsLikes      = "collectionsLikes"
	AchieveCollectionItemsLikes  = "collectionItemsLikes"
	AchievePhotosLikes           = "photosLikes"
	AchievePostsReactions        = "postsReactions"
	AchieveCollectionsSum        = "collectionsSum"
	AchieveCollectionShipSum     = "collectionShipSum"
)

const (
	XPBetaTester                 = 100
	XPBetaTesterDonate           = 500
	XPFirstCollectionCreate      = 25
	XPDonateLevel1               = 100
	XPDonateLevel2               = 1200
	XPSubscribersLevel1          = 10
	XPSubscribersLevel2          = 200
	XPSubscribersLevel3          = 500
	XPYearlyRegister             = 1000
	XPCollectionItemsAddLevel1   = 20
	XPCollectionItemsAddLevel2   = 250
	XPCollectionItemsAddLevel3   = 3000
	XPPhotosAddLevel1            = 20
	XPPhotosAddLevel2            = 250
	XPPhotosAddLevel3            = 3000
	XPTagsAddLevel1              = 50
	XPTagsAddLevel2              = 400
	XPTagsAddLevel3              = 1500
	XPPostsCreateLevel1          = 25
	XPPostsCreateLevel2          = 300
	XPPostsCreateLevel3          = 1000
	XPCollectionsLikesLevel1     = 20
	XPCollectionsLikesLevel2     = 150
	XPCollectionsLikesLevel3     = 500
	XPCollectionItemsLikesLevel1 = 50
	XPCollectionItemsLikesLevel2 = 300
	XPCollectionItemsLikesLevel3 = 1500
	XPPhotosLikesLevel1          = 50
	XPPhotosLikesLevel2          = 300
	XPPhotosLikesLevel3          = 1500
	XPPostsReactionsLevel1       = 50
	XPPostsReactionsLevel2       = 300
	XPPostsReactionsLevel3       = 1500
	XPCollectionsSumLevel1       = 50
	XPCollectionsSumLevel2       = 300
	XPCollectionsSumLevel3       = 1500
	XPCollectionsShipSumLevel1   = 50
	XPCollectionsShipSumLevel2   = 250
	XPCollectionsShipSumLevel3   = 1000
)

const (
	SubscribersLevel1          = 1
	SubscribersLevel2          = 20
	SubscribersLevel3          = 50
	CollectionItemsAddLevel1   = 10
	CollectionItemsAddLevel2   = 75
	CollectionItemsAddLevel3   = 200
	PhotosAddLevel1            = 10
	PhotosAddLevel2            = 75
	PhotosAddLevel3            = 200
	TagsAddLevel1              = 5
	TagsAddLevel2              = 30
	TagsAddLevel3              = 100
	PostsCreateLevel1          = 1
	PostsCreateLevel2          = 15
	PostsCreateLevel3          = 75
	CollectionsLikesLevel1     = 10
	CollectionsLikesLevel2     = 50
	CollectionsLikesLevel3     = 200
	CollectionItemsLikesLevel1 = 50
	CollectionItemsLikesLevel2 = 250
	CollectionItemsLikesLevel3 = 1000
	PhotosLikesLevel1          = 50
	PhotosLikesLevel2          = 250
	PhotosLikesLevel3          = 1000
	PostsReactionsLevel1       = 50
	PostsReactionsLevel2       = 250
	PostsReactionsLevel3       = 1000
	CollectionsSumLevel1       = 50000
	CollectionsSumLevel2       = 250000
	CollectionsSumLevel3       = 1000000
	CollectionsShipSumLevel1   = 4000
	CollectionsShipSumLevel2   = 20000
	CollectionsShipSumLevel3   = 80000
)

func GetAdminLogins() []string {
	return []string{Seymor, Admin, GodSeymor, EgRifYkkurIBita, GSeymor}
}
