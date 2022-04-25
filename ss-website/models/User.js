const mongoose = require("mongoose");
const Schema = mongoose.Schema;

const userSchema = new Schema({
  username: {
    type: String,
    required: true,
  },
  password: {
    type: String,
    required: true,
  },
  private_tokens: {
    type: Array,
    required: true,
  },
  slot: {
    type: Number,
    required: true
  },
  sniped: {
    type: Number,
    required: true
  },
  webhook: {
    type: String,
    required: true
  },
  claim_token: {
    type: String,
    required: true
  }
});

module.exports = mongoose.model("User", userSchema);
