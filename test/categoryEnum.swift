//enum list
enum Category: String {

	//case
	case cars, realestate, watch, jewelry

	//name
	var string: String {
		get {
			switch self {
				case .cars:          return "cars"
				case .realestate:    return "realestate"
				case .watch:         return "watch"
				case .jewelry:       return "jewelry"
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .cars:          return "Cars"
				case .realestate:    return "Realestate"
				case .watch:         return "Watch"
				case .jewelry:       return "Jewelry"
			}
		}
	}
	
}
