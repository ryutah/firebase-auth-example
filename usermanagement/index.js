import firebase from "firebase/app";
import "firebase/auth";

const config = {
  apiKey: process.env.API_KEY,
  authDomain: `${process.env.PROJECT_ID}.firebaseapp.com`,
  projectId: process.env.PROJECT_ID,
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
      const data = new URLSearchParams();
      data.append("idToken", idToken);
      data.append(
        "csrfToken",
        document.querySelector("input[name=csrfToken]").value
      );
      fetch("/signin", {
        credentials: "same-origin",
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
          "X-CSRF-Token": document.querySelector("input[name=csrfToken]").value,
        },
        body: data,
      });
    })
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

document.querySelector("#create").addEventListener("click", (e) => {
  e.preventDefault();
  console.log(`Create ${$email.value}:${$pass.value}`);

  firebase
    .auth()
    .createUserWithEmailAndPassword($email.value, $pass.value)
    .catch(console.error);
});

firebase.auth().onAuthStateChanged((usr) => {
  if (usr) {
    console.log(usr);
  }
  window.currentUser = usr;
});
