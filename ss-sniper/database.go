package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Slot struct {
	ID            primitive.ObjectID `bson:"_id"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	Slot          int                `bson:"slot"`
	Sniped        int                `bson:"sniped"`
	ClaimToken    string             `bson:"claim_token"`
	WebhookUrl    string             `bson:"webhook"`
	Username      string             `bson:"username"`
	PrivateTokens []string           `bson:"private_tokens"`
}

type ClaimerAccount struct {
	webhook         *hook
	AccSlot         Slot
	base_request    fasthttp.Request
	payement_method string
}

var collection *mongo.Collection
var ctx = context.TODO()

func init_db() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://the_user:the_passsword@cluster0.r9vxk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority") // replace by your creds
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log(err.Error())
	}

	client.Connect(ctx)
	collection = client.Database("myFirstDatabase").Collection("users")
}

func get_slot() []*ClaimerAccount {
	var slots []*ClaimerAccount

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log(err.Error())
	}

	for cursor.Next(ctx) {
		var slot Slot
		err := cursor.Decode(&slot)
		if err != nil {
			log(err.Error())
		}

		if slot.ClaimToken != "None" {
			if check_token(slot.ClaimToken) {
				if slot.Sniped < slot.Slot {
					AccSlot := &ClaimerAccount{
						AccSlot: slot,
					}

					AccSlot.initialize_header()
					AccSlot.get_fingerprint()
					AccSlot.get_cookies()
					AccSlot.get_payement_source_id()

					if slot.WebhookUrl != "None" {
						url := strings.Split(slot.WebhookUrl, "discord.com/api/webhooks/")[1]
						h, _ := new_hook(strings.Split(url, "/")[1], strings.Split(url, "/")[0])
						AccSlot.webhook = h
					}

					slots = append(slots, AccSlot)
				}
			} else {
				collection.UpdateOne(ctx, bson.M{"_id": slot.ID}, bson.M{"$set": bson.M{"claim_token": "None"}})
				if slot.WebhookUrl != "None" {
					url := strings.Split(slot.WebhookUrl, "discord.com/api/webhooks/")[1]
					h, _ := new_hook(strings.Split(url, "/")[1], strings.Split(url, "/")[0])
					h.send_hook("Your token was invalid, please change it.")
					log(fmt.Sprintf("Invalid token for %s", slot.Username))
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log(err.Error())
	}

	return slots
}

func (account *ClaimerAccount) update_sniped() {
	collection.UpdateOne(ctx, bson.M{"_id": account.AccSlot.ID}, bson.M{"$set": bson.M{"sniped": account.AccSlot.Sniped}})
}
