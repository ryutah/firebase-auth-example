import firebase from "firebase/app";
import "firebase/auth";

const config = {
  apiKey: "[API_KEY]",
  authDomain: "[PROJECT].firebaseapp.com",
  projectId: "[PROJECT]",
};

firebase.initializeApp(config);

window.firebase = firebase;

const $email = document.querySelector("#email");
const $pass = document.querySelector("#pass");

document.querySelector("#login").addEventListener("click", (e) => {
  e.preventDefault();
  console.log(`Login ${$email.value}:${$pass.value}`);

  firebase
    .auth()
    .signInWithEmailAndPassword($email.value, $pass.value)
    .then((result) => result.user.getIdToken())
    .then((idToken) => {
      const body = { idToken: idToken };
      fetch("/createsession", {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json; charset=utf-8",
        },
        body: JSON.stringify(body),
      });
    })
    .then(console.log)
    .then(() => firebase.auth().signOut())
    .catch(console.error);
});

document.querySelector("#logout").addEventListener("click", (e) => {
  e.preventDefault();
  console.log("Click logout");

  firebase
    .auth()
    .signOut()
    .then(() => console.log("Success sign out"))
    .catch(console.error);
});

document.querySelector("#verify").addEventListener("click", (e) => {
  e.preventDefault();

  fetch("/verifysession", { credentials: "include" })
    .then(console.log)
    .catch(console.error);
});
