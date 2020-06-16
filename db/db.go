package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/strosel/sprak/cards"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func http500(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}

type server struct {
	coll *mongo.Collection
}

func NewServer() (*server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	return &server{client.Database("sprak").Collection("decks")}, err
}

func (s *server) getDecks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := s.coll.Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$group", bson.M{
			"_id": "null",
			"res": bson.M{"$addToSet": "$name"},
		}}},
	})

	if http500(w, err) {
		return
	}

	decks := []bson.M{}
	err = res.All(ctx, &decks)
	if http500(w, err) {
		return
	}

	data, err := json.Marshal(decks[0]["res"])
	if http500(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *server) handleDeck(w http.ResponseWriter, r *http.Request) {
	re, _ := regexp.Compile(`^\/?Decks\/([^\/]+)$`)
	name := re.FindAllStringSubmatch(r.URL.Path, -1)[0][1]

	deckString := []byte(r.FormValue("deck"))
	var deck cards.Deck
	if len(deckString) > 0 {
		err := json.Unmarshal(deckString, &deck)
		if http500(w, err) {
			return
		}

		if name != deck.Name {
			http500(w, fmt.Errorf("Mismatched names in path and object"))
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		s.getDeck(w, ctx, name)
	case http.MethodPost:
		s.addDeck(w, ctx, deck)
	case http.MethodPut:
		s.updateDeck(w, ctx, deck)
	case http.MethodDelete:
		s.removeDeck(w, ctx, name)
	default:
		http500(w, fmt.Errorf("Unsupported method"))
	}
}

func (s *server) getDeck(w http.ResponseWriter, ctx context.Context, name string) {
	res := s.coll.FindOne(ctx, bson.M{
		"name": name,
	})

	d := cards.Deck{}
	err := res.Decode(&d)
	if http500(w, err) {
		return
	}

	data, err := json.Marshal(d)
	if http500(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *server) addDeck(w http.ResponseWriter, ctx context.Context, deck cards.Deck) error {
	_, err := s.coll.InsertOne(ctx, deck)
	return err
}

func (s *server) removeDeck(w http.ResponseWriter, ctx context.Context, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s.coll.FindOneAndDelete(ctx, bson.M{
		"name": name,
	})
}

func (s *server) updateDeck(w http.ResponseWriter, ctx context.Context, deck cards.Deck) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s.coll.FindOneAndUpdate(ctx, bson.M{
		"name": deck.Name,
	}, deck)
}

func main() {
	s, err := NewServer()
	if err != nil {
		log.Fatal()
	}
	http.HandleFunc("/Decks", s.getDecks)
	http.HandleFunc("/Decks/", s.handleDeck)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
