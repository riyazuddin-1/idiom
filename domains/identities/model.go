package identities

type Identity struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	ProjectID string `json:"project_id" bson:"project_id"`
}
