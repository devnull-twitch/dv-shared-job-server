package sharedjob

import "github.com/sirupsen/logrus"

type (
	CargoType     string
	CargoCategory string
)

const (
	None                CargoType = "None"
	Coal                CargoType = "Coal"
	IronOre             CargoType = "IronOre"
	CrudeOil            CargoType = "CrudeOil"
	Diesel              CargoType = "Diesel"
	Gasoline            CargoType = "Gasoline"
	Methane             CargoType = "Methane"
	Logs                CargoType = "Logs"
	Boards              CargoType = "Boards"
	Plywood             CargoType = "Plywood"
	Wheat               CargoType = "Wheat"
	Corn                CargoType = "Corn"
	Pigs                CargoType = "Pigs"
	Cows                CargoType = "Cows"
	Chickens            CargoType = "Chickens"
	Sheep               CargoType = "Sheep"
	Goats               CargoType = "Goats"
	Bread               CargoType = "Bread"
	DairyProducts       CargoType = "DairyProducts"
	MeatProducts        CargoType = "MeatProducts"
	CannedFood          CargoType = "CannedFood"
	CatFood             CargoType = "CatFood"
	SteelRolls          CargoType = "SteelRolls"
	SteelBillets        CargoType = "SteelBillets"
	SteelSlabs          CargoType = "SteelSlabs"
	SteelBentPlates     CargoType = "SteelBentPlates"
	SteelRails          CargoType = "SteelRails"
	ScrapMetal          CargoType = "ScrapMetal"
	ElectronicsIskar    CargoType = "ElectronicsIskar"
	ElectronicsKrugmann CargoType = "ElectronicsKrugmann"
	ElectronicsAAG      CargoType = "ElectronicsAAG"
	ElectronicsNovae    CargoType = "ElectronicsNovae"
	ElectronicsTraeg    CargoType = "ElectronicsTraeg"
	ToolsIskar          CargoType = "ToolsIskar"
	ToolsBrohm          CargoType = "ToolsBrohm"
	ToolsAAG            CargoType = "ToolsAAG"
	ToolsNovae          CargoType = "ToolsNovae"
	ToolsTraeg          CargoType = "ToolsTraeg"
	Furniture           CargoType = "Furniture"
	Pipes               CargoType = "Pipes"
	ClothingObco        CargoType = "ClothingObco"
	ClothingNeoGamma    CargoType = "ClothingNeoGamma"
	ClothingNovae       CargoType = "ClothingNovae"
	ClothingTraeg       CargoType = "ClothingTraeg"
	Medicine            CargoType = "Medicine"
	ChemicalsIskar      CargoType = "ChemicalsIskar"
	ChemicalsSperex     CargoType = "ChemicalsSperex"
	NewCars             CargoType = "NewCars"
	ImportedNewCars     CargoType = "ImportedNewCars"
	Tractors            CargoType = "Tractors"
	Excavators          CargoType = "Excavators"
	Alcohol             CargoType = "Alcohol"
	Acetylene           CargoType = "Acetylene"
	CryoOxygen          CargoType = "CryoOxygen"
	CryoHydrogen        CargoType = "CryoHydrogen"
	Argon               CargoType = "Argon"
	Nitrogen            CargoType = "Nitrogen"
	Ammonia             CargoType = "Ammonia"
	SodiumHydroxide     CargoType = "SodiumHydroxide"
	SpentNuclearFuel    CargoType = "SpentNuclearFuel"
	Ammunition          CargoType = "Ammunition"
	Biohazard           CargoType = "Biohazard"
	Tanks               CargoType = "Tanks"
	MilitaryTrucks      CargoType = "MilitaryTrucks"
	MilitarySupplies    CargoType = "MilitarySupplies"
	EmptySunOmni        CargoType = "EmptySunOmni"
	EmptyIskar          CargoType = "EmptyIskar"
	EmptyObco           CargoType = "EmptyObco"
	EmptyGoorsk         CargoType = "EmptyGoorsk"
	EmptyKrugmann       CargoType = "EmptyKrugmann"
	EmptyBrohm          CargoType = "EmptyBrohm"
	EmptyAAG            CargoType = "EmptyAAG"
	EmptySperex         CargoType = "EmptySperex"
	EmptyNovae          CargoType = "EmptyNovae"
	EmptyTraeg          CargoType = "EmptyTraeg"
	EmptyChemlek        CargoType = "EmptyChemlek"
	EmptyNeoGamma       CargoType = "EmptyNeoGamma"

	CategoryRaw     CargoCategory = "Raw"
	CategoryDanger  CargoCategory = "Danger"
	CategoryEasy    CargoCategory = "Easy"
	CategoryComplex CargoCategory = "Complex"
)

var cargoCategory = map[CargoType]CargoCategory{
	None:                CategoryRaw,
	Coal:                CategoryRaw,
	IronOre:             CategoryRaw,
	CrudeOil:            CategoryDanger,
	Diesel:              CategoryDanger,
	Gasoline:            CategoryDanger,
	Methane:             CategoryDanger,
	Logs:                CategoryRaw,
	Boards:              CategoryEasy,
	Plywood:             CategoryEasy,
	Wheat:               CategoryRaw,
	Corn:                CategoryRaw,
	Pigs:                CategoryRaw,
	Cows:                CategoryRaw,
	Chickens:            CategoryRaw,
	Sheep:               CategoryRaw,
	Goats:               CategoryRaw,
	Bread:               CategoryEasy,
	DairyProducts:       CategoryEasy,
	MeatProducts:        CategoryEasy,
	CannedFood:          CategoryEasy,
	CatFood:             CategoryEasy,
	SteelRolls:          CategoryEasy,
	SteelBillets:        CategoryEasy,
	SteelSlabs:          CategoryEasy,
	SteelBentPlates:     CategoryEasy,
	SteelRails:          CategoryEasy,
	ScrapMetal:          CategoryRaw,
	ElectronicsIskar:    CategoryComplex,
	ElectronicsKrugmann: CategoryComplex,
	ElectronicsAAG:      CategoryComplex,
	ElectronicsNovae:    CategoryComplex,
	ElectronicsTraeg:    CategoryComplex,
	ToolsIskar:          CategoryComplex,
	ToolsBrohm:          CategoryComplex,
	ToolsAAG:            CategoryComplex,
	ToolsNovae:          CategoryComplex,
	ToolsTraeg:          CategoryComplex,
	Furniture:           CategoryComplex,
	Pipes:               CategoryEasy,
	ClothingObco:        CategoryComplex,
	ClothingNeoGamma:    CategoryComplex,
	ClothingNovae:       CategoryComplex,
	ClothingTraeg:       CategoryComplex,
	Medicine:            CategoryComplex,
	ChemicalsIskar:      CategoryComplex,
	ChemicalsSperex:     CategoryComplex,
	NewCars:             CategoryComplex,
	ImportedNewCars:     CategoryComplex,
	Tractors:            CategoryComplex,
	Excavators:          CategoryComplex,
	Alcohol:             CategoryEasy,
	Acetylene:           CategoryDanger,
	CryoOxygen:          CategoryDanger,
	CryoHydrogen:        CategoryDanger,
	Argon:               CategoryDanger,
	Nitrogen:            CategoryDanger,
	Ammonia:             CategoryDanger,
	SodiumHydroxide:     CategoryDanger,
	SpentNuclearFuel:    CategoryDanger,
	Ammunition:          CategoryDanger,
	Biohazard:           CategoryDanger,
	Tanks:               CategoryDanger,
	MilitaryTrucks:      CategoryComplex,
	MilitarySupplies:    CategoryComplex,
	EmptySunOmni:        CategoryRaw,
	EmptyIskar:          CategoryRaw,
	EmptyObco:           CategoryRaw,
	EmptyGoorsk:         CategoryRaw,
	EmptyKrugmann:       CategoryRaw,
	EmptyBrohm:          CategoryRaw,
	EmptyAAG:            CategoryRaw,
	EmptySperex:         CategoryRaw,
	EmptyNovae:          CategoryRaw,
	EmptyTraeg:          CategoryRaw,
	EmptyChemlek:        CategoryRaw,
	EmptyNeoGamma:       CategoryRaw,
}

func (ct CargoType) String() string {
	return string(ct)
}

func (c CargoType) BaseWage() int {
	switch cargoCategory[c] {
	case CategoryRaw:
		return 400
	case CategoryDanger:
		return 1000
	case CategoryEasy:
		return 600
	case CategoryComplex:
		return 800
	}

	logrus.WithField("cargo", c).Info("unmapped cargo type for wage calc")
	return 0
}
