package config

// MongoConfiguration is used to store configuration for mongo
type MongoConfiguration struct {
	Addresses          []string
	Username, Password string
	AuthDB             string
	Database           string
}

// URI returns uri to connect to mongo
// mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
func (m MongoConfiguration) URI() string {
	uri := "mongodb://"
	if m.Username != "" && m.Password != "" {
		uri += (m.Username + ":" + m.Password + "@")
	}
	for i, address := range m.Addresses {
		if i != 0 {
			uri += ","
		}
		uri += address
	}
	if m.Database != "" {
		uri += ("/" + m.AuthDB)
	}
	return uri
}
