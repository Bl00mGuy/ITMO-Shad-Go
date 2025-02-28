//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	changes := map[int]int{}
	for _, guest := range guests {
		changes[guest.CheckInDate]++
		changes[guest.CheckOutDate]--
	}

	days := []int{}
	for day := range changes {
		days = append(days, day)
	}
	sort.Ints(days)

	var result []Load
	current := 0
	for _, day := range days {
		updated := current + changes[day]
		if updated != current {
			result = append(result, Load{
				StartDate:  day,
				GuestCount: updated,
			})
		}
		current = updated
	}
	return result
}
