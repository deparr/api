package cache

// There's gotta be a better way to use these types??
type pinnedRepoData struct {
	Data struct {
		User struct {
			PinnedItems struct {
				TotalCount int
				Edges      []struct {
					Node struct {
						Owner          struct{ Login string }
						Name           string
						Url            string
						Description    string
						StargazerCount int
						PushedAt       string
						Languages      struct {
							Edges []struct {
								Node struct {
									Name  string
									Color string
								}
								Size int
							}
						}
					}
				}
			}
		}
	}
}
