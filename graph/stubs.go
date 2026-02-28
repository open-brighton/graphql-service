package graph

import (
	"encoding/base64"
	"fmt"
	"strconv"

	model1 "github.com/openbrighton/graphql-service/graph/model"
)

// ─── Stub data ────────────────────────────────────────────────────────────────

var stubFeedItems = []*model1.FeedItem{
	{ID: "feed-1", Title: "Library Event", Description: "Story time & activities for all ages at Brighton Memorial Library.", Category: model1.FeedItemCategoryCommunityEvent, Date: "2025-03-08"},
	{ID: "feed-2", Title: "Town Hall Meeting", Description: "Community discussion on the upcoming budget and local projects.", Category: model1.FeedItemCategoryAnnouncement, Date: "2025-03-12"},
	{ID: "feed-3", Title: "Bulk Yard Pickup", Description: "Curbside collection dates for large items — schedule yours now.", Category: model1.FeedItemCategoryNews, Date: "2025-03-15"},
	{ID: "feed-4", Title: "Road Closure: Westfall Rd", Description: "Westfall Road will be closed between 5–9 AM for utility work.", Category: model1.FeedItemCategoryAlert, Date: "2025-03-20"},
}

var stubEvents = []*model1.Event{
	{ID: "evt-1", Title: "Brighton Food Festival", Description: "A celebration of local food and drink from Brighton's finest vendors.", Date: "Sat 7 Mar · 10:00 – 18:00", Location: "New Road, Brighton", Category: model1.EventCategoryFood},
	{ID: "evt-2", Title: "Open Mic Night", Description: "Local musicians and comedians take the stage every Friday night.", Date: "Fri 13 Mar · 19:30", Location: "The Haunt, Pool Valley", Category: model1.EventCategoryArts},
	{ID: "evt-3", Title: "Beach Clean-Up", Description: "Join volunteers for a community beach clean along the seafront.", Date: "Sun 15 Mar · 09:00 – 12:00", Location: "Brighton Beach, West Pier", Category: model1.EventCategoryEnvironment},
	{ID: "evt-4", Title: "Council Budget Consultation", Description: "Have your say on how Brighton's budget is spent this year.", Date: "Wed 18 Mar · 18:00 – 20:00", Location: "Hove Town Hall", Category: model1.EventCategoryPolitics},
	{ID: "evt-5", Title: "Brighton Half Marathon", Description: "Annual half marathon along the iconic Brighton seafront.", Date: "Sun 22 Mar · 09:30", Location: "Madeira Drive, Seafront", Category: model1.EventCategorySports},
	{ID: "evt-6", Title: "Community Gardening Day", Description: "Help plant and tend the neighbourhood community garden.", Date: "Sat 29 Mar · 10:00", Location: "Preston Park, Brighton", Category: model1.EventCategoryCommunity},
}

var stubBusinesses = []*model1.Business{
	{ID: "biz-1", Slug: "the-flour-pot-bakery", Name: "The Flour Pot Bakery", Category: model1.BusinessCategoryBakery, Address: "40 Sydney St, North Laine", Description: "Award-winning artisan bakery in the heart of North Laine, known for sourdough and pastries.", Longitude: -0.1376, Latitude: 50.8257, AccentColor: "#C8835A"},
	{ID: "biz-2", Slug: "choccywoccydoodah", Name: "Choccywoccydoodah", Category: model1.BusinessCategoryChocolatier, Address: "24 Duke St, Brighton", Description: "Brighton's most extravagant chocolate boutique — cakes and truffles like nowhere else.", Longitude: -0.1419, Latitude: 50.8225, AccentColor: "#7B4226"},
	{ID: "biz-3", Slug: "snoopers-paradise", Name: "Snoopers Paradise", Category: model1.BusinessCategoryVintage, Address: "7-8 Kensington Gardens", Description: "A legendary two-floor treasure trove of vintage clothing, furniture, and collectibles.", Longitude: -0.1368, Latitude: 50.8263, AccentColor: "#6A5ACD"},
	{ID: "biz-4", Slug: "infinity-foods", Name: "Infinity Foods", Category: model1.BusinessCategoryHealthFood, Address: "25 North Rd, Brighton", Description: "Brighton's original wholefoods co-op, stocking organic and ethical products since 1971.", Longitude: -0.1389, Latitude: 50.8258, AccentColor: "#4A8C49"},
	{ID: "biz-5", Slug: "brighton-fishing-museum", Name: "Brighton Fishing Museum", Category: model1.BusinessCategoryMuseum, Address: "201 King's Rd Arches", Description: "Free museum celebrating Brighton's fishing heritage, housed in the historic seafront arches.", Longitude: -0.1567, Latitude: 50.8196, AccentColor: "#2E6DA4"},
}

// ─── Pagination helpers ───────────────────────────────────────────────────────

func encodeCursor(index int) string {
	return base64.StdEncoding.EncodeToString([]byte("cursor:" + strconv.Itoa(index)))
}

func decodeCursor(cursor string) (int, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	s := string(b)
	if len(s) < 7 {
		return 0, fmt.Errorf("invalid cursor")
	}
	return strconv.Atoi(s[7:])
}

func pageSize(first *int32, total int) int {
	if first == nil {
		return total
	}
	n := int(*first)
	if n > total {
		return total
	}
	return n
}

func startIndex(after *string) int {
	if after == nil {
		return 0
	}
	idx, err := decodeCursor(*after)
	if err != nil {
		return 0
	}
	return idx + 1
}
