const express = require("express");
const session = require("express-session");
const MongoDBStore = require("connect-mongodb-session")(session);
const appController = require("./controllers/appController");
const private_api = require('./api/private_api');
const isAuth = require("./middleware/is-auth");
const connectDB = require("./config/db");
const mongoURI = 'mongodb+srv://your_username:your_pass@cluster0.r9vxk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority';

const app = express();
connectDB();

const store = new MongoDBStore({
  uri: mongoURI,
  collection: "mySessions",
});


app.set("view engine", "ejs");
app.use(express.static(__dirname + '/public'));
app.use(express.urlencoded({ extended: true }));

app.use(
  session({
    secret: "secret",
    resave: false,
    saveUninitialized: false,
    store: store,
  })
);

//=================== Routes
// Landing Page
app.get("/", appController.landing_page);

// Login Page
app.get("/login", appController.login_get);
app.post("/login", appController.login_post);

// Register Page
app.get("/register", appController.register_get);
app.post("/register", appController.register_post);

// Dashboard Page
app.get("/dashboard", isAuth, appController.dashboard_get);
app.post("/update_hook/", isAuth, appController.update_hook_post);
app.post("/update_token/", isAuth, appController.update_token_post);
app.post("/logout", appController.logout_post);
app.post('/claim_key/', isAuth, appController.claim_key)

// Private api
app.get('/private_api/send_account', private_api.send_account)
app.get('/private_api/report_claim', private_api.report_claim)

app.listen(8888, console.log("App Running on http://localhost:8888"));



/*
  i am not webdev so the most part if this website was taked from git and edited
*/