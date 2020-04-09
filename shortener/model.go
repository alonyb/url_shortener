package shortener

type Redirect struct {
	Code     string	`json:"code" bson:"code" msgpack:"code"`
	URL      string	`json:"url" bson:"url" msgpack:"url" validated:"empty=false & format=url"`
	CreateAt string	`json:"created_at" bson:"created_at" msgpack:"created_at"`
}
