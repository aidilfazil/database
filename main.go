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

// Car model
type Car struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Make      string             `json:"make"`
	Model     string             `json:"model"`
	Year      int                `json:"year"`
	Type      string             `json:"type"`
	Available bool               `json:"available"`
}

// Customer model
type Customer struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Email         string             `json:"email"`
	PhoneNumber   string             `json:"phone_number"`
	DriversLicense string            `json:"drivers_license,omitempty"`
}

// Rental model
type Rental struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CarID            primitive.ObjectID `json:"car_id"`
	CustomerID       primitive.ObjectID `json:"customer_id"`
	RentalStartDate  string             `json:"rental_start_date"`
	RentalEndDate    string             `json:"rental_end_date"`
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

	// Define routes here

	log.Fatal(app.Listen(":5000"))
}

// Get all cars
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

// Create a new car
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

// Update a car
func updateCar(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}
	update := bson.M{"$set": c.BodyParser(new(Car))}
	_, err = carCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// Delete a car
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


// Get all customers
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

// Create a new customer
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

// Update a customer
func updateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}
	update := bson.M{"$set": c.BodyParser(new(Customer))}
	_, err = customerCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// Delete a customer
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

// Get all rentals
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

// Create a new rental
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

// Update a rental
func updateRental(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid rental ID"})
	}
	update := bson.M{"$set": c.BodyParser(new(Rental))}
	_, err = rentalCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// Delete a rental
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