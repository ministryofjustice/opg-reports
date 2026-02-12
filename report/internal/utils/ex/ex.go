package ex

func FilterValue(src map[string]string, remove string) map[string]string {
	var updated = src
	for _, v := range updated {
		if v == remove {
			delete(updated, v)
		}
	}
	return updated
}
