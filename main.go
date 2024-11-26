package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB setup
var client *mongo.Client
var patientCollection *mongo.Collection
var appointmentCollection *mongo.Collection

type Patient struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Age   int                `json:"age" bson:"age"`
	Phone string             `json:"phone" bson:"phone"`
}

type Appointment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PatientID primitive.ObjectID `json:"patientId" bson:"patientId"`
	Date      time.Time          `json:"date" bson:"date"`
	Reason    string             `json:"reason" bson:"reason"`
}

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with a specific domain if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Connect to MongoDB
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://pesatoken:sBTiy1lKhZvAaZta@newcluster.sm1ec.mongodb.net/?retryWrites=true&w=majority&appName=NewCluster"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	// Collections
	patientCollection = client.Database("clinic").Collection("patients")
	appointmentCollection = client.Database("clinic").Collection("appointments")

	// Create a new multiplexer
	mux := http.NewServeMux()

	// // Routes
	// router := mux.NewRouter()

	mux.HandleFunc("/patients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPatients(w, r)
		case http.MethodPost:
			addPatient(w, r)
		case http.MethodPut:
			editPatient(w, r)
		case http.MethodDelete:
			deletePatient(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/appointments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAppointments(w, r)
		case http.MethodPost:
			scheduleAppointment(w, r)
		case http.MethodPut:
			rescheduleAppointment(w, r)
		case http.MethodDelete:
			deleteAppointment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Wrap routes with CORS middleware
	handlerWithCORS := enableCORS(mux)

	// Start the server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

func addPatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	_ = json.NewDecoder(r.Body).Decode(&patient)

	// Assign a random ObjectID
	patient.ID = primitive.NewObjectID()

	_, err := patientCollection.InsertOne(context.Background(), patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the generated ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": patient.ID.Hex()})
}

func editPatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var patient Patient
	_ = json.NewDecoder(r.Body).Decode(&patient)
	update := bson.M{"$set": patient}
	_, err := patientCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getPatients(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	nameFilter := r.URL.Query().Get("name")   // Filter by name
	phoneFilter := r.URL.Query().Get("phone") // Filter by phone
	ageFilter := r.URL.Query().Get("age")     // Filter by age

	// Build the filter based on provided query parameters
	filter := bson.M{}
	if nameFilter != "" {
		filter["name"] = nameFilter
	}
	if phoneFilter != "" {
		filter["phone"] = phoneFilter
	}
	if ageFilter != "" {
		filter["age"] = ageFilter
	}

	// Retrieve documents from the collection
	collection := client.Database("clinic").Collection("patients")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to fetch patients", http.StatusInternalServerError)
		log.Printf("Error fetching patients: %v", err)
		return
	}
	defer cursor.Close(context.Background())

	// Decode documents into a slice
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		http.Error(w, "Failed to decode patients", http.StatusInternalServerError)
		log.Printf("Error decoding patients: %v", err)
		return
	}

	// Convert the results to JSON and send as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	_, err := patientCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func scheduleAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment Appointment
	_ = json.NewDecoder(r.Body).Decode(&appointment)

	// Assign a random ObjectID
	appointment.ID = primitive.NewObjectID()

	_, err := appointmentCollection.InsertOne(context.Background(), appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the generated ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": appointment.ID.Hex()})
}

func rescheduleAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var appointment Appointment
	_ = json.NewDecoder(r.Body).Decode(&appointment)
	update := bson.M{"$set": bson.M{"date": appointment.Date, "reason": appointment.Reason}}
	_, err := appointmentCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getAppointments(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	patientIdFilter := r.URL.Query().Get("patientId") // Filter by patientId
	dateFilter := r.URL.Query().Get("date")           // Filter by date
	reasonFilter := r.URL.Query().Get("reason")       // Filter by reason

	// Build the filter based on provided query parameters
	filter := bson.M{}
	if patientIdFilter != "" {
		filter["patientId"] = patientIdFilter
	}
	if dateFilter != "" {
		filter["date"] = dateFilter
	}
	if reasonFilter != "" {
		filter["reason"] = reasonFilter
	}

	// Retrieve documents from the collection
	collection := client.Database("clinic").Collection("appointments")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
		log.Printf("Error fetching appointments: %v", err)
		return
	}
	defer cursor.Close(context.Background())

	// Decode documents into a slice
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		http.Error(w, "Failed to decode appointments", http.StatusInternalServerError)
		log.Printf("Error decoding appointments: %v", err)
		return
	}

	// Convert the results to JSON and send as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func deleteAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	_, err := appointmentCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
