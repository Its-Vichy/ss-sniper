const mongoose = require("mongoose");

const connectDB = async () => {
  try {
    await mongoose.connect('mongodb+srv://your_username:your_pass@cluster0.r9vxk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority', {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });

    console.log("MongoDB connected");
  } catch (error) {
    console.log("Something went wrong with Database connection", error);
    process.exit(1);
  }
};

module.exports = connectDB;
