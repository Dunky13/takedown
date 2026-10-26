package source

import (
	"errors"
	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/drive/v3"
	"sort"
	"strings"
)

// Resource groups map 1:1 to OAuth scopes, as documented at
// https://developers.google.com/data-portability/user-guide/methods .
// Keep this explicit SDK-backed allowlist: a syntactically valid arbitrary scope
// is not evidence that Google supports exporting that resource group.
var resourceScopes = map[string]string{
	"alerts.subscriptions":                       dataportability.DataportabilityAlertsSubscriptionsScope,
	"chrome.autofill":                            dataportability.DataportabilityChromeAutofillScope,
	"chrome.bookmarks":                           dataportability.DataportabilityChromeBookmarksScope,
	"chrome.dictionary":                          dataportability.DataportabilityChromeDictionaryScope,
	"chrome.extensions":                          dataportability.DataportabilityChromeExtensionsScope,
	"chrome.history":                             dataportability.DataportabilityChromeHistoryScope,
	"chrome.reading_list":                        dataportability.DataportabilityChromeReadingListScope,
	"chrome.settings":                            dataportability.DataportabilityChromeSettingsScope,
	"discover.follows":                           dataportability.DataportabilityDiscoverFollowsScope,
	"discover.likes":                             dataportability.DataportabilityDiscoverLikesScope,
	"discover.not_interested":                    dataportability.DataportabilityDiscoverNotInterestedScope,
	"maps.aliased_places":                        dataportability.DataportabilityMapsAliasedPlacesScope,
	"maps.commute_routes":                        dataportability.DataportabilityMapsCommuteRoutesScope,
	"maps.commute_settings":                      dataportability.DataportabilityMapsCommuteSettingsScope,
	"maps.ev_profile":                            dataportability.DataportabilityMapsEvProfileScope,
	"maps.factual_contributions":                 dataportability.DataportabilityMapsFactualContributionsScope,
	"maps.offering_contributions":                dataportability.DataportabilityMapsOfferingContributionsScope,
	"maps.photos_videos":                         dataportability.DataportabilityMapsPhotosVideosScope,
	"maps.questions_answers":                     dataportability.DataportabilityMapsQuestionsAnswersScope,
	"maps.reviews":                               dataportability.DataportabilityMapsReviewsScope,
	"maps.starred_places":                        dataportability.DataportabilityMapsStarredPlacesScope,
	"maps.vehicle_profile":                       dataportability.DataportabilityMapsVehicleProfileScope,
	"myactivity.maps":                            dataportability.DataportabilityMyactivityMapsScope,
	"myactivity.myadcenter":                      dataportability.DataportabilityMyactivityMyadcenterScope,
	"myactivity.play":                            dataportability.DataportabilityMyactivityPlayScope,
	"myactivity.search":                          dataportability.DataportabilityMyactivitySearchScope,
	"myactivity.shopping":                        dataportability.DataportabilityMyactivityShoppingScope,
	"myactivity.youtube":                         dataportability.DataportabilityMyactivityYoutubeScope,
	"mymaps.maps":                                dataportability.DataportabilityMymapsMapsScope,
	"nest.camera_event":                          dataportability.DataportabilityNestCameraEventScope,
	"nest.camera_feature":                        dataportability.DataportabilityNestCameraFeatureScope,
	"nest.camera_video":                          dataportability.DataportabilityNestCameraVideoScope,
	"nest.store":                                 dataportability.DataportabilityNestStoreScope,
	"nest.user":                                  dataportability.DataportabilityNestUserScope,
	"order_reserve.purchases_reservations":       dataportability.DataportabilityOrderReservePurchasesReservationsScope,
	"pixel.device_data":                          dataportability.DataportabilityPixelDeviceDataScope,
	"play.devices":                               dataportability.DataportabilityPlayDevicesScope,
	"play.grouping":                              dataportability.DataportabilityPlayGroupingScope,
	"play.installs":                              dataportability.DataportabilityPlayInstallsScope,
	"play.library":                               dataportability.DataportabilityPlayLibraryScope,
	"play.playpoints":                            dataportability.DataportabilityPlayPlaypointsScope,
	"play.promotions":                            dataportability.DataportabilityPlayPromotionsScope,
	"play.purchases":                             dataportability.DataportabilityPlayPurchasesScope,
	"play.redemptions":                           dataportability.DataportabilityPlayRedemptionsScope,
	"play.subscriptions":                         dataportability.DataportabilityPlaySubscriptionsScope,
	"play.usersettings":                          dataportability.DataportabilityPlayUsersettingsScope,
	"saved.collections":                          dataportability.DataportabilitySavedCollectionsScope,
	"search_ugc.comments":                        dataportability.DataportabilitySearchUgcCommentsScope,
	"search_ugc.media.reviews_and_stars":         dataportability.DataportabilitySearchUgcMediaReviewsAndStarsScope,
	"search_ugc.media.streaming_video_providers": dataportability.DataportabilitySearchUgcMediaStreamingVideoProvidersScope,
	"search_ugc.media.thumbs":                    dataportability.DataportabilitySearchUgcMediaThumbsScope,
	"search_ugc.media.watched":                   dataportability.DataportabilitySearchUgcMediaWatchedScope,
	"searchnotifications.settings":               dataportability.DataportabilitySearchnotificationsSettingsScope,
	"searchnotifications.subscriptions":          dataportability.DataportabilitySearchnotificationsSubscriptionsScope,
	"shopping.addresses":                         dataportability.DataportabilityShoppingAddressesScope,
	"shopping.reviews":                           dataportability.DataportabilityShoppingReviewsScope,
	"streetview.imagery":                         dataportability.DataportabilityStreetviewImageryScope,
	"youtube.channel":                            dataportability.DataportabilityYoutubeChannelScope,
	"youtube.clips":                              dataportability.DataportabilityYoutubeClipsScope,
	"youtube.comments":                           dataportability.DataportabilityYoutubeCommentsScope,
	"youtube.conversations":                      dataportability.DataportabilityYoutubeConversationsScope,
	"youtube.live_chat":                          dataportability.DataportabilityYoutubeLiveChatScope,
	"youtube.music":                              dataportability.DataportabilityYoutubeMusicScope,
	"youtube.playable":                           dataportability.DataportabilityYoutubePlayableScope,
	"youtube.posts":                              dataportability.DataportabilityYoutubePostsScope,
	"youtube.private_playlists":                  dataportability.DataportabilityYoutubePrivatePlaylistsScope,
	"youtube.private_videos":                     dataportability.DataportabilityYoutubePrivateVideosScope,
	"youtube.public_playlists":                   dataportability.DataportabilityYoutubePublicPlaylistsScope,
	"youtube.public_videos":                      dataportability.DataportabilityYoutubePublicVideosScope,
	"youtube.shopping":                           dataportability.DataportabilityYoutubeShoppingScope,
	"youtube.subscriptions":                      dataportability.DataportabilityYoutubeSubscriptionsScope,
	"youtube.unlisted_playlists":                 dataportability.DataportabilityYoutubeUnlistedPlaylistsScope,
	"youtube.unlisted_videos":                    dataportability.DataportabilityYoutubeUnlistedVideosScope,
}

func normalizeResources(resources []string) ([]string, error) {
	result := make([]string, 0, len(resources))
	seen := map[string]bool{}
	for _, r := range resources {
		if _, ok := resourceScopes[r]; !ok {
			return nil, errors.New("unsupported Data Portability resource group; select an official resource group such as myactivity.search, without a scope URL")
		}
		if seen[r] {
			return nil, errors.New("duplicate Data Portability resource group")
		}
		seen[r] = true
		result = append(result, r)
	}
	sort.Strings(result)
	return result, nil
}

// Scopes deliberately requests only the selected provider's scopes. In
// particular Data Portability grants cannot be mixed with Drive/email/openid.
// Empty selections are allowed for setup; operations reject missing resources.
func Scopes(provider string, resources []string) ([]string, error) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	switch provider {
	case "", "drive":
		if len(resources) > 0 {
			return nil, errors.New("Drive cannot be combined with Data Portability resources; select a single provider")
		}
		// drive.file only covers individually authorized files. It does not grant
		// access to arbitrary current/future Takeout children of a selected folder.
		return []string{drive.DriveReadonlyScope}, nil
	case "portability":
		rs, err := normalizeResources(resources)
		if err != nil {
			return nil, err
		}
		scopes := make([]string, 0, len(rs))
		for _, r := range rs {
			scopes = append(scopes, resourceScopes[r])
		}
		return scopes, nil
	default:
		return nil, errors.New("unknown archive source provider; use drive or portability")
	}
}
