const bcrypt = require("bcryptjs");
const User = require("../models/User");
const Key = require("../models/Licence");

const private_key = "13374859411567JFHIEHFIUHNZ";

exports.send_account = async (req, res) => {
    if (req.query.key != private_key) {
        res.send('nop');
        return;
    }

    const users = await User.find({});
    console.log(users);

    res.send(users);
}

exports.report_claim = async (req, res) => {
    if (req.query.key != private_key) {
        res.send('nop');
        return;
    }

    try {
        const user = await User.findOne({ username: req.query.username });
        const query_sniped = parseInt(req.query.sniped, 10);

        console.log(req.query)
        
        var sniped = user.sniped + query_sniped;
        console.log(sniped)
        console.log(user.slot)

        if (sniped > user.slot) {
            console.log('max snip resized')
            sniped = user.slot;
        } else {
            console.log('ok')
        }

        await User.findOneAndUpdate({ username: req.query.username }, {
            sniped: sniped
        });

        res.send('ok');
    } catch (err) {
        res.send(`Error: ${err}`);
    }
}

