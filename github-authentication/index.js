import firebase from "firebase/app";
import "firebase/auth";

const config = {
  apiKey: process.env.API_KEY,
  authDomain: `${process.env.PROJECT_ID}.firebaseapp.com`,
  projectId: process.env.PROJECT_ID,
};

firebase.initializeApp(config);

window.firebase = firebase;

document.querySelector("#login").addEventListener("click", () => {
  const provider = new firebase.auth.GithubAuthProvider();
  firebase.auth().signInWithRedirect(provider);
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

document.querySelector("#getoauth").addEventListener("click", () => {
  const provider = new firebase.auth.GithubAuthProvider();
  firebase.auth().currentUser.reauthenticateWithRedirect(provider);
});

firebase.auth().onAuthStateChanged((usr) => {
  if (usr) {
    console.log(usr);
  }
  window.currentUser = usr;
});

firebase
  .auth()
  .getRedirectResult()
  .then((result) => {
    console.log(result);
    window.redirectResult = result;
  })
  .catch(console.error);
