package gotag

const (
	JSON = "json"
	BSON = "bson"
)

type TagTransHandle func(string) string

var tagHandleSet = map[string]TagTransHandle{}

func AddTag(tag string, handle TagTransHandle) {
	tagHandleSet[tag] = handle
}
func GetTagHandle(tag string) TagTransHandle {
	return tagHandleSet[tag]
}

func ToSnake(in string) string {
	runes := []rune(in)
	out := make([]rune, 0)
	for i, r := range runes {
		if 'A' <= r && r <= 'Z' {
			s := r + 32 // 'a'-'A' = 32
			if i == 0 {
				out = append(out, s)
			} else {
				out = append(out, '_')
				out = append(out, s)
			}
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

func init() {
	AddTag(JSON, ToSnake)
	AddTag(BSON, ToSnake)
}
