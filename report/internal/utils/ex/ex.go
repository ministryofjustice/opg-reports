package ex

import "slices"

func FilterValue(src map[string]string, remove ...string) map[string]string {
	var updated = src
	for k, v := range updated {
		if slices.Contains(remove, v) {
			delete(updated, k)
		}
	}
	return updated
}

func FilterKeys(src map[string]string, remove ...string) map[string]string {
	var updated = src
	for k, _ := range updated {
		if slices.Contains(remove, k) {
			delete(updated, k)
		}
	}
	return updated
}
