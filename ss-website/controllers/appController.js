const bcrypt = require("bcryptjs");

const User = require("../models/User");
const Key = require("../models/Licence");

exports.landing_page = (req, res) => {
    res.render("index");
};

exports.login_get = (req, res) => {
    const error = req.session.error;
    delete req.session.error;
    res.render("login", { err: error });
};

exports.login_post = async (req, res) => {
    const { username, password } = req.body;

    const user = await User.findOne({ username });

    if (!user) {
        req.session.error = "Invalid Credentials";
        return res.redirect("/login");
    }

    const isMatch = await bcrypt.compare(password, user.password);

    if (!isMatch) {
        req.session.error = "Invalid Credentials";
        return res.redirect("/login");
    }

    req.session.isAuth = true;
    req.session.username = user.username;
    req.session.slot = user.slot;
    req.session.private_tokens = user.private_tokens;
    req.session.sniped = user.sniped;
    req.session.webhook = user.webhook;
    req.session.token = user.claim_token;
    res.redirect("/dashboard");
};

exports.register_get = (req, res) => {
    const error = req.session.error;
    delete req.session.error;
    res.render("register", { err: error });
};

exports.register_post = async (req, res) => {
    const { username, password, licence_key } = req.body;

    let user = await User.findOne({ username });

    if (user) {
        req.session.error = "User already exists";
        return res.redirect("/register");
    }

    console.log(licence_key)
    let key = await Key.findOne({ key: licence_key });
    console.log(key)

    if (!key) {
        req.session.error = "Licence key invalid";
        return res.redirect("/register");
    }

    if (key.used === true) {
        req.session.error = "Licence key already used";
        return res.redirect("/register");
    }

    const hasdPsw = await bcrypt.hash(password, 12);

    user = new User({
        username,
        password: hasdPsw,
        private_tokens: [],
        slot: key.slot,
        sniped: 0,
        webhook: 'None',
        claim_token: 'None'
    });


    await user.save();

    await Key.findOneAndUpdate({ key: licence_key }, {
        used: true
    })

    res.redirect("/login");
};

exports.dashboard_get = async (req, res) => {
    let user = await User.findOne({ username: req.session.username });
    req.session.slot = user.slot;
    req.session.private_tokens = user.private_tokens;
    req.session.sniped = user.sniped;
    req.session.webhook = user.webhook;
    req.session.token = user.claim_token;

    console.log(user);

    const username = req.session.username;
    const slot = req.session.slot;
    const private_tokens = req.session.private_tokens;
    const sniped = req.session.sniped;
    var webhook = req.session.webhook;
    var token = req.session.token;

    if (webhook.includes('discord.com')) {
        webhook = '✔️'
    } else {
        webhook = '❌'
    }

    if (['None', 'null'].includes(token)) {
        token = '❌'
    } else {
        token = '✔️'
    }

    console.log(req.session);
    res.render("dashboard", { name: username, slot: slot, private_tokens: private_tokens, sniped: sniped, webhook: webhook, token: token });
};

exports.update_hook_post = async (req, res) => {
    console.log(req.query.hook)
    const hook = req.query.hook;
    const username = req.session.username;
    console.log(username)

    await User.findOneAndUpdate({ username }, {
        webhook: hook
    });

    req.session.webhook = hook;
    res.redirect("/dashboard");
};

exports.update_token_post = async (req, res) => {
    console.log(req.query.token)
    const token = req.query.token;
    const username = req.session.username;
    console.log(username)

    await User.findOneAndUpdate({ username }, {
        claim_token: token
    });

    req.session.token = token;
    res.redirect("/dashboard");
};

exports.logout_post = (req, res) => {
    req.session.destroy((err) => {
        if (err) throw err;
        res.redirect("/login");
    });
};

exports.claim_key = async (req, res) => {
    const username = req.session.username;
    let key = await Key.findOne({ key: req.query.key });

    const user = await User.findOne({ username });

    if (!key) {
        req.session.error = "Licence key invalid";
        return res.redirect("/dashboard");
    }

    if (key.used == true) {
        req.session.error = "Licence key already used";
        return res.redirect("/dashboard");
    }

    await User.findOneAndUpdate({ username }, {
        slot: user.slot + key.slot
    });

    await Key.findOneAndUpdate({ key: req.query.key }, {
        used: true
    });

    res.send('/dashboard');
}