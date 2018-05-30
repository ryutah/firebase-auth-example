import firebase from "firebase/app";
import "firebase/auth";

const config = {
  apiKey: "[API_KEY]",
  authDomain: "[PROJECT_ID].firebaseapp.com",
  projectId: "[PROJECT_ID]",
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
