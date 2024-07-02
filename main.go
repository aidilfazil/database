package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

type Car struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Make      string             `json:"make"`
	Model     string             `json:"model"`
	Year      int                `json:"year"`
	Type      string             `json:"type"`
	Available bool               `json:"available"`
}

type Customer struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	PhoneNumber    string             `json:"phone_number"`
	DriversLicense string             `json:"drivers_license,omitempty"`
}

type Rental struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CarID           primitive.ObjectID `json:"car_id"`
	CustomerID      primitive.ObjectID `json:"customer_id"`
	RentalStartDate string             `json:"rental_start_date"`
	RentalEndDate   string             `json:"rental_end_date"`
}

var carCollection *mongo.Collection
var customerCollection *mongo.Collection
var rentalCollection *mongo.Collection

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	carCollection = client.Database("car_rental_db").Collection("cars")
	customerCollection = client.Database("car_rental_db").Collection("customers")
	rentalCollection = client.Database("car_rental_db").Collection("rentals")

	app := fiber.New()

	app.Get("/api/car/:id", getCarByID)
	app.Get("/api/cars", getCars)
	app.Post("/api/cars", createCar)
	app.Patch("/api/cars/:id", updateCar)
	app.Delete("/api/cars/:id", deleteCar)

	app.Get("/api/customer/:id", getCustomerByID)
	app.Get("/api/customers", getCustomers)
	app.Post("/api/customers", createCustomer)
	app.Patch("/api/customers/:id", updateCustomer)
	app.Delete("/api/customers/:id", deleteCustomer)

	app.Get("/api/rental/:id", getRentalByID)
	app.Get("/api/rentals", getRentals)
	app.Post("/api/rentals", createRental)
	app.Patch("/api/rentals/:id", updateRental)
	app.Delete("/api/rentals/:id", deleteRental)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getCarByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}

	var car Car
	err = carCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&car)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Car not found"})
		}
		return err
	}

	return c.JSON(car)
}

func getCars(c *fiber.Ctx) error {
	var cars []Car
	cursor, err := carCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var car Car
		cursor.Decode(&car)
		cars = append(cars, car)
	}
	return c.JSON(cars)
}

func createCar(c *fiber.Ctx) error {
	car := new(Car)
	if err := c.BodyParser(car); err != nil {
		return err
	}
	car.ID = primitive.NewObjectID()
	_, err := carCollection.InsertOne(context.Background(), car)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(car)
}

func updateCar(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}

	var car Car
	if err := c.BodyParser(&car); err != nil {
		return err
	}

	update := bson.M{}
	if car.Make != "" {
		update["make"] = car.Make
	}
	if car.Model != "" {
		update["model"] = car.Model
	}
	if car.Year != 0 {
		update["year"] = car.Year
	}
	if car.Type != "" {
		update["type"] = car.Type
	}
	// Check presence of the 'available' field in the request body
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if available, ok := body["available"]; ok {
		update["available"] = available
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No valid fields to update"})
	}

	updateSet := bson.M{"$set": update}
	_, err = carCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, updateSet)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteCar(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}
	_, err = carCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func getCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	var customer Customer
	err = customerCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
		}
		return err
	}

	return c.JSON(customer)
}

func getCustomers(c *fiber.Ctx) error {
	var customers []Customer
	cursor, err := customerCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var customer Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	return c.JSON(customers)
}

func createCustomer(c *fiber.Ctx) error {
	customer := new(Customer)
	if err := c.BodyParser(customer); err != nil {
		return err
	}
	customer.ID = primitive.NewObjectID()
	_, err := customerCollection.InsertOne(context.Background(), customer)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(customer)
}

func updateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	var customer Customer
	if err := c.BodyParser(&customer); err != nil {
		return err
	}

	update := bson.M{}
	if customer.Name != "" {
		update["name"] = customer.Name
	}
	if customer.Email != "" {
		update["email"] = customer.Email
	}
	if customer.PhoneNumber != "" {
		update["phone_number"] = customer.PhoneNumber
	}
	if customer.DriversLicense != "" {
		update["drivers_license"] = customer.DriversLicense
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No valid fields to update"})
	}

	updateSet := bson.M{"$set": update}
	_, err = customerCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, updateSet)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}
	_, err = customerCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func getRentalByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rental ID"})
	}

	var rental Rental
	err = rentalCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&rental)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Rental not found"})
		}
		return err
	}

	return c.JSON(rental)
}

func getRentals(c *fiber.Ctx) error {
	var rentals []Rental
	cursor, err := rentalCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var rental Rental
		cursor.Decode(&rental)
		rentals = append(rentals, rental)
	}
	return c.JSON(rentals)
}

func createRental(c *fiber.Ctx) error {
	rental := new(Rental)
	if err := c.BodyParser(rental); err != nil {
		return err
	}
	rental.ID = primitive.NewObjectID()
	_, err := rentalCollection.InsertOne(context.Background(), rental)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(rental)
}

func updateRental(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rental ID"})
	}

	var rental Rental
	if err := c.BodyParser(&rental); err != nil {
		return err
	}

	update := bson.M{}
	if !rental.CarID.IsZero() {
		update["car_id"] = rental.CarID
	}
	if !rental.CustomerID.IsZero() {
		update["customer_id"] = rental.CustomerID
	}
	if rental.RentalStartDate != "" {
		update["rental_start_date"] = rental.RentalStartDate
	}

	if rental.RentalEndDate != "" {
		update["rental_end_date"] = rental.RentalEndDate
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No valid fields to update"})
	}

	updateSet := bson.M{"$set": update}
	_, err = rentalCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, updateSet)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteRental(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rental ID"})
	}
	_, err = rentalCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}
