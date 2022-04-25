const mongoose = require("mongoose");
const Schema = mongoose.Schema;

const KeySchema = new Schema({
  key: {
    type: String,
    required: true,
  },
  used: {
      type: Boolean,
      required: true
  },
  slot: {
    type: Number,
    required: true
  }
});

module.exports = mongoose.model("Key", KeySchema);
