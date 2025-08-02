//enum list
enum Category: String {

	//case
	case id, created, money, active, kyc, oid, type, verified, title, html, tags, channels, channels64, floats, keys, features, likes, providers, stats, price, meta, timeout, value, raw

	//name
	var string: String {
		get {
			switch self {
				case .id:          return "id"
				case .created:     return "created"
				case .money:       return "money"
				case .active:      return "active"
				case .kyc:         return "kyc"
				case .oid:         return "oid"
				case .type:        return "type"
				case .verified:    return "verified"
				case .title:       return "title"
				case .html:        return "html"
				case .tags:        return "tags"
				case .channels:    return "channels"
				case .channels64:  return "channels64"
				case .floats:      return "floats"
				case .keys:        return "keys"
				case .features:    return "features"
				case .likes:       return "likes"
				case .providers:   return "providers"
				case .stats:       return "stats"
				case .price:       return "price"
				case .meta:        return "meta"
				case .timeout:     return "timeout"
				case .value:       return "value"
				case .raw:         return "raw"
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .id:          return "Account ID"
				case .created:     return "Created"
				case .money:       return "Count"
				case .active:      return "Active"
				case .kyc:         return "KYC"
				case .oid:         return "OID"
				case .type:        return "Type"
				case .verified:    return "Verify"
				case .title:       return "Title"
				case .html:        return "Html"
				case .tags:        return "Tags"
				case .channels:    return "Channels"
				case .channels64:  return "Channels64"
				case .floats:      return "Floats"
				case .keys:        return "Keys"
				case .features:    return "Features"
				case .likes:       return "Likes"
				case .providers:   return "Providers"
				case .stats:       return "Stats"
				case .price:       return "Price"
				case .meta:        return "Meta"
				case .timeout:     return "Timeout"
				case .value:       return "Value"
				case .raw:         return "Raw"
			}
		}
	}
	
}
