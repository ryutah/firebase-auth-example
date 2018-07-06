const functions = require("firebase-functions");

exports.handleNewUser = functions.auth.user().onCreate((user) => {
  console.log("Receive new user!!");
  console.log(user);
  return "Success Call!!";
});

exports.handleDeleteUser = functions.auth.user().onDelete((user) => {
  console.log("Delete user!!");
  console.log(user);
  return "Success Call!!";
});
