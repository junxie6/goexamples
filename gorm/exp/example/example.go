package example

type Product struct {
	gorm.Model
	SKU         string
	Price       uint
	Description string
}

func create() {
}

func migrate() {
	// Migrate the schema
	db.AutoMigrate(&Product{})
}

func addProduct() {
	// Create
	db.Create(&Product{
		SKU:   "ABC123",
		Price: 1000,
	})
}

func updateProduct() {
	// Read
	var product Product

	//db.First(&product, 1)                   // find product with id 1
	db.First(&product, "SKU = ?", "WXR123") // find product with code l1212

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)
}

func deleteProduct() {
	// Delete - delete product
	//db.Delete(&product)
}

func changeColumnType() {
	db.Model(&Product{}).ModifyColumn("description", "text")
}

// 1st param : foreignkey field
// 2nd param : destination table(id)
// 3rd param : ONDELETE
// 4th param : ONUPDATE
func addForeignKey() {
	//db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
}
