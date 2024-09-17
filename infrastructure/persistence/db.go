package persistence

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repositories struct holds the repositories for the application
type Repositories struct {
	User               repository.UserRepository
	Category           repository.CategoryRepository
	Location           repository.LocationRepository
	Size               repository.SizeRepository
	City               repository.CityRepository
	ShipmentContent    repository.ShipmentContentRepository
	ExtraService       repository.ExtraServiceRepository
	TransportationMode repository.TransportationModeRepository
	TruckType          repository.TruckTypeRepository
	TruckModel         repository.TruckModelRepository
	DeliveryTime       repository.DeliveryTimeRepository
	Driver             repository.DriverRepository
	Order              repository.OrderRepository
	Rating             repository.RatingRepository
	Page               repository.PageRepository
	FAQ                repository.FAQRepository
	Setting            repository.SettingRepository
	PasswordReset      repository.PasswordResetRepository
	PhoneVerification  repository.PhoneVerificationRepository
	IdentityDocument   repository.IdentityDocumentRepository
	Balance            repository.BalanceRepository
	Offer              repository.OfferRepository
	Device             repository.DeviceRepository
	db                 *gorm.DB
}

// NewRepositories returns a new instance of Repositories
func NewRepositories(host, port, user, password, dbname, sslMode, timeZone string) (*Repositories, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslMode, timeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repositories{
		User:               NewUserRepository(db),
		Category:           NewCategoryRepository(db),
		Location:           NewLocationRepository(db),
		Size:               NewSizeRepository(db),
		City:               NewCityRepository(db),
		ShipmentContent:    NewShipmentContentRepository(db),
		ExtraService:       NewExtraServiceRepository(db),
		TransportationMode: NewTransportationModeRepository(db),
		TruckType:          NewTruckTypeRepository(db),
		TruckModel:         NewTruckModelRepository(db),
		DeliveryTime:       NewDeliveryTimeRepository(db),
		Driver:             NewDriverRepository(db),
		Order:              NewOrderRepository(db),
		Rating:             NewRatingRepository(db),
		Page:               NewPageRepository(db),
		FAQ:                NewFAQRepository(db),
		Setting:            NewSettingRepository(db),
		PasswordReset:      NewPasswordResetRepository(db),
		PhoneVerification:  NewPhoneVerificationRepository(db),
		IdentityDocument:   NewIdentityDocumentRepository(db),
		Balance:            NewBalanceRepository(db),
		Offer:              NewOfferRepository(db),
		Device:             NewDeviceRepository(db),
		db:                 db,
	}, nil
}

// AutoMigrate creates the necessary tables in the database
func (r *Repositories) AutoMigrate() error {
	return r.db.AutoMigrate(&entity.User{}, &entity.Category{}, &entity.Location{}, &entity.Size{}, &entity.ShipmentContent{}, &entity.ExtraService{}, &entity.TransportationMode{}, &entity.TruckType{}, &entity.TruckModel{}, &entity.DeliveryTime{}, &entity.Driver{}, &entity.Order{}, &entity.Rating{}, &entity.Page{}, &entity.FAQ{}, &entity.OrderShipmentContent{}, &entity.OrderExtraService{}, &entity.OrderDriverPool{}, &entity.OrderDriverPool{}, &entity.Setting{}, &entity.PasswordReset{}, &entity.PhoneVerification{}, &entity.IdentityDocument{}, &entity.Balance{}, &entity.Offer{}, &entity.Device{}, &entity.City{})
}

// SeedCategories seeds the categories into the database.
// It creates a list of categories and inserts each category into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedCategories() {
	// List of categories to be seeded into the database
	categories := []entity.Category{
		{Name: `{"en": "Express", "ar": "سريع"}`, MenuName: `{"en": "Express Parcels", "ar": "شحن سريع"}`, Icon: "delivery.png", MenuIcon: ".png", Order: &[]int64{1}[0], MenuOrder: &[]int64{2}[0]},
		{Name: `{"en": "Truck", "ar": "شاحنة"}`, MenuName: `{"en": "Trucks", "ar": "شاحنات"}`, Icon: "freezer-truck.png", MenuIcon: ".png", IsTruck: &[]bool{true}[0], Order: &[]int64{2}[0], MenuOrder: &[]int64{3}[0]},
		{Name: `{"en": "Ridesharing", "ar": "مشوار"}`, MenuName: `{"en": "Ridesharing", "ar": "توصيل ركاب"}`, Icon: "driving.png", MenuIcon: ".png", Order: &[]int64{3}[0], MenuOrder: &[]int64{1}[0]},
		{Name: `{"en": "Pickup Truck", "ar": "بيك اب"}`, MenuName: `{"en": "Pickup", "ar": "نقل بيك اب"}`, Icon: "delivery-truck.png", MenuIcon: "pickup-truck-1.png", Order: &[]int64{4}[0], MenuOrder: &[]int64{4}[0]},
		{Name: `{"en": "Flat Truck", "ar": "سطحة"}`, MenuName: `{"en": "Flat Trucks", "ar": "سطحة"}`, Icon: "flat-bed-truck.png", MenuIcon: "car 1.png", Order: &[]int64{5}[0], MenuOrder: &[]int64{5}[0]},
	}

	// Iterate through the list of categories and insert each category into the database
	for _, category := range categories {
		r.db.Create(&category)
	}

	log.Println("Categories have been seeded successfully into the database.")
}

// SeedLocations seeds the delivery times into the database.
// It creates a list of delivery times and inserts each delivery time into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedDeliveryTimes() {
	// List of delivery times to be seeded into the database
	deliveryTimes := []entity.DeliveryTime{
		{Name: `{"en": "4 Hours", "ar": "4 ساعات"}`, Duration: 14400},
		{Name: `{"en": "8 Hours", "ar": "8 ساعات"}`, Duration: 28800},
		{Name: `{"en": "12 Hours", "ar": "12 ساعة"}`, Duration: 43200},
		{Name: `{"en": "24 Hours", "ar": "24 ساعة"}`, Duration: 86400},
	}

	// Iterate through the list of delivery times and insert each delivery time into the database
	for _, deliveryTime := range deliveryTimes {
		r.db.Create(&deliveryTime)
	}

	log.Println("Delivery times have been seeded successfully into the database.")
}

// SeedTruckTypes seeds the truck types into the database.
// It creates a list of truck types and inserts each truck type into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedTruckTypes() {
	// List of truck types to be seeded into the database
	truckTypes := []entity.TruckType{
		{Name: `{"en": "Sides", "ar": "جوانب"}`},
		{Name: `{"en": "Closed", "ar": "مغلقة"}`},
	}

	// Iterate through the list of truck types and insert each truck type into the database
	for _, truckType := range truckTypes {
		r.db.Create(&truckType)
	}

	log.Println("Truck types have been seeded successfully into the database.")
}

// SeedTruckModels seeds the truck models into the database.
// It creates a list of truck models and inserts each truck model into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedTruckModels() {
	// List of truck models to be seeded into the database
	truckModels := []entity.TruckModel{
		{Name: `{"en": "Diana", "ar": "دينا"}`},
		{Name: `{"en": "Lorry", "ar": "لوري"}`},
		{Name: `{"en": "Traila", "ar": "تريلا"}`},
	}

	// Iterate through the list of truck models and insert each truck model into the database
	for _, truckModel := range truckModels {
		r.db.Create(&truckModel)
	}

	log.Println("Truck models have been seeded successfully into the database.")
}

// SeedLocations seeds the locations into the database.
// It creates a list of locations and inserts each location into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedLocations() {

	print("I triend entering here")
	// List of locations to be seeded into the database
	locations := []entity.Location{
		{Name: `{"en":"Riyadh","ar":"الرياض"}`, Latitude: 24.7136, Longitude: 46.6753},
		{Name: `{"en":"Jeddah","ar":"جدة"}`, Latitude: 21.5169, Longitude: 39.2192},
		{Name: `{"en":"Mecca","ar":"مكة المدينة"}`, Latitude: 21.4278, Longitude: 39.8256},
		{Name: `{"en":"Medina","ar":"المدينة المنورة"}`, Latitude: 24.4977, Longitude: 39.6042},
		{Name: `{"en":"Dammam","ar":"الدمام"}`, Latitude: 26.3924, Longitude: 50.0993},
		{Name: `{"en":"Al Khobar","ar":"الخبر"}`, Latitude: 26.2866, Longitude: 50.1719},
		{Name: `{"en":"Buraidah","ar":"بريدة"}`, Latitude: 26.3357, Longitude: 43.9824},
		{Name: `{"en":"Taif","ar":"الطائف"}`, Latitude: 21.2779, Longitude: 40.4158},
		{Name: `{"en":"Tabuk","ar":"تبوك"}`, Latitude: 28.3833, Longitude: 36.5667},
		{Name: `{"en":"Hofuf","ar":"الهفوف"}`, Latitude: 25.3889, Longitude: 49.5595},
	}

	// Iterate through the list of locations and insert each location into the database
	for _, location := range locations {
		r.db.Create(&location)
	}

	log.Println("Locations have been seeded successfully into the database.")
}

// SeedLocations seeds the users into the database.
// It creates a list of users and inserts each user into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedUsers() error {
	// List of users to be seeded into the database
	users := []entity.User{
		{LocationID: rand.Uint64()%10 + 1, Name: "Super Admin", Phone: randomPhoneNumber(), Password: "password", Role: "admin"},
		{LocationID: rand.Uint64()%10 + 1, Name: "User", Phone: randomPhoneNumber(), Password: "password", Role: "user"},
	}

	// Iterate through the list of users and insert each user into the database
	for _, user := range users {
		if err := r.db.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Println("Users have been seeded successfully into the database.")
	return nil
}

// SeedShipmentContents seeds the shipment contents into the database.
// It creates a list of shipment contents and inserts each shipment content into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedShipmentContents() {
	// List of shipment contents to be seeded into the database
	shipmentContents := []entity.ShipmentContent{
		{Name: `{"en": "Dry", "ar": "جاف"}`, Icon: "box.png"},
		{Name: `{"en": "Cool", "ar": "مبرد"}`, Icon: "snowflake.png"},
		{Name: `{"en": "Medicines", "ar": "أدوية"}`, Icon: "herbal.png"},
		{Name: `{"en": "Beauty", "ar": "العطور والتجميل"}`, Icon: "perfume.png"},
		{Name: `{"en": "Frozen", "ar": "مجمد"}`, Icon: "frost.png"},
		{Name: `{"en": "Plants", "ar": "نباتات"}`, Icon: "plants.png"},
		{Name: `{"en": "Vegetables & Fruits", "ar": "خضار و فواكه"}`, Icon: "vegetable.png"},
		{Name: `{"en": "Dates", "ar": "تمور"}`, Icon: "date-palm.png"},
		{Name: `{"en": "Animals", "ar": "حيوانات"}`, Icon: "lamb.png"},
		{Name: `{"en": "Documents", "ar": "مستندات"}`, Icon: "documents.png"},
		{Name: `{"en": "Sensitive Document", "ar": "وثائق خاصة"}`, Icon: "document.png"},
		{Name: `{"en": "Furniture", "ar": "عفش و أثاث"}`, Icon: "sofa.png"},
	}

	// Iterate through the list of shipment contents and insert each shipment content into the database
	for _, shipmentContent := range shipmentContents {
		r.db.Create(&shipmentContent)
	}

	log.Println("Shipment content have been seeded successfully into the database.")
}

// SeedExtraServices seeds the extra services into the database.
// It creates a list of extra services and inserts each extra service into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedExtraServices() {
	// List of extra services to be seeded into the database
	extraServices := []entity.ExtraService{
		{Name: `{"en": "Pickup & Drop off", "ar": "تحميل و تنزيل"}`, Icon: "box-car-estate-loader.png"},
	}

	// Iterate through the list of extra services and insert each extra service into the database
	for _, extraService := range extraServices {
		r.db.Create(&extraService)
	}

	log.Println("Extra service have been seeded successfully into the database.")
}

// SeedSizes seeds the sizes into the database.
// It creates a list of sizes and inserts each size into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedSizes() {
	// List of sizes to be seeded into the database
	sizes := []entity.Size{
		{Name: `{"en": "Small", "ar": "صغير"}`, Description: `{"en": "Up to 10 kg", "ar": "حتى 10 كيلو للطلب الواحد"}`},
		{Name: `{"en": "Medium", "ar": "وسط"}`, Description: `{"en": "Up to 20 kg", "ar": "حتى 20 كيلو للطلب الواحد"}`},
		{Name: `{"en": "Large", "ar": "كبير"}`, Description: `{"en": "Up to 30 kg", "ar": "حتى 30 كيلو للطلب الواحد"}`},
	}

	// Iterate through the list of sizes and insert each size into the database
	for _, size := range sizes {
		r.db.Create(&size)
	}

	log.Println("Sizes have been seeded successfully into the database.")
}

// SeedTransportationModes seeds the transportation modes into the database.
// It creates a list of transportation modes and inserts each transportation mode into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedTransportationModes() {
	// List of transportation modes to be seeded into the database
	transportationModes := []entity.TransportationMode{
		{Name: `{"en": "Car", "ar": "سيارة"}`, Marker: "car.png"},
		{Name: `{"en": "Truck", "ar": "شاحنة"}`, Marker: "truck.png"},
	}

	// Iterate through the list of transportation modes and insert each transportation mode into the database
	for _, transportationMode := range transportationModes {
		r.db.Create(&transportationMode)
	}

	log.Println("Transportation modes have been seeded successfully into the database.")
}

// SeedPages seeds the pages into the database.
// It creates a list of pages and inserts each page into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedPages() {
	// List of pages to be seeded into the database
	pages := []entity.Page{
		{Name: `{"en": "Terms and Conditions", "ar": "الشروط والأحكام"}`, Body: `{"en": "<p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.</p><p>It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing</p><p>Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p><p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.</p><p>Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words,</p>", "ar": "<p>لوريم إيبسوم هو ببساطة نص شكلي يستخدم في صناعة الطباعة والتنضيد. كان Lorem Ipsum هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ، عندما أخذت طابعة غير معروفة لوحًا من النوع وتدافعت عليه لعمل كتاب عينة.</p><p>لقد صمد ليس فقط لخمسة قرون ، ولكن أيضًا القفزة في التنضيد الإلكتروني ، وظل دون تغيير جوهري. تم نشره في الستينيات مع إصدار أوراق Letraset التي تحتوي على</p><p>مقاطع Lorem Ipsum ، ومؤخرًا مع برامج النشر المكتبي مثل Aldus PageMaker بما في ذلك إصدارات Lorem Ipsum.</p><p>خلافًا للاعتقاد الشائع ، فإن Lorem Ipsum ليس مجرد نص عشوائي. لها جذور في قطعة من الأدب اللاتيني الكلاسيكي من 45 قبل الميلاد ، مما يجعلها أكثر من 2000 عام.</p><p>ريتشارد مكلينتوك ، أستاذ اللغة اللاتينية في كلية هامبدن-سيدني في فيرجينيا ، بحث عن واحدة من أكثر الكلمات اللاتينية غموضًا ،</p>"}`},
		{Name: `{"en": "Privacy Policy", "ar": "سياسة الخصوصية"}`, Body: `{"en": "<p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.</p><p>It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing</p><p>Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p><p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.</p><p>Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words,</p>", "ar": "<p>لوريم إيبسوم هو ببساطة نص شكلي يستخدم في صناعة الطباعة والتنضيد. كان Lorem Ipsum هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ، عندما أخذت طابعة غير معروفة لوحًا من النوع وتدافعت عليه لعمل كتاب عينة.</p><p>لقد صمد ليس فقط لخمسة قرون ، ولكن أيضًا القفزة في التنضيد الإلكتروني ، وظل دون تغيير جوهري. تم نشره في الستينيات مع إصدار أوراق Letraset التي تحتوي على</p><p>مقاطع Lorem Ipsum ، ومؤخرًا مع برامج النشر المكتبي مثل Aldus PageMaker بما في ذلك إصدارات Lorem Ipsum.</p><p>خلافًا للاعتقاد الشائع ، فإن Lorem Ipsum ليس مجرد نص عشوائي. لها جذور في قطعة من الأدب اللاتيني الكلاسيكي من 45 قبل الميلاد ، مما يجعلها أكثر من 2000 عام.</p><p>ريتشارد مكلينتوك ، أستاذ اللغة اللاتينية في كلية هامبدن-سيدني في فيرجينيا ، بحث عن واحدة من أكثر الكلمات اللاتينية غموضًا ،</p>"}`},
		{Name: `{"en": "Service Rules", "ar": "شروط الخدمة"}`, Body: `{"en": "<h2>Rules of condact for passengers</h2><p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.</p><p>It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing</p><p>Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p><p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.</p><p>Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words,</p>", "ar": "<h2>قواعد السلوك للركاب</h2><p>لوريم إيبسوم هو ببساطة نص شكلي يستخدم في صناعة الطباعة والتنضيد. كان Lorem Ipsum هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ، عندما أخذت طابعة غير معروفة لوحًا من النوع وتدافعت عليه لعمل كتاب عينة.</p><p>لقد صمد ليس فقط لخمسة قرون ، ولكن أيضًا القفزة في التنضيد الإلكتروني ، وظل دون تغيير جوهري. تم نشره في الستينيات مع إصدار أوراق Letraset التي تحتوي على</p><p>مقاطع Lorem Ipsum ، ومؤخرًا مع برامج النشر المكتبي مثل Aldus PageMaker بما في ذلك إصدارات Lorem Ipsum.</p><p>خلافًا للاعتقاد الشائع ، فإن Lorem Ipsum ليس مجرد نص عشوائي. لها جذور في قطعة من الأدب اللاتيني الكلاسيكي من 45 قبل الميلاد ، مما يجعلها أكثر من 2000 عام.</p><p>ريتشارد مكلينتوك ، أستاذ اللغة اللاتينية في كلية هامبدن-سيدني في فيرجينيا ، بحث عن واحدة من أكثر الكلمات اللاتينية غموضًا ،</p>"}`},
	}

	// Iterate through the list of pages and insert each page into the database
	for _, page := range pages {
		r.db.Create(&page)
	}

	log.Println("Pages have been seeded successfully into the database.")
}

// SeedFAQs seeds the FAQs into the database.
// It creates a list of FAQs and inserts each FAQ into the database using the GORM library.
// If the seed operation is successful, it logs a message indicating so.
func (r *Repositories) SeedFAQs() {
	// List of pages to be seeded into the database
	faqs := []entity.FAQ{
		{Question: `{"en": "What is Jayeek?", "ar": "ماذا يعني تطبيق جايك؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ،"}`},
		{Question: `{"en": "Which is the Services That Provide?", "ar": "ما هي الخدمة التي يقدمها التطبيق؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry.", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر  ،"}`},
		{Question: `{"en": "How I Can Pay to Jayeek?", "ar": "كيف يمكنني الدفع الى الخدمة؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ،"}`},
		{Question: `{"en": "What Are the Borders of Cities?", "ar": "ما هي حدود المدينة المتاحة؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ،"}`},
		{Question: `{"en": "How Can I Become a Driver?", "ar": "كيف يمكنني ان اكون سائقاً؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر  ،"}`},
		{Question: `{"en": "How Can I Get Support?", "ar": "كيف يمكنني الحصول على دعم فني؟"}`, Answer: `{"en": "is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,", "ar": "هو مجرد نص وهمي لصناعة الطباعة والتنضيد. لوريم إيبسوم هو النص الوهمي القياسي في الصناعة منذ القرن الخامس عشر الميلادي ،"}`},
	}

	// Iterate through the list of pages and insert each page into the database
	for _, faq := range faqs {
		r.db.Create(&faq)
	}

	log.Println("FAQs have been seeded successfully into the database.")
}

func (r *Repositories) SeedSettings() {
	// List of settings to be seeded into the database
	settings := []entity.Setting{
		{Key: "iso2", Value: "SA"},
		{Key: "currency_symbol", Value: "SAR"},
		{Key: "police_url", Value: "tel:999"},
		{Key: "support_url", Value: "https://jayeek.net/#contact"},
		{Key: "terms_page_id", Value: "1"},
		{Key: "privacy_page_id", Value: "2"},
		{Key: "rules_page_id", Value: "3"},
		{Key: "shipment_contents_max_selections", Value: "3"},
		{Key: "max_orders_per_trip", Value: "5"},
	}

	// Iterate through the list of transportation modes and insert each transportation mode into the database
	for _, setting := range settings {
		r.db.Create(&setting)
	}

	log.Println("Sizes have been seeded successfully into the database.")
}

func randomPhoneNumber() string {
	return "+966" + strconv.Itoa(rand.Intn(1000000000) + 1000000000)[1:]
}
