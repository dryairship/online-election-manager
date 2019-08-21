const mongoose = require('mongoose');

var CEOSchema = new mongoose.Schema({
    roll : String,
    name : String,
    email : String,
    username : String,
    password : String,
    authcode : String,
    publickey : String,
    privatekey : String
});

var CEO = mongoose.model('CEO', CEOSchema, 'ceo');

function getCEO() {
    return new Promise((resolve, reject) => {
        CEO.findOne(
            null,
            (err, ceo) => {
                if(err || !ceo)
                    reject(err);
                else
                    resolve(ceo);
            }
        );
    });
}

module.exports = {getCEO}