package params

// Define available fleets.
const (
	FleetUndefined = ""
	FleetBeta      = "eth.beta"
	FleetStaging   = "eth.staging"
)

// Cluster defines a list of Ethereum nodes.
type Cluster struct {
	StaticNodes     []string `json:"staticnodes"`
	BootNodes       []string `json:"bootnodes"`
	MailServers     []string `json:"mailservers"` // list of trusted mail servers
	RendezvousNodes []string `json:"rendezvousnodes"`
}

// Consult this list with http://fleets.status.im/.
var clusters = map[string]func() Cluster{
	FleetStaging: func() Cluster {
		return Cluster{
			BootNodes: []string{
				"enode://10a78c17929a7019ef4aa2249d7302f76ae8a06f40b2dc88b7b31ebff4a623fbb44b4a627acba296c1ced3775d91fbe18463c15097a6a36fdb2c804ff3fc5b35@35.238.97.234:30404",   // boot-01.gc-us-central1-a.eth.staging
				"enode://f79fb3919f72ca560ad0434dcc387abfe41e0666201ebdada8ede0462454a13deb05cda15f287d2c4bd85da81f0eb25d0a486bbbc8df427b971ac51533bd00fe@174.138.107.239:30404", // boot-01.do-ams3.eth.staging
			},
			StaticNodes: []string{
				"enode://914c0b30f27bab30c1dfd31dad7652a46fda9370542aee1b062498b1345ee0913614b8b9e3e84622e84a7203c5858ae1d9819f63aece13ee668e4f6668063989@167.99.19.148:30305", // node-01.do-ams3.eth.staging
				"enode://2d897c6e846949f9dcf10279f00e9b8325c18fe7fa52d658520ad7be9607c83008b42b06aefd97cfe1fdab571f33a2a9383ff97c5909ed51f63300834913237e@35.192.0.86:30305",   // "node-01.gc-us-central1-a.eth.staging"
			},
			MailServers: []string{
				"enode://69f72baa7f1722d111a8c9c68c39a31430e9d567695f6108f31ccb6cd8f0adff4991e7fdca8fa770e75bc8a511a87d24690cbc80e008175f40c157d6f6788d48@206.189.240.16:30504", // mail-01.do-ams3.eth.staging
				"enode://e4fc10c1f65c8aed83ac26bc1bfb21a45cc1a8550a58077c8d2de2a0e0cd18e40fd40f7e6f7d02dc6cd06982b014ce88d6e468725ffe2c138e958788d0002a7f@35.239.193.41:30504",  // mail-01.gc-us-central1-a.eth.staging
			},
			RendezvousNodes: []string{
				"/ip4/174.138.107.239/tcp/30703/ethv4/16Uiu2HAkyJHeetQ4DNpd4NZ2ntzxMo25zcdpvGQRqkD5pB9BE6RU",
				"/ip4/35.238.97.234/tcp/30703/ethv4/16Uiu2HAm1sVyXmkMNjdeDWqK2urbyC3oBHi8MDpCdYkns1nYafqz",
			},
		}
	},
	FleetBeta: func() Cluster {
		return Cluster{
			BootNodes: []string{
				"enode://436cc6f674928fdc9a9f7990f2944002b685d1c37f025c1be425185b5b1f0900feaf1ccc2a6130268f9901be4a7d252f37302c8335a2c1a62736e9232691cc3a@174.138.105.243:30404", // boot-01.do-ams3.eth.beta
				"enode://5395aab7833f1ecb671b59bf0521cf20224fe8162fc3d2675de4ee4d5636a75ec32d13268fc184df8d1ddfa803943906882da62a4df42d4fccf6d17808156a87@206.189.243.57:30404",  // boot-02.do-ams3.eth.beta
				"enode://7427dfe38bd4cf7c58bb96417806fab25782ec3e6046a8053370022cbaa281536e8d64ecd1b02e1f8f72768e295d06258ba43d88304db068e6f2417ae8bcb9a6@104.154.88.123:30404",  // boot-01.gc-us-central1-a.eth.beta
				"enode://ebefab39b69bbbe64d8cd86be765b3be356d8c4b24660f65d493143a0c44f38c85a257300178f7845592a1b0332811542e9a58281c835babdd7535babb64efc1@35.202.99.224:30404",   // boot-02.gc-us-central1-a.eth.beta
			},
			StaticNodes: []string{
				"enode://a6a2a9b3a7cbb0a15da74301537ebba549c990e3325ae78e1272a19a3ace150d03c184b8ac86cc33f1f2f63691e467d49308f02d613277754c4dccd6773b95e8@206.189.243.176:30304", // node-01.do-ams3.eth.beta
				"enode://207e53d9bf66be7441e3daba36f53bfbda0b6099dba9a865afc6260a2d253fb8a56a72a48598a4f7ba271792c2e4a8e1a43aaef7f34857f520c8c820f63b44c8@35.224.15.65:30304",    // node-01.gc-us-central1-a.eth.beta
			},
			MailServers: []string{
				"enode://c42f368a23fa98ee546fd247220759062323249ef657d26d357a777443aec04db1b29a3a22ef3e7c548e18493ddaf51a31b0aed6079bd6ebe5ae838fcfaf3a49@206.189.243.162:30504", // mail-01.do-ams3.eth.beta
				"enode://7aa648d6e855950b2e3d3bf220c496e0cae4adfddef3e1e6062e6b177aec93bc6cdcf1282cb40d1656932ebfdd565729da440368d7c4da7dbd4d004b1ac02bf8@206.189.243.169:30504", // mail-02.do-ams3.eth.beta
				"enode://8a64b3c349a2e0ef4a32ea49609ed6eb3364be1110253c20adc17a3cebbc39a219e5d3e13b151c0eee5d8e0f9a8ba2cd026014e67b41a4ab7d1d5dd67ca27427@206.189.243.168:30504", // mail-03.do-ams3.eth.beta
				"enode://7de99e4cb1b3523bd26ca212369540646607c721ad4f3e5c821ed9148150ce6ce2e72631723002210fac1fd52dfa8bbdf3555e05379af79515e1179da37cc3db@35.188.19.210:30504",   // mail-01.gc-us-central1-a.eth.beta
				"enode://015e22f6cd2b44c8a51bd7a23555e271e0759c7d7f52432719665a74966f2da456d28e154e836bee6092b4d686fe67e331655586c57b718be3997c1629d24167@35.226.21.19:30504",    // mail-02.gc-us-central1-a.eth.beta
				"enode://531e252ec966b7e83f5538c19bf1cde7381cc7949026a6e499b6e998e695751aadf26d4c98d5a4eabfb7cefd31c3c88d600a775f14ed5781520a88ecd25da3c6@35.225.227.79:30504",   // mail-03.gc-us-central1-a.eth.beta
			},
			RendezvousNodes: []string{
				"/ip4/174.138.105.243/tcp/30703/ethv4/16Uiu2HAmRHPzF3rQg55PgYPcQkyvPVH9n2hWsYPhUJBZ6kVjJgdV", // boot-01.do-ams3.eth.beta
				"/ip4/206.189.243.57/tcp/30703/ethv4/16Uiu2HAmLqTXuY4Sb6G28HNooaFUXUKzpzKXCcgyJxgaEE2i5vnf",  // boot-02.do-ams3.eth.beta
			},
		}
	},
}

// ClusterForFleet returns a cluster for a given fleet.
func ClusterForFleet(fleet string) (Cluster, bool) {
	cluster, ok := clusters[fleet]
	if ok {
		return cluster(), true
	}
	return Cluster{}, false
}
